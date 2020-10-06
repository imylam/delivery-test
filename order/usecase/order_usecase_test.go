package usecase

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/imylam/delivery-test/domain"
	"github.com/imylam/delivery-test/googlemap"

	"github.com/imylam/delivery-test/domain/mocks"
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
		mockOrderRepo.On("Create", mock.AnythingOfType("*domain.Order")).Return(nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		order, err := uc.PlaceOrder([]string{"22.300789", "114.167815"}, []string{"22.33540", "114.176155"})

		if err != nil {
			t.Errorf("TestPlaceOrder() fails, expect no error, got: %s", err.Error())
		}
		if order.Distance != distance {
			t.Errorf("TestPlaceOrder() fails, expect order.Distance: %d, got: %d", order.Distance, distance)
		}
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
		mockOrderRepo.On("Create", mock.AnythingOfType("*domain.Order")).Return(&mysql.MySQLError{}).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.PlaceOrder([]string{"22.300789", "114.167815"}, []string{"22.33540", "114.176155"})

		if err == nil {
			t.Errorf("TestPlaceOrder() fails, expect no error, got: %s", err.Error())
		}
		if _, ok := err.(*mysql.MySQLError); !ok {
			t.Errorf("TestPlaceOrder() fails: Expected mysql error")
		}

		mockMapClient.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestTakeOrder(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockMapClient := new(googlemap.MockMapClient)

	mockOrder := domain.Order{Status: domain.StatusUnassigned}

	t.Run("success", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()
		mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		status, err := uc.TakeOrder(mockOrderID)

		if err != nil {
			t.Errorf("TestTakeOrder() fails, expect no error, got: %s", err.Error())
		}
		if status != statusUpdateOrderStatusSuccess {
			t.Errorf("TestTakeOrder() fails, expect status: %s, got: %s", statusUpdateOrderStatusSuccess, status)
		}

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order-taken", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder
		tempOrder.Status = domain.StatusTaken

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		if err == nil {
			t.Errorf("TestTakeOrder() fails, expect an error, got none")
			return
		}
		if err.Error() != ErrorOrderTaken {
			t.Errorf("TestTakeOrder() fails, expect error msg: %s, got:%s", ErrorOrderTaken, err.Error())
		}

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order-taken-when-update", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()
		mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(sql.ErrNoRows).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		if err == nil {
			t.Errorf("TestTakeOrder() fails, expect an error, got none")
			return
		}
		if err.Error() != ErrorOrderTaken {
			t.Errorf("TestTakeOrder() fails, expect error msg: %s, got:%s", ErrorOrderTaken, err.Error())
		}

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("no-such-order", func(t *testing.T) {
		mockOrderID := int64(1)

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(nil, sql.ErrNoRows).Once()
		// mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		if err == nil {
			t.Errorf("TestTakeOrder() fails, expect an error, got none")
			return
		}
		if err != sql.ErrNoRows {
			t.Errorf("TestTakeOrder() fails, expect sql.ErrNoRows, got:%s", err.Error())
		}

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("update-failure", func(t *testing.T) {
		mockOrderID := int64(1)
		tempOrder := mockOrder

		mockOrderRepo.On("FindByID", mock.AnythingOfType("int64")).Return(&tempOrder, nil).Once()
		mockOrderRepo.On("UpdateStatusByID", mock.AnythingOfType("int64")).Return(&mysql.MySQLError{}).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.TakeOrder(mockOrderID)

		if err == nil {
			t.Errorf("TestTakeOrder() fails, expect an error, got none")
			return
		}
		if _, ok := err.(*mysql.MySQLError); !ok {
			t.Errorf("TestTakeOrder() fails: Expected mysql error")
		}

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestListOrders(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)
	mockMapClient := new(googlemap.MockMapClient)

	mockPage := 1
	mockLimit := 4
	mockOrders := []domain.Order{
		domain.Order{ID: 1, Distance: 100, Status: domain.StatusTaken},
		domain.Order{ID: 2, Distance: 200, Status: domain.StatusUnassigned},
		domain.Order{ID: 3, Distance: 300, Status: domain.StatusUnassigned},
		domain.Order{ID: 4, Distance: 400, Status: domain.StatusTaken},
	}

	t.Run("success", func(t *testing.T) {
		tempOrders := mockOrders

		mockOrderRepo.On("FindRange", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(&tempOrders, nil).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		orders, err := uc.ListOrders(mockPage, mockLimit)

		if err != nil {
			t.Errorf("TestListOrders() fails, expect no error, got: %s", err.Error())
		}
		if len(*orders) != len(tempOrders) {
			t.Errorf("TestListOrders() fails, expect number of orders: %d, got: %d",
				len(tempOrders), len(*orders))
		}

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("db-error", func(t *testing.T) {
		mockOrderRepo.On("FindRange", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(nil, &mysql.MySQLError{}).Once()

		uc := NewOrderUsecase(mockOrderRepo, mockMapClient)
		_, err := uc.ListOrders(mockPage, mockLimit)

		if err == nil {
			t.Errorf("TestListOrders() fails, expect an error, got none")
			return
		}
		if _, ok := err.(*mysql.MySQLError); !ok {
			t.Errorf("TestListOrders() fails: Expected mysql error")
		}

		mockOrderRepo.AssertExpectations(t)
	})
}
