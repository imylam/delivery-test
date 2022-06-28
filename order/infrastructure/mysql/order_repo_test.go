package mysql

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/imylam/delivery-test/order"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type AnyInt struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyInt) Match(v driver.Value) bool {
	_, ok := v.(int)
	return ok
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected error '%s' when opening sqlmock database connection", err.Error())
		return
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	qInsert := "INSERT INTO orders"
	qSelect := "SELECT (.+) FROM orders"

	mockOrder := order.Order{
		Distance: 1000,
		Status:   order.StatusUnassigned,
	}

	t.Run("success", func(t *testing.T) {
		tempOrder := mockOrder
		mockOrderID := int64(8)

		prepInsert := mock.ExpectPrepare(qInsert)
		prepInsert.ExpectExec().WithArgs(tempOrder.Distance, tempOrder.Status).
			WillReturnResult(sqlmock.NewResult(mockOrderID, 1))

		rows := sqlmock.NewRows([]string{"id", "distance", "status", "created_at", "updated_at"}).
			AddRow(mockOrderID, tempOrder.Distance, tempOrder.Status, time.Now(), time.Now())
		mock.ExpectQuery(qSelect).WithArgs(mockOrderID).WillReturnRows(rows)

		repo := NewOrderRepositoryMysql(sqlxDB)
		err = repo.Create(&tempOrder)

		assert.Equal(t, true, err == nil)
		assert.Equal(t, mockOrderID, tempOrder.ID)
	})

	t.Run("insert-error", func(t *testing.T) {
		tempOrder := mockOrder

		prepInsert := mock.ExpectPrepare(qInsert)
		prepInsert.ExpectExec().WithArgs(tempOrder.Distance, tempOrder.Status).
			WillReturnError(&mysql.MySQLError{})

		repo := NewOrderRepositoryMysql(sqlxDB)
		err = repo.Create(&tempOrder)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
	})

	t.Run("select-error", func(t *testing.T) {
		tempOrder := mockOrder
		mockOrderID := int64(10)

		prepInsert := mock.ExpectPrepare(qInsert)
		prepInsert.ExpectExec().WithArgs(tempOrder.Distance, tempOrder.Status).
			WillReturnResult(sqlmock.NewResult(mockOrderID, 1))

		mock.ExpectQuery(qSelect).WithArgs(mockOrderID).WillReturnError(&mysql.MySQLError{})

		repo := NewOrderRepositoryMysql(sqlxDB)
		err = repo.Create(&tempOrder)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
	})
}

func TestUpdateStatusByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected error '%s' when opening sqlmock database connection", err.Error())
		return
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	q := "UPDATE orders SET"

	t.Run("success", func(t *testing.T) {
		mockOrderID := int64(8)

		prepUpdate := mock.ExpectPrepare(q)
		prepUpdate.ExpectExec().WithArgs(order.StatusTaken, mockOrderID, order.StatusUnassigned).
			WillReturnResult(sqlmock.NewResult(mockOrderID, 1))

		repo := NewOrderRepositoryMysql(sqlxDB)
		err = repo.UpdateStatusByID(mockOrderID)

		assert.Equal(t, true, err == nil)
	})

	t.Run("no-update", func(t *testing.T) {
		mockOrderID := int64(8)

		prepUpdate := mock.ExpectPrepare(q)
		prepUpdate.ExpectExec().WithArgs(order.StatusTaken, mockOrderID, order.StatusUnassigned).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewOrderRepositoryMysql(sqlxDB)
		err = repo.UpdateStatusByID(mockOrderID)

		assert.Equal(t, false, err == nil)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("update-error", func(t *testing.T) {
		mockOrderID := int64(8)

		prepUpdate := mock.ExpectPrepare(q)
		prepUpdate.ExpectExec().WithArgs(order.StatusTaken, mockOrderID, order.StatusUnassigned).
			WillReturnError(&mysql.MySQLError{})

		repo := NewOrderRepositoryMysql(sqlxDB)
		err = repo.UpdateStatusByID(mockOrderID)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
	})
}

func TestFindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected error '%s' when opening sqlmock database connection", err.Error())
		return
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	q := "SELECT (.+) FROM orders"

	t.Run("success", func(t *testing.T) {
		mockOrderID := int64(8)
		mockDistance := 888
		mockStatus := order.StatusUnassigned

		rows := sqlmock.NewRows([]string{"id", "distance", "status", "created_at", "updated_at"}).
			AddRow(mockOrderID, mockDistance, mockStatus, time.Now(), time.Now())
		mock.ExpectQuery(q).WithArgs(mockOrderID).WillReturnRows(rows)

		repo := NewOrderRepositoryMysql(sqlxDB)
		order, err := repo.FindByID(mockOrderID)

		assert.Equal(t, true, err == nil)
		assert.Equal(t, mockOrderID, order.ID)
	})

	t.Run("select-error", func(t *testing.T) {
		mockOrderID := int64(99)

		mock.ExpectQuery(q).WithArgs(mockOrderID).WillReturnError(&mysql.MySQLError{})

		repo := NewOrderRepositoryMysql(sqlxDB)
		_, err := repo.FindByID(mockOrderID)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
	})
}

func TestFindRange(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected error '%s' when opening sqlmock database connection", err.Error())
		return
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	q := "SELECT (.+) FROM orders"

	t.Run("success", func(t *testing.T) {
		mockLimit := 3
		mockPage := 1

		rows := sqlmock.NewRows([]string{"id", "distance", "status", "created_at", "updated_at"}).
			AddRow(1, 100, order.StatusTaken, time.Now(), time.Now()).
			AddRow(2, 200, order.StatusUnassigned, time.Now(), time.Now()).
			AddRow(3, 300, order.StatusUnassigned, time.Now(), time.Now())
		mock.ExpectQuery(q).WithArgs(mockLimit, mockPage).WillReturnRows(rows)

		repo := NewOrderRepositoryMysql(sqlxDB)
		orders, err := repo.FindRange(mockLimit, mockPage)

		assert.Equal(t, true, err == nil)
		assert.Equal(t, 3, len(*orders))
	})

	t.Run("select-error", func(t *testing.T) {
		mockLimit := 4
		mockPage := 2

		mock.ExpectQuery(q).WithArgs(mockLimit, mockPage).WillReturnError(&mysql.MySQLError{})

		repo := NewOrderRepositoryMysql(sqlxDB)
		_, err := repo.FindRange(mockLimit, mockPage)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
	})
}
