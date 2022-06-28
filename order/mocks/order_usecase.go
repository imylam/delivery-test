package mocks

import (
	"github.com/imylam/delivery-test/order"
	"github.com/stretchr/testify/mock"
)

// OrderUsecase is a mock type for the OrderUsecase type
type OrderUsecase struct {
	mock.Mock
}

// PlaceOrder provides a mock function with given fields: origins, destinations
func (_m *OrderUsecase) PlaceOrder(origins, destinations []string) (*order.Order, error) {
	ret := _m.Called(origins, destinations)

	var r0 *order.Order
	if rf, ok := ret.Get(0).(func([]string, []string) *order.Order); ok {
		r0 = rf(origins, destinations)
	} else {
		if _, ok := ret.Get(0).(*order.Order); ok {
			r0 = ret.Get(0).(*order.Order)
		} else {
			r0 = nil
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, []string) error); ok {
		r1 = rf(origins, destinations)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TakeOrder provides a mock function with given fields: id
func (_m *OrderUsecase) TakeOrder(id int64) (string, error) {
	ret := _m.Called(id)

	var r0 string
	if rf, ok := ret.Get(0).(func(int64) string); ok {
		r0 = rf(id)
	} else {
		if _, ok := ret.Get(0).(string); ok {
			r0 = ret.Get(0).(string)
		} else {
			r0 = ""
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOrders provides a mock function with given fields: page, limit
func (_m *OrderUsecase) ListOrders(page, limit int) (*[]order.Order, error) {
	ret := _m.Called(page, limit)

	var r0 *[]order.Order
	if rf, ok := ret.Get(0).(func(int, int) *[]order.Order); ok {
		r0 = rf(page, limit)
	} else {
		if _, ok := ret.Get(0).(*[]order.Order); ok {
			r0 = ret.Get(0).(*[]order.Order)
		} else {
			r0 = nil
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
