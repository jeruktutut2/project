package helpers

import "time"

type TimeHelper interface {
	NowUnixMili() int64
}

type TimeHelperImplementation struct {
}

func NewTimeHelper() TimeHelper {
	return &TimeHelperImplementation{}
}

func (helper *TimeHelperImplementation) NowUnixMili() int64 {
	return time.Now().UnixMilli()
}
