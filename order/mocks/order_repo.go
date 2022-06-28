package mocks

import (
	"github.com/imylam/delivery-test/order"
	"github.com/stretchr/testify/mock"
)

// OrderRepository is a mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: order
func (_m *OrderRepository) Create(newOrder *order.Order) error {
	ret := _m.Called(newOrder)

	var r0 error
	if rf, ok := ret.Get(0).(func(*order.Order) error); ok {
		r0 = rf(newOrder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStatusByID provides a mock function with given fields: id
func (_m *OrderRepository) UpdateStatusByID(id int64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByID provides a mock function with given fields: id
func (_m *OrderRepository) FindByID(id int64) (*order.Order, error) {
	ret := _m.Called(id)

	var r0 *order.Order
	if rf, ok := ret.Get(0).(func(int64) *order.Order); ok {
		r0 = rf(id)
	} else {
		if _, ok := ret.Get(0).(*order.Order); ok {
			r0 = ret.Get(0).(*order.Order)
		} else {
			r0 = nil
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

// FindRange provides a mock function with given fields: limit, offset
func (_m *OrderRepository) FindRange(limit, offset int) (*[]order.Order, error) {
	ret := _m.Called(limit, offset)

	var r0 *[]order.Order
	if rf, ok := ret.Get(0).(func(int, int) *[]order.Order); ok {
		r0 = rf(limit, offset)
	} else {
		if _, ok := ret.Get(0).(*[]order.Order); ok {
			r0 = ret.Get(0).(*[]order.Order)
		} else {
			r0 = nil
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
