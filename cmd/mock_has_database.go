// Code generated by mockery v2.35.3. DO NOT EDIT.

package cmd

import (
	mock "github.com/stretchr/testify/mock"
	server "github.com/tscolari/servicetools/server"
)

// MockHasDatabase is an autogenerated mock type for the HasDatabase type
type MockHasDatabase struct {
	mock.Mock
}

// ConfigureDatabase provides a mock function with given fields: _a0
func (_m *MockHasDatabase) ConfigureDatabase(_a0 *server.WithDB) {
	_m.Called(_a0)
}

// NewMockHasDatabase creates a new instance of MockHasDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHasDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHasDatabase {
	mock := &MockHasDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
