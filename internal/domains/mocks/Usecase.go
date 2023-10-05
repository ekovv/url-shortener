// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// CheckConn provides a mock function with given fields:
func (_m *UseCase) CheckConn() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLong provides a mock function with given fields: shortURL
func (_m *UseCase) GetLong(shortURL string) (string, error) {
	ret := _m.Called(shortURL)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(shortURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(shortURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShort provides a mock function with given fields: path
func (_m *UseCase) GetShort(path string) (string, error) {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveLog provides a mock function with given fields: id, path
func (_m *UseCase) SaveWithoutGenerate(id string, path string) (string, error) {
	ret := _m.Called(id, path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(id, path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUseCase(t mockConstructorTestingTNewUseCase) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
