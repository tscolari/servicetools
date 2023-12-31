// Code generated by mockery v2.35.3. DO NOT EDIT.

package cmd

import (
	mock "github.com/stretchr/testify/mock"
	server "github.com/tscolari/servicetools/server"
)

// MockHasGRPC is an autogenerated mock type for the HasGRPC type
type MockHasGRPC struct {
	mock.Mock
}

// ConfigureGRPC provides a mock function with given fields: _a0
func (_m *MockHasGRPC) ConfigureGRPC(_a0 *server.WithGRPC) {
	_m.Called(_a0)
}

// NewMockHasGRPC creates a new instance of MockHasGRPC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHasGRPC(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHasGRPC {
	mock := &MockHasGRPC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
