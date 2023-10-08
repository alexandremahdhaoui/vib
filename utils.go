/*
Copyright 2023 Alexandre Mahdhaoui

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vib

import (
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"strings"
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

func JoinLine(buffer string, line string) string {
	if line == "" {
		return buffer
	}

	if buffer == "" {
		return line
	}

	return strings.Join([]string{buffer, line}, "\n")
}
