package service

import (
	"context"
	"database/sql"
	"testing"

	mockDAO "dealer/internal/mock/dao"
	mockService "dealer/internal/mock/service"
	"dealer/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DealerTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	db           *sql.DB
	mockDB       sqlmock.Sqlmock
	mockGormDB   *gorm.DB
	mockBuyBook  *mockService.MockOrderBookInterface
	mockSellBook *mockService.MockOrderBookInterface
	mockOrderDAO *mockDAO.MockOrderInterface
	mockDealDAO  *mockDAO.MockDealInterface
	svc          *Dealer
}

func (t *DealerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	var err error
	t.db, t.mockDB, err = sqlmock.New()
	if err != nil {
		t.Failf("err", "an error '%s' was not expected when opening a stub database connection", err)
	}

	t.mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      t.db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		t.Failf("err", "an error '%s' was not expected when opening a stub database connection", err)
	}

	t.mockBuyBook = mockService.NewMockOrderBookInterface(t.ctrl)
	t.mockSellBook = mockService.NewMockOrderBookInterface(t.ctrl)
	t.mockOrderDAO = mockDAO.NewMockOrderInterface(t.ctrl)
	t.mockDealDAO = mockDAO.NewMockDealInterface(t.ctrl)
	t.svc = &Dealer{
		db:       t.mockGormDB,
		orderDAO: t.mockOrderDAO,
		dealDAO:  t.mockDealDAO,
		buyBook:  t.mockBuyBook,
		sellBook: t.mockSellBook,
	}
}

func (t *DealerTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.db.Close()
}

func TestDealerTestSuite(t *testing.T) {
	suite.Run(t, new(DealerTestSuite))
}

