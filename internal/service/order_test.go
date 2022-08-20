package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	mockDAO "dealer/internal/mock/dao"
	mockSDK "dealer/internal/mock/sdk"
	"dealer/internal/models"
)

type OrderTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	db           *sql.DB
	mockDB       sqlmock.Sqlmock
	mockGormDB   *gorm.DB
	mockOrderDAO *mockDAO.MockOrderInterface
	mockChannel  *mockSDK.MockAMQPChannel
	svc          *OrderProcessor
}

func (t *OrderTestSuite) SetupTest() {
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

	t.mockOrderDAO = mockDAO.NewMockOrderInterface(t.ctrl)
	t.mockChannel = mockSDK.NewMockAMQPChannel(t.ctrl)
	t.svc = NewOrderProcessor(t.mockChannel, "name", t.mockGormDB, t.mockOrderDAO)
}

func (t *OrderTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.db.Close()
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (t *OrderTestSuite) TestNewOrder() {
	order := &models.Order{
		ID:             1,
		OrderType:      models.OrderTypeBuy,
		Quantity:       10,
		RemainQuantity: 10,
		PriceType:      models.PriceTypeLimit,
		Price:          10,
		IsCancel:       false,
	}

	tests := []struct {
		name     string
		fn       func()
		hasError bool
	}{
		{
			name: "New order normal",
			fn: func() {
				t.mockOrderDAO.EXPECT().
					Insert(context.Background(), t.mockGormDB, order).
					Return(nil)
				data, _ := json.Marshal(order)
				t.mockChannel.EXPECT().
					PublishWithContext(context.Background(), "", "name", false, false, amqp.Publishing{ContentType: "application/json", Body: data}).
					Return(nil)
			},
			hasError: false,
		},
		{
			name: "New order insert database failed",
			fn: func() {
				t.mockOrderDAO.EXPECT().
					Insert(context.Background(), t.mockGormDB, order).
					Return(errors.New(""))
			},
			hasError: true,
		},
		{
			name: "New order publish message queue failed",
			fn: func() {
				t.mockOrderDAO.EXPECT().
					Insert(context.Background(), t.mockGormDB, order).
					Return(nil)
				data, _ := json.Marshal(order)
				t.mockChannel.EXPECT().
					PublishWithContext(context.Background(), "", "name", false, false, amqp.Publishing{ContentType: "application/json", Body: data}).
					Return(errors.New(""))
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			err := t.svc.NewOrder(context.Background(), order)
			t.Equal(test.hasError, err != nil)
		})
	}
}

func (t *OrderTestSuite) TestCancelOrder() {
	order := &models.Order{ID: 1, IsCancel: true}
	tests := []struct {
		name     string
		fn       func()
		hasError bool
	}{
		{
			name: "Cancel order normal",
			fn: func() {
				t.mockOrderDAO.EXPECT().
					Update(context.Background(), t.mockGormDB, order).
					Return(nil)
				data, _ := json.Marshal(order)
				t.mockChannel.EXPECT().
					PublishWithContext(context.Background(), "", "name", false, false, amqp.Publishing{ContentType: "application/json", Body: data}).
					Return(nil)
			},
			hasError: false,
		},
		{
			name: "Cancel order update database failed",
			fn: func() {
				t.mockOrderDAO.EXPECT().
					Update(context.Background(), t.mockGormDB, order).
					Return(errors.New(""))
			},
			hasError: true,
		},
		{
			name: "Cancel order send message queue failed",
			fn: func() {
				t.mockOrderDAO.EXPECT().
					Update(context.Background(), t.mockGormDB, order).
					Return(nil)
				data, _ := json.Marshal(order)
				t.mockChannel.EXPECT().
					PublishWithContext(context.Background(), "", "name", false, false, amqp.Publishing{ContentType: "application/json", Body: data}).
					Return(errors.New(""))
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			err := t.svc.CancelOrder(context.Background(), 1)
			t.Equal(test.hasError, err != nil)
		})
	}
}
