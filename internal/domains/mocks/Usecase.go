// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	storage "url-shortener/internal/storage"

	mock "github.com/stretchr/testify/mock"
)

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

// GetAllUrls provides a mock function with given fields: user
func (_m *UseCase) GetAllUrls(user string) ([]storage.URL, error) {
	ret := _m.Called(user)

	var r0 []storage.URL
	if rf, ok := ret.Get(0).(func(string) []storage.URL); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]storage.URL)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLong provides a mock function with given fields: user, shortURL
func (_m *UseCase) GetLong(user string, shortURL string) (string, error) {
	ret := _m.Called(user, shortURL)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(user, shortURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(user, shortURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShort provides a mock function with given fields: user, path
func (_m *UseCase) GetShort(user string, path string) (string, error) {
	ret := _m.Called(user, path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(user, path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(user, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveAndGetSessionMap provides a mock function with given fields: session
func (_m *UseCase) SaveAndGetSessionMap(session string) int {
	ret := _m.Called(session)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// SaveWithoutGenerate provides a mock function with given fields: user, id, path
func (_m *UseCase) SaveWithoutGenerate(user string, id string, path string) (string, error) {
	ret := _m.Called(user, id, path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(user, id, path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(user, id, path)
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
