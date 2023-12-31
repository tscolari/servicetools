// Code generated by mockery v2.35.3. DO NOT EDIT.

package server

import (
	mock "github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

// MockGRPCRegisterFunc is an autogenerated mock type for the GRPCRegisterFunc type
type MockGRPCRegisterFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *MockGRPCRegisterFunc) Execute(_a0 *grpc.Server) {
	_m.Called(_a0)
}

// NewMockGRPCRegisterFunc creates a new instance of MockGRPCRegisterFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGRPCRegisterFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGRPCRegisterFunc {
	mock := &MockGRPCRegisterFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
