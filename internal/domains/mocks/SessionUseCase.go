// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SessionUseCase is an autogenerated mock type for the SessionUseCase type
type SessionUseCase struct {
	mock.Mock
}

// CreateIfNotExists provides a mock function with given fields: session
func (_m *SessionUseCase) CreateIfNotExists(session string) int {
	ret := _m.Called(session)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

type mockConstructorTestingTNewSessionUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewSessionUseCase creates a new instance of SessionUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSessionUseCase(t mockConstructorTestingTNewSessionUseCase) *SessionUseCase {
	mock := &SessionUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
