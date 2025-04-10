// Copyright 2018 The casbin Authors. All Rights Reserved.
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

package casbin

import (
	"testing"

	fileadapter "github.com/ApicaSystem/casbin/v2/persist/file-adapter"
	"github.com/ApicaSystem/casbin/v2/util"
)

func TestInitFilteredAdapter(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

	// policy should not be loaded yet
	testHasPolicy(t, e, []string{"admin", "domain1", "data1", "read"}, false)
}

func TestLoadFilteredPolicy(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)
	if err := e.LoadPolicy(); err != nil {
		t.Errorf("unexpected error in LoadPolicy: %v", err)
	}

	// validate initial conditions
	testHasPolicy(t, e, []string{"admin", "domain1", "data1", "read"}, true)
	testHasPolicy(t, e, []string{"admin", "domain2", "data2", "read"}, true)

	if err := e.LoadFilteredPolicy(&fileadapter.Filter{
		P: []string{"", "domain1"},
		G: []string{"", "", "domain1"},
	}); err != nil {
		t.Errorf("unexpected error in LoadFilteredPolicy: %v", err)
	}
	if !e.IsFiltered() {
		t.Errorf("adapter did not set the filtered flag correctly")
	}

	// only policies for domain1 should be loaded
	testHasPolicy(t, e, []string{"admin", "domain1", "data1", "read"}, true)
	testHasPolicy(t, e, []string{"admin", "domain2", "data2", "read"}, false)

	if err := e.SavePolicy(); err == nil {
		t.Errorf("enforcer did not prevent saving filtered policy")
	}
	if err := e.GetAdapter().SavePolicy(e.GetModel()); err == nil {
		t.Errorf("adapter did not prevent saving filtered policy")
	}
}

func TestLoadMoreTypeFilteredPolicy(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_pattern_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_pattern_model.conf", adapter)
	if err := e.LoadPolicy(); err != nil {
		t.Errorf("unexpected error in LoadPolicy: %v", err)
	}
	e.AddNamedMatchingFunc("g2", "matching func", util.KeyMatch2)
	_ = e.BuildRoleLinks()

	testEnforce(t, e, "alice", "/book/1", "GET", true)

	// validate initial conditions
	testHasPolicy(t, e, []string{"book_admin", "book_group", "GET"}, true)
	testHasPolicy(t, e, []string{"pen_admin", "pen_group", "GET"}, true)

	if err := e.LoadFilteredPolicy(&fileadapter.Filter{
		P:  []string{"book_admin"},
		G:  []string{"alice"},
		G2: []string{"", "book_group"},
	}); err != nil {
		t.Errorf("unexpected error in LoadFilteredPolicy: %v", err)
	}
	if !e.IsFiltered() {
		t.Errorf("adapter did not set the filtered flag correctly")
	}

	testHasPolicy(t, e, []string{"alice", "/pen/1", "GET"}, false)
	testHasPolicy(t, e, []string{"alice", "/pen2/1", "GET"}, false)
	testHasPolicy(t, e, []string{"pen_admin", "pen_group", "GET"}, false)
	testHasGroupingPolicy(t, e, []string{"alice", "book_admin"}, true)
	testHasGroupingPolicy(t, e, []string{"bob", "pen_admin"}, false)
	testHasGroupingPolicy(t, e, []string{"cathy", "pen_admin"}, false)
	testHasGroupingPolicy(t, e, []string{"cathy", "/book/1/2/3/4/5"}, false)

	testEnforce(t, e, "alice", "/book/1", "GET", true)
	testEnforce(t, e, "alice", "/pen/1", "GET", false)
}

func TestAppendFilteredPolicy(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)
	if err := e.LoadPolicy(); err != nil {
		t.Errorf("unexpected error in LoadPolicy: %v", err)
	}

	// validate initial conditions
	testHasPolicy(t, e, []string{"admin", "domain1", "data1", "read"}, true)
	testHasPolicy(t, e, []string{"admin", "domain2", "data2", "read"}, true)

	if err := e.LoadFilteredPolicy(&fileadapter.Filter{
		P: []string{"", "domain1"},
		G: []string{"", "", "domain1"},
	}); err != nil {
		t.Errorf("unexpected error in LoadFilteredPolicy: %v", err)
	}
	if !e.IsFiltered() {
		t.Errorf("adapter did not set the filtered flag correctly")
	}

	// only policies for domain1 should be loaded
	testHasPolicy(t, e, []string{"admin", "domain1", "data1", "read"}, true)
	testHasPolicy(t, e, []string{"admin", "domain2", "data2", "read"}, false)

	// disable clear policy and load second domain
	if err := e.LoadIncrementalFilteredPolicy(&fileadapter.Filter{
		P: []string{"", "domain2"},
		G: []string{"", "", "domain2"},
	}); err != nil {
		t.Errorf("unexpected error in LoadFilteredPolicy: %v", err)
	}

	// both domain policies should be loaded
	testHasPolicy(t, e, []string{"admin", "domain1", "data1", "read"}, true)
	testHasPolicy(t, e, []string{"admin", "domain2", "data2", "read"}, true)
}

func TestFilteredPolicyInvalidFilter(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

	if err := e.LoadFilteredPolicy([]string{"", "domain1"}); err == nil {
		t.Errorf("expected error in LoadFilteredPolicy, but got nil")
	}
}

func TestFilteredPolicyEmptyFilter(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

	if err := e.LoadFilteredPolicy(nil); err != nil {
		t.Errorf("unexpected error in LoadFilteredPolicy: %v", err)
	}
	if e.IsFiltered() {
		t.Errorf("adapter did not reset the filtered flag correctly")
	}
	if err := e.SavePolicy(); err != nil {
		t.Errorf("unexpected error in SavePolicy: %v", err)
	}
}

func TestUnsupportedFilteredPolicy(t *testing.T) {
	e, _ := NewEnforcer("examples/rbac_with_domains_model.conf", "examples/rbac_with_domains_policy.csv")

	err := e.LoadFilteredPolicy(&fileadapter.Filter{
		P: []string{"", "domain1"},
		G: []string{"", "", "domain1"},
	})
	if err == nil {
		t.Errorf("encorcer should have reported incompatibility error")
	}
}

func TestFilteredAdapterEmptyFilepath(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

	if err := e.LoadFilteredPolicy(nil); err != nil {
		t.Errorf("unexpected error in LoadFilteredPolicy: %v", err)
	}
}

func TestFilteredAdapterInvalidFilepath(t *testing.T) {
	e, _ := NewEnforcer()

	adapter := fileadapter.NewFilteredAdapter("examples/does_not_exist_policy.csv")
	_ = e.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

	if err := e.LoadFilteredPolicy(nil); err == nil {
		t.Errorf("expected error in LoadFilteredPolicy, but got nil")
	}
}
