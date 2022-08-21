package dao

import (
	"context"
	"database/sql"
	"dealer/internal/models"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OrderTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	db         *sql.DB
	mockDB     sqlmock.Sqlmock
	mockGormDB *gorm.DB
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
}

func (t *OrderTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.db.Close()
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (t *OrderTestSuite) TestInsert() {
	tests := []struct {
		name     string
		order    *models.Order
		fn       func()
		hasError bool
	}{
		{
			name: "Insert order success",
			order: &models.Order{
				OrderType:      models.OrderTypeBuy,
				Quantity:       5,
				RemainQuantity: 5,
				PriceType:      models.PriceTypeLimit,
				Price:          10,
			},
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `order` (`order_type`,`quantity`,`remain_quantity`,`price_type`,`price`,`is_cancel`) VALUES (?,?,?,?,?,?)")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name:     "Insert order nil",
			order:    nil,
			fn:       func() {},
			hasError: false,
		},
		{
			name: "Insert order failed",
			order: &models.Order{
				OrderType:      models.OrderTypeBuy,
				Quantity:       5,
				RemainQuantity: 5,
				PriceType:      models.PriceTypeLimit,
				Price:          10,
			},
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `order` (`order_type`,`quantity`,`remain_quantity`,`price_type`,`price`,`is_cancel`) VALUES (?,?,?,?,?,?)")).
					WillReturnError(errors.New(""))
				t.mockDB.ExpectCommit()
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			err := NewOrder().Insert(context.Background(), t.mockGormDB, test.order)
			t.Equal(test.hasError, err != nil)
		})
	}
}

func (t *OrderTestSuite) TestUpdate() {
	tests := []struct {
		name     string
		fn       func()
		hasError bool
	}{
		{
			name: "Update order success",
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("UPDATE `order` SET `is_cancel`=? WHERE `id` = ?")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name: "Update order failed",
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("UPDATE `order` SET `is_cancel`=? WHERE `id` = ?")).
					WillReturnError(errors.New(""))
				t.mockDB.ExpectCommit()
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			err := NewOrder().Update(context.Background(), t.mockGormDB, &models.Order{ID: 1, IsCancel: true})
			t.Equal(test.hasError, err != nil)
		})
	}
}

func (t *OrderTestSuite) TestBulkUpdate() {
	tests := []struct {
		name     string
		fn       func()
		hasError bool
	}{
		{
			name: "Bulk update order success",
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `order` (`order_type`,`quantity`,`remain_quantity`,`price_type`,`price`,`is_cancel`,`id`) VALUES (?,?,?,?,?,?,?),(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `remain_quantity`=VALUES(`remain_quantity`)")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name: "Bulk update order failed",
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `order` (`order_type`,`quantity`,`remain_quantity`,`price_type`,`price`,`is_cancel`,`id`) VALUES (?,?,?,?,?,?,?),(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `remain_quantity`=VALUES(`remain_quantity`)")).
					WillReturnError(errors.New(""))
				t.mockDB.ExpectCommit()
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			orders := []*models.Order{
				{
					ID:             1,
					RemainQuantity: 3,
				},
				{
					ID:             2,
					RemainQuantity: 4,
				},
			}
			err := NewOrder().BulkUpdate(context.Background(), t.mockGormDB, orders)
			t.Equal(test.hasError, err != nil)
		})
	}
}
