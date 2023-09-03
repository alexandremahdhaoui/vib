package vib

import (
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

func ToPointer[T any](t T) *T {
	value := t

	return &value
}

func removeIndexFromSliceFast[T any](sl []T, i int) []T {
	sl[i] = sl[len(sl)-1]
	return sl[:len(sl)-1]
}

func Debug(v interface{}) {
	logger.SugaredLoggerWithCallerSkip(1).Debugf("%#v", v)
}
