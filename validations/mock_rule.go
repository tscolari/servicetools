// Code generated by mockery v2.35.3. DO NOT EDIT.

package validations

import mock "github.com/stretchr/testify/mock"

// MockRule is an autogenerated mock type for the Rule type
type MockRule struct {
	mock.Mock
}

// Validate provides a mock function with given fields: value
func (_m *MockRule) Validate(value interface{}) error {
	ret := _m.Called(value)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRule creates a new instance of MockRule. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRule(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRule {
	mock := &MockRule{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}