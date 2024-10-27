// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// PostServiceI is an autogenerated mock type for the PostServiceI type
type PostServiceI struct {
	mock.Mock
}

// IsPrivate provides a mock function with given fields: postId
func (_m *PostServiceI) IsPrivate(postId uuid.UUID) (bool, error) {
	ret := _m.Called(postId)

	if len(ret) == 0 {
		panic("no return value specified for IsPrivate")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (bool, error)); ok {
		return rf(postId)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(postId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPostServiceI creates a new instance of PostServiceI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostServiceI(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostServiceI {
	mock := &PostServiceI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
