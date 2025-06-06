// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	user "github.com/miqxzz/miqxzzforum/auth_service/internal/proto"
)

// UserServiceClient is an autogenerated mock type for the UserServiceClient type
type UserServiceClient struct {
	mock.Mock
}

// GetUsername provides a mock function with given fields: ctx, in, opts
func (_m *UserServiceClient) GetUsername(ctx context.Context, in *user.UserRequest, opts ...grpc.CallOption) (*user.UserResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetUsername")
	}

	var r0 *user.UserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.UserRequest, ...grpc.CallOption) (*user.UserResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *user.UserRequest, ...grpc.CallOption) *user.UserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.UserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *user.UserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserServiceClient creates a new instance of UserServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserServiceClient {
	mock := &UserServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
