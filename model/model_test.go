// Copyright 2019 The casbin Authors. All Rights Reserved.
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

package model

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ApicaSystem/casbin/v2/config"
	"github.com/ApicaSystem/casbin/v2/constant"
)

var (
	basicExample = filepath.Join("..", "examples", "basic_model.conf")
	basicConfig  = &MockConfig{
		data: map[string]string{
			"request_definition::r": "sub, obj, act",
			"policy_definition::p":  "sub, obj, act",
			"policy_effect::e":      "some(where (p.eft == allow))",
			"matchers::m":           "r.sub == p.sub && r.obj == p.obj && r.act == p.act",
		},
	}
)

type MockConfig struct {
	data map[string]string
	config.ConfigInterface
}

func (mc *MockConfig) String(key string) string {
	return mc.data[key]
}

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m == nil {
		t.Error("new model should not be nil")
	}
}

func TestNewModelFromFile(t *testing.T) {
	m, err := NewModelFromFile(basicExample)
	if err != nil {
		t.Errorf("model failed to load from file: %s", err)
	}
	if m == nil {
		t.Error("model should not be nil")
	}
}

func TestNewModelFromString(t *testing.T) {
	modelBytes, _ := ioutil.ReadFile(basicExample)
	modelString := string(modelBytes)
	m, err := NewModelFromString(modelString)
	if err != nil {
		t.Errorf("model failed to load from string: %s", err)
	}
	if m == nil {
		t.Error("model should not be nil")
	}
}

func TestLoadModelFromConfig(t *testing.T) {
	m := NewModel()
	err := m.loadModelFromConfig(basicConfig)
	if err != nil {
		t.Error("basic config should not return an error")
	}
	m = NewModel()
	err = m.loadModelFromConfig(&MockConfig{})
	if err == nil {
		t.Error("empty config should return error")
	} else {
		// check for missing sections in message
		for _, rs := range requiredSections {
			if !strings.Contains(err.Error(), sectionNameMap[rs]) {
				t.Errorf("section name: %s should be in message", sectionNameMap[rs])
			}
		}
	}
}

func TestHasSection(t *testing.T) {
	m := NewModel()
	_ = m.loadModelFromConfig(basicConfig)
	for _, sec := range requiredSections {
		if !m.hasSection(sec) {
			t.Errorf("%s section was expected in model", sec)
		}
	}
	m = NewModel()
	_ = m.loadModelFromConfig(&MockConfig{})
	for _, sec := range requiredSections {
		if m.hasSection(sec) {
			t.Errorf("%s section was not expected in model", sec)
		}
	}
}

func TestModel_AddDef(t *testing.T) {
	m := NewModel()
	s := "r"
	v := "sub, obj, act"
	ok := m.AddDef(s, s, v)
	if !ok {
		t.Errorf("non empty assertion should be added")
	}
	ok = m.AddDef(s, s, "")
	if ok {
		t.Errorf("empty assertion value should not be added")
	}
}

func TestModelToTest(t *testing.T) {
	testModelToText(t, "r.sub == p.sub && r.obj == p.obj && r_func(r.act, p.act) && testr_func(r.act, p.act)", "r_sub == p_sub && r_obj == p_obj && r_func(r_act, p_act) && testr_func(r_act, p_act)")
	testModelToText(t, "r.sub == p.sub && r.obj == p.obj && p_func(r.act, p.act) && testp_func(r.act, p.act)", "r_sub == p_sub && r_obj == p_obj && p_func(r_act, p_act) && testp_func(r_act, p_act)")
}

func testModelToText(t *testing.T, mData, mExpected string) {
	m := NewModel()
	data := map[string]string{
		"r": "sub, obj, act",
		"p": "sub, obj, act",
		"e": "some(where (p.eft == allow))",
		"m": mData,
	}
	expected := map[string]string{
		"r": "sub, obj, act",
		"p": "sub, obj, act",
		"e": constant.AllowOverrideEffect,
		"m": mExpected,
	}
	addData := func(ptype string) {
		m.AddDef(ptype, ptype, data[ptype])
	}
	for ptype := range data {
		addData(ptype)
	}
	newM := NewModel()
	print(m.ToText())
	_ = newM.LoadModelFromText(m.ToText())
	for ptype := range data {
		if newM[ptype][ptype].Value != expected[ptype] {
			t.Errorf("\"%s\" assertion value changed, current value: %s, it should be: %s", ptype, newM[ptype][ptype].Value, expected[ptype])
		}
	}
}
