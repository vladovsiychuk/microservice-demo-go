// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	post "github.com/vladovsiychuk/microservice-demo-go/internal/post"
)

// PostI is an autogenerated mock type for the PostI type
type PostI struct {
	mock.Mock
}

// Update provides a mock function with given fields: req
func (_m *PostI) Update(req post.PostRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(post.PostRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPostI creates a new instance of PostI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostI(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostI {
	mock := &PostI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}