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

type DealTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	db         *sql.DB
	mockDB     sqlmock.Sqlmock
	mockGormDB *gorm.DB
}

func (t *DealTestSuite) SetupTest() {
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

func (t *DealTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.db.Close()
}

func TestDealTestSuite(t *testing.T) {
	suite.Run(t, new(DealTestSuite))
}

func (t *DealTestSuite) TestInsert() {
	tests := []struct {
		name     string
		deals    []*models.Deal
		fn       func()
		hasError bool
	}{
		{
			name: "Insert deals success",
			deals: []*models.Deal{
				{ID: 1},
				{ID: 2},
			},
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `deal` (`taker_order_id`,`maker_order_id`,`quantity`,`price`,`id`) VALUES (?,?,?,?,?),(?,?,?,?,?)")).
					WillReturnResult(sqlmock.NewResult(2, 2))
				t.mockDB.ExpectCommit()
			},
			hasError: false,
		},
		{
			name:     "Insert deals no deal",
			deals:    nil,
			fn:       func() {},
			hasError: false,
		},
		{
			name: "Insert deals failed",
			deals: []*models.Deal{
				{ID: 1},
				{ID: 2},
			},
			fn: func() {
				t.mockDB.ExpectBegin()
				t.mockDB.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `deal` (`taker_order_id`,`maker_order_id`,`quantity`,`price`,`id`) VALUES (?,?,?,?,?),(?,?,?,?,?)")).
					WillReturnError(errors.New(""))
				t.mockDB.ExpectCommit()
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			err := NewDeal().Insert(context.Background(), t.mockGormDB, test.deals)
			t.Equal(test.hasError, err != nil)
		})
	}
}

func (t *DealTestSuite) TestList() {
	tests := []struct {
		name     string
		fn       func()
		expected []*models.Deal
		hasError bool
	}{
		{
			name: "List deals success",
			fn: func() {
				t.mockDB.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `deal`")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "taker_order_id", "maker_order_id", "quantity", "price"}).
						AddRow(1, 2, 3, 4, 5))
			},
			expected: []*models.Deal{
				{
					ID:           1,
					TakerOrderID: 2,
					MakerOrderID: 3,
					Quantity:     4,
					Price:        5,
				},
			},
			hasError: false,
		},
		{
			name: "List deals failed",
			fn: func() {
				t.mockDB.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `deal`")).
					WillReturnError(errors.New(""))
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn()
			actual, err := NewDeal().List(context.Background(), t.mockGormDB, nil)
			t.Equal(test.hasError, err != nil)
			t.Equal(test.expected, actual)
		})
	}
}
