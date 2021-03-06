// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	types "github.com/kyma-incubator/hydroform/types"
	mock "github.com/stretchr/testify/mock"
)

// Builder is an autogenerated mock type for the Builder type
type Builder struct {
	mock.Mock
}

// CleanUp provides a mock function with given fields:
func (_m *Builder) CleanUp() {
	_m.Called()
}

// Create provides a mock function with given fields:
func (_m *Builder) Create() (*types.Cluster, *types.Provider, error) {
	ret := _m.Called()

	var r0 *types.Cluster
	if rf, ok := ret.Get(0).(func() *types.Cluster); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Cluster)
		}
	}

	var r1 *types.Provider
	if rf, ok := ret.Get(1).(func() *types.Provider); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*types.Provider)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
