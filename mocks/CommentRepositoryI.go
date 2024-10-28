// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	comment "github.com/vladovsiychuk/microservice-demo-go/internal/comment"

	uuid "github.com/google/uuid"
)

// CommentRepositoryI is an autogenerated mock type for the CommentRepositoryI type
type CommentRepositoryI struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *CommentRepositoryI) Create(_a0 *comment.Comment) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*comment.Comment) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByKey provides a mock function with given fields: _a0, commentId
func (_m *CommentRepositoryI) FindByKey(_a0 *comment.Comment, commentId uuid.UUID) error {
	ret := _m.Called(_a0, commentId)

	if len(ret) == 0 {
		panic("no return value specified for FindByKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*comment.Comment, uuid.UUID) error); ok {
		r0 = rf(_a0, commentId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0
func (_m *CommentRepositoryI) Update(_a0 *comment.Comment) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*comment.Comment) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCommentRepositoryI creates a new instance of CommentRepositoryI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommentRepositoryI(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentRepositoryI {
	mock := &CommentRepositoryI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}