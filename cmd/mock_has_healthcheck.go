// Code generated by mockery v2.35.3. DO NOT EDIT.

package cmd

import (
	mock "github.com/stretchr/testify/mock"
	server "github.com/tscolari/servicetools/server"
)

// MockHasHealthcheck is an autogenerated mock type for the HasHealthcheck type
type MockHasHealthcheck struct {
	mock.Mock
}

// ConfigureHealthcheck provides a mock function with given fields: _a0
func (_m *MockHasHealthcheck) ConfigureHealthcheck(_a0 *server.WithHealthcheck) {
	_m.Called(_a0)
}

// NewMockHasHealthcheck creates a new instance of MockHasHealthcheck. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHasHealthcheck(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHasHealthcheck {
	mock := &MockHasHealthcheck{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
