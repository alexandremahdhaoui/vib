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

package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

type (
	Kind interface {
		APIVersion() types.APIVersion
		Kind() types.Kind
	}
)

func NewAPIKind(apiVersion types.APIVersion, kind types.Kind) APIKind {
	return &concreteAPIKind{
		apiVersion: apiVersion,
		kind:       kind,
	}
}

type concreteAPIKind struct {
	apiVersion types.APIVersion
	kind       types.Kind
}

func (c *concreteAPIKind) APIVersion() types.APIVersion {
	return c.apiVersion
}

func (c *concreteAPIKind) Kind() types.Kind {
	return c.kind
}
