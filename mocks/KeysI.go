// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// KeysI is an autogenerated mock type for the KeysI type
type KeysI struct {
	mock.Mock
}

// NewKeysI creates a new instance of KeysI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKeysI(t interface {
	mock.TestingT
	Cleanup(func())
}) *KeysI {
	mock := &KeysI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
