// Code generated by mockery v2.35.3. DO NOT EDIT.

package cmd

import (
	context "context"
	slog "log/slog"

	mock "github.com/stretchr/testify/mock"
)

// MockServer is an autogenerated mock type for the Server type
type MockServer struct {
	mock.Mock
}

// Start provides a mock function with given fields: _a0, _a1
func (_m *MockServer) Start(_a0 context.Context, _a1 *slog.Logger) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *slog.Logger) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields: _a0, _a1
func (_m *MockServer) Stop(_a0 context.Context, _a1 *slog.Logger) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *slog.Logger) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockServer creates a new instance of MockServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockServer {
	mock := &MockServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
