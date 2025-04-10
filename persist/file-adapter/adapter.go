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
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/ApicaSystem/casbin/v2/model"
	"github.com/ApicaSystem/casbin/v2/persist"
	"github.com/ApicaSystem/casbin/v2/util"
)

// Adapter is the file adapter for Casbin.
// It can load policy from file or save policy to file.
type Adapter struct {
	filePath string
}

func (a *Adapter) UpdatePolicy(sec string, ptype string, oldRule, newRule []string) error {
	return errors.New("not implemented")
}

func (a *Adapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	return errors.New("not implemented")
}

func (a *Adapter) UpdateFilteredPolicies(sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return nil, errors.New("not implemented")
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(filePath string) *Adapter {
	return &Adapter{filePath: filePath}
}

// LoadPolicy loads all policy rules from the storage.
func (a *Adapter) LoadPolicy(model model.Model) error {
	if a.filePath == "" {
		return errors.New("invalid file path, file path cannot be empty")
	}

	return a.loadPolicyFile(model, persist.LoadPolicyLine)
}

// SavePolicy saves all policy rules to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	if a.filePath == "" {
		return errors.New("invalid file path, file path cannot be empty")
	}

	var tmp bytes.Buffer

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}

	return a.savePolicyFile(strings.TrimRight(tmp.String(), "\n"))
}

func (a *Adapter) loadPolicyFile(model model.Model, handler func(string, model.Model) error) error {
	f, err := os.Open(a.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		err = handler(line, model)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (a *Adapter) savePolicyFile(text string) error {
	f, err := os.Create(a.filePath)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	_, err = w.WriteString(text)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return f.Close()
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// AddPolicies adds policy rules to the storage.
func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicies removes policy rules from the storage.
func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
