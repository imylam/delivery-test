package usecase

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/imylam/delivery-test/order"
	"github.com/imylam/delivery-test/order/infrastructure/googlemap"

	"github.com/imylam/delivery-test/order/mocks"
	"github.com/stretchr/testify/mock"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestPlaceOrder(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockMapClient := new(googlemap.MockMapClient)

	t.Run("success", func(t *testing.T) {
		distance := 888

		mockMapClient.On("GetDistance", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(distance, nil).Once()
		mockOrderRepo.On("Create", mock.AnythingOfType("*order.Order")).Return(nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		order, err := uc.PlaceOrder([]string{"22.300789", "114.167815"}, []string{"22.33540", "114.176155"})

		assert.Equal(t, true, err == nil)
		assert.Equal(t, distance, order.Distance)
		mockMapClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("map-api-error", func(t *testing.T) {
		mapErrMsg := "service unavailable"

		mockMapClient.On("GetDistance", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(0, errors.New(mapErrMsg)).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.PlaceOrder([]string{"22.300789", "114.167815"}, []string{"22.33540", "114.176155"})

		if err == nil {
			t.Errorf("TestPlaceOrder() fails, expect an error, got none")
			return
		}
		if err.Error() != mapErrMsg {
			t.Errorf("TestPlaceOrder() fails, expect error msg: %s, got: %s", mapErrMsg, err.Error())
		}
		mockMapClient.AssertExpectations(t)
	})

	t.Run("db-error", func(t *testing.T) {
		distance := 941

		mockMapClient.On("GetDistance", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(distance, nil).Once()
		mockOrderRepo.On("Create", mock.AnythingOfType("*order.Order")).Return(&mysql.MySQLError{}).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.PlaceOrder([]string{"22.300789", "114.167815"}, []string{"22.33540", "114.176155"})

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
		mockMapClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestTakeOrder(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockMapClient := new(googlemap.MockMapClient)

	mockOrder := order.Order{Status: order.StatusUnassigned}

	t.Run("success", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()
		mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		status, err := uc.TakeOrder(mockOrderID)

		assert.Equal(t, true, err == nil)
		assert.Equal(t, statusUpdateOrderStatusSuccess, status)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order-taken", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder
		tempOrder.Status = order.StatusTaken

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		assert.Equal(t, false, err == nil)
		assert.Equal(t, ErrorOrderTaken, err.Error())
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order-taken-when-update", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()
		mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(sql.ErrNoRows).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		assert.Equal(t, false, err == nil)
		assert.Equal(t, ErrorOrderTaken, err.Error())
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("no-such-order", func(t *testing.T) {
		mockOrderID := int64(1)

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(nil, sql.ErrNoRows).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		assert.Equal(t, false, err == nil)
		assert.Equal(t, sql.ErrNoRows, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("update-failure", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()
		mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(&mysql.MySQLError{}).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestListOrders(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockMapClient := new(googlemap.MockMapClient)

	mockPage := 1
	mockLimit := 4
	mockOrders := []order.Order{
		{ID: 1, Distance: 100, Status: order.StatusTaken},
		{ID: 2, Distance: 200, Status: order.StatusUnassigned},
		{ID: 3, Distance: 300, Status: order.StatusUnassigned},
		{ID: 4, Distance: 400, Status: order.StatusTaken},
	}

	t.Run("success", func(t *testing.T) {
		tempOrders := mockOrders

		mockOrderRepo.On("FindRange", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(&tempOrders, nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		orders, err := uc.ListOrders(mockPage, mockLimit)

		assert.Equal(t, true, err == nil)
		assert.Equal(t, len(tempOrders), len(*orders))
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("db-error", func(t *testing.T) {
		mockOrderRepo.On("FindRange", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(nil, &mysql.MySQLError{}).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.ListOrders(mockPage, mockLimit)

		assert.Equal(t, false, err == nil)

		_, isMysqlError := err.(*mysql.MySQLError)
		assert.Equal(t, true, isMysqlError)
		mockOrderRepo.AssertExpectations(t)
	})
}
