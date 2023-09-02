package vib

import (
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"os"
	"path/filepath"
)

type Raw interface{}

func mkBaseDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func fileExist(path string) (bool, error) {
	err := mkBaseDir(path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			logger.Error(err)
			return false, err
		}
	}

	return true, nil
}

func ToPointer[T any](t T) *T {
	value := t

	return &value
}

func removeIndexFromSliceFast[T any](sl []T, i int) []T {
	sl[i] = sl[len(sl)-1]
	return sl[:len(sl)-1]
}

func Debug(v interface{}) {
	logger.Debug("%#v", v)
}
