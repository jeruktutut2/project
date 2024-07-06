package mockhelpers

import "github.com/stretchr/testify/mock"

type TimeHelperMock struct {
	Mock mock.Mock
}

func (helper *TimeHelperMock) NowUnixMili() int64 {
	arguments := helper.Mock.Called()
	return arguments.Get(0).(int64)
}
