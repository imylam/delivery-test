package googlemap

import (
	"github.com/stretchr/testify/mock"
)

// MockMapClient is a mock type for the MapClient  type
type MockMapClient struct {
	mock.Mock
}

// GetDistance provides a mock function with given fields: origin, destination
func (_m *MockMapClient) GetDistance(origin string, destination string) (distance int, err error) {
	ret := _m.Called(origin, destination)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, string) int); ok {
		r0 = rf(origin, destination)
	} else {
		if _, ok := ret.Get(0).(int); ok {
			r0 = ret.Get(0).(int)
		} else {
			r0 = 0
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(origin, destination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