func (t *DealerTestSuite) TestProcessOrder() {
	tests := []struct {
		name     string
		order    *models.Order
		fn       func()
		hasError bool
	}{
		{
			name:  "Process order cancel",
			order: &models.Order{ID: 1, IsCancel: true},
			fn: func() {
				t.mockBuyBook.EXPECT().RemoveOrder(int64(1))
				t.mockSellBook.EXPECT().RemoveOrder(int64(1))
			},
		},
		{
			name: "Process buy order on market price not fulfil",
			order: &models.Order{
				ID:             1,
				OrderType:      models.OrderTypeBuy,
				Quantity:       1,
				RemainQuantity: 1,
				PriceType:      models.PriceTypeMarket,
			},
			fn: func() {
				t.mockSellBook.EXPECT().Peek().Return(nil)
				t.mockBuyBook.EXPECT().AddOrder(&models.Order{
					ID:             1,
					OrderType:      models.OrderTypeBuy,
					Quantity:       1,
					RemainQuantity: 1,
					PriceType:      models.PriceTypeMarket,
				})
				t.mockDB.ExpectBegin()
				t.mockOrderDAO.EXPECT().
					BulkUpdate(context.Background(), gomock.Any(), []*models.Order{
						{
							ID:             1,
							OrderType:      models.OrderTypeBuy,
							Quantity:       1,
							RemainQuantity: 1,
							PriceType:      models.PriceTypeMarket,
						},
					}).
					Return(nil)
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name: "Process buy order on limit price not fulfil price not match",
			order: &models.Order{
				ID:             1,
				OrderType:      models.OrderTypeBuy,
				Quantity:       1,
				RemainQuantity: 1,
				PriceType:      models.PriceTypeLimit,
				Price:          5,
			},
			fn: func() {
				t.mockSellBook.EXPECT().Peek().Return(&models.Order{
					ID:             1,
					OrderType:      models.OrderTypeSell,
					Quantity:       1,
					RemainQuantity: 1,
					PriceType:      models.PriceTypeLimit,
					Price:          10,
				})
				t.mockBuyBook.EXPECT().AddOrder(&models.Order{
					ID:             1,
					OrderType:      models.OrderTypeBuy,
					Quantity:       1,
					RemainQuantity: 1,
					PriceType:      models.PriceTypeLimit,
					Price:          5,
				})
				t.mockDB.ExpectBegin()
				t.mockOrderDAO.EXPECT().
					BulkUpdate(context.Background(), gomock.Any(), []*models.Order{
						{
							ID:             1,
							OrderType:      models.OrderTypeBuy,
							Quantity:       1,
							RemainQuantity: 1,
							PriceType:      models.PriceTypeLimit,
							Price:          5,
						},
					}).
					Return(nil)
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name: "Process buy order on market price fulfil seller partial fulfil",
			order: &models.Order{
				ID:             1,
				OrderType:      models.OrderTypeBuy,
				Quantity:       1,
				RemainQuantity: 1,
				PriceType:      models.PriceTypeMarket,
			},
			fn: func() {
				t.mockSellBook.EXPECT().
					Peek().
					Return(&models.Order{
						ID:             2,
						OrderType:      models.OrderTypeSell,
						Quantity:       2,
						RemainQuantity: 2,
						PriceType:      models.PriceTypeLimit,
						Price:          10,
					})
				t.mockDB.ExpectBegin()
				t.mockOrderDAO.EXPECT().
					BulkUpdate(context.Background(), gomock.Any(), []*models.Order{
						{
							ID:             2,
							OrderType:      models.OrderTypeSell,
							Quantity:       2,
							RemainQuantity: 1,
							PriceType:      models.PriceTypeLimit,
							Price:          10,
						},
						{
							ID:             1,
							OrderType:      models.OrderTypeBuy,
							Quantity:       1,
							RemainQuantity: 0,
							PriceType:      models.PriceTypeMarket,
						},
					}).
					Return(nil)
				t.mockDealDAO.EXPECT().
					Insert(context.Background(), gomock.Any(), []*models.Deal{
						{
							TakerOrderID: 1,
							MakerOrderID: 2,
							Quantity:     1,
							Price:        10,
						},
					})
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name: "Process buy order on market price fulfil seller fulfil",
			order: &models.Order{
				ID:             1,
				OrderType:      models.OrderTypeBuy,
				Quantity:       1,
				RemainQuantity: 1,
				PriceType:      models.PriceTypeMarket,
			},
			fn: func() {
				t.mockSellBook.EXPECT().
					Peek().
					Return(&models.Order{
						ID:             2,
						OrderType:      models.OrderTypeSell,
						Quantity:       2,
						RemainQuantity: 1,
						PriceType:      models.PriceTypeLimit,
						Price:          10,
					})
				t.mockSellBook.EXPECT().Dequeue().Return(
					&models.Order{
						ID:             2,
						OrderType:      models.OrderTypeSell,
						Quantity:       2,
						RemainQuantity: 0,
						PriceType:      models.PriceTypeLimit,
						Price:          10,
					})
				t.mockDB.ExpectBegin()
				t.mockOrderDAO.EXPECT().
					BulkUpdate(context.Background(), gomock.Any(), []*models.Order{
						{
							ID:             2,
							OrderType:      models.OrderTypeSell,
							Quantity:       2,
							RemainQuantity: 0,
							PriceType:      models.PriceTypeLimit,
							Price:          10,
						},
						{
							ID:             1,
							OrderType:      models.OrderTypeBuy,
							Quantity:       1,
							RemainQuantity: 0,
							PriceType:      models.PriceTypeMarket,
						},
					}).
					Return(nil)
				t.mockDealDAO.EXPECT().
					Insert(context.Background(), gomock.Any(), []*models.Deal{
						{
							TakerOrderID: 1,
							MakerOrderID: 2,
							Quantity:     1,
							Price:        10,
						},
					})
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name: "Process buy order on market price fulfil seller market price partial fulfil",
			order: &models.Order{
				ID:             1,
				OrderType:      models.OrderTypeBuy,
				Quantity:       1,
				RemainQuantity: 1,
				PriceType:      models.PriceTypeMarket,
			},
			fn: func() {
				t.svc.lastTradingPrice = 20
				t.mockSellBook.EXPECT().
					Peek().
					Return(&models.Order{
						ID:             2,
						OrderType:      models.OrderTypeSell,
						Quantity:       2,
						RemainQuantity: 2,
						PriceType:      models.PriceTypeMarket,
					})
				t.mockDB.ExpectBegin()
				t.mockOrderDAO.EXPECT().
					BulkUpdate(context.Background(), gomock.Any(), []*models.Order{
						{
							ID:             2,
							OrderType:      models.OrderTypeSell,
							Quantity:       2,
							RemainQuantity: 1,
							PriceType:      models.PriceTypeMarket,
						},
						{
							ID:             1,
							OrderType:      models.OrderTypeBuy,
							Quantity:       1,
							RemainQuantity: 0,
							PriceType:      models.PriceTypeMarket,
						},
					}).
					Return(nil)
				t.mockDealDAO.EXPECT().
					Insert(context.Background(), gomock.Any(), []*models.Deal{
						{
							TakerOrderID: 1,
							MakerOrderID: 2,
							Quantity:     1,
							Price:        20,
						},
					})
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			err := t.svc.ProcessOrder(context.Background(), test.order)
			t.Equal(test.hasError, err != nil)
		})
	}
}

func (t *DealerTestSuite) TestIsPriceMatch() {
	tests := []struct {
		name       string
		takerOrder *models.Order
		price      float64
		expected   bool
	}{
		{
			name: "Buyer limit price same price match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeBuy,
				PriceType: models.PriceTypeLimit,
				Price:     10,
			},
			price:    10,
			expected: true,
		},
		{
			name: "Buyer limit price lower price match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeBuy,
				PriceType: models.PriceTypeLimit,
				Price:     10,
			},
			price:    9,
			expected: true,
		},
		{
			name: "Buyer limit price higher price not match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeBuy,
				PriceType: models.PriceTypeLimit,
				Price:     10,
			},
			price:    11,
			expected: false,
		},
		{
			name: "Buyer market price match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeBuy,
				PriceType: models.PriceTypeMarket,
			},
			price:    11,
			expected: true,
		},
		{
			name: "Seller limit price same price match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeSell,
				PriceType: models.PriceTypeLimit,
				Price:     10,
			},
			price:    10,
			expected: true,
		},
		{
			name: "Seller limit price lower price not match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeSell,
				PriceType: models.PriceTypeLimit,
				Price:     10,
			},
			price:    9,
			expected: false,
		},
		{
			name: "Buyer limit price higher price match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeSell,
				PriceType: models.PriceTypeLimit,
				Price:     10,
			},
			price:    11,
			expected: true,
		},
		{
			name: "Seller market price match",
			takerOrder: &models.Order{
				OrderType: models.OrderTypeSell,
				PriceType: models.PriceTypeMarket,
			},
			price:    11,
			expected: true,
		},
	}

	for _, test := range tests {
		actual := isPriceMatch(test.takerOrder, test.price)
		t.Equal(test.expected, actual)
	}
}
