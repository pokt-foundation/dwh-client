// Code generated by mockery v2.32.3. DO NOT EDIT.

package client

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockIDWHClient is an autogenerated mock type for the IDWHClient type
type MockIDWHClient struct {
	mock.Mock
}

// GetTotalRelaysForAccountIDs provides a mock function with given fields: ctx, params
func (_m *MockIDWHClient) GetTotalRelaysForAccountIDs(ctx context.Context, params GetTotalRelaysForAccountIDsParams) ([]AnalyticsRelaysTotal, error) {
	ret := _m.Called(ctx, params)

	var r0 []AnalyticsRelaysTotal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, GetTotalRelaysForAccountIDsParams) ([]AnalyticsRelaysTotal, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, GetTotalRelaysForAccountIDsParams) []AnalyticsRelaysTotal); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]AnalyticsRelaysTotal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, GetTotalRelaysForAccountIDsParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalRelaysForPortalAppIDs provides a mock function with given fields: ctx, params
func (_m *MockIDWHClient) GetTotalRelaysForPortalAppIDs(ctx context.Context, params GetTotalRelaysForPortalAppIDsParams) ([]AnalyticsRelaysTotal, error) {
	ret := _m.Called(ctx, params)

	var r0 []AnalyticsRelaysTotal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, GetTotalRelaysForPortalAppIDsParams) ([]AnalyticsRelaysTotal, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, GetTotalRelaysForPortalAppIDsParams) []AnalyticsRelaysTotal); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]AnalyticsRelaysTotal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, GetTotalRelaysForPortalAppIDsParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockIDWHClient creates a new instance of MockIDWHClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIDWHClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIDWHClient {
	mock := &MockIDWHClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
