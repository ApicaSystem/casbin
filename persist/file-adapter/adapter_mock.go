// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileadapter

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ApicaSystem/casbin/v2/model"
	"github.com/ApicaSystem/casbin/v2/persist"
)

// AdapterMock is the file adapter for Casbin.
// It can load policy from file or save policy to file.
type AdapterMock struct {
	filePath   string
	errorValue string
}

// NewAdapterMock is the constructor for AdapterMock.
func NewAdapterMock(filePath string) *AdapterMock {
	a := AdapterMock{}
	a.filePath = filePath
	return &a
}

// LoadPolicy loads all policy rules from the storage.
func (a *AdapterMock) LoadPolicy(model model.Model) error {
	err := a.loadPolicyFile(model, persist.LoadPolicyLine)
	return err
}

// SavePolicy saves all policy rules to the storage.
func (a *AdapterMock) SavePolicy(model model.Model) error {
	return nil
}

func (a *AdapterMock) loadPolicyFile(model model.Model, handler func(string, model.Model) error) error {
	f, err := os.Open(a.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err2 := handler(line, model); err2 != nil {
			return err2
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// SetMockErr sets string to be returned by of the mock during testing.
func (a *AdapterMock) SetMockErr(errorToSet string) {
	a.errorValue = errorToSet
}

// GetMockErr returns a mock error or nil.
func (a *AdapterMock) GetMockErr() error {
	var returnError error
	if a.errorValue != "" {
		returnError = errors.New(a.errorValue)
	}
	return returnError
}

// AddPolicy adds a policy rule to the storage.
func (a *AdapterMock) AddPolicy(sec string, ptype string, rule []string) error {
	return a.GetMockErr()
}

// AddPolicies removes policy rules from the storage.
func (a *AdapterMock) AddPolicies(sec string, ptype string, rules [][]string) error {
	return a.GetMockErr()
}

// RemovePolicy removes a policy rule from the storage.
func (a *AdapterMock) RemovePolicy(sec string, ptype string, rule []string) error {
	return a.GetMockErr()
}

// RemovePolicies removes policy rules from the storage.
func (a *AdapterMock) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return a.GetMockErr()
}

// UpdatePolicy removes a policy rule from the storage.
func (a *AdapterMock) UpdatePolicy(sec string, ptype string, oldRule, newPolicy []string) error {
	return a.GetMockErr()
}

func (a *AdapterMock) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	return a.GetMockErr()
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *AdapterMock) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return a.GetMockErr()
}
