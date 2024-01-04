// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StorageService is an autogenerated mock type for the StorageService type
type StorageService struct {
	mock.Mock
}

// ReadFromBucketObject provides a mock function with given fields: objectId
func (_m *StorageService) ReadFromBucketObject(objectId string) ([]byte, error) {
	ret := _m.Called(objectId)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(objectId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(objectId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}