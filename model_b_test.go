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

package casbin

import (
	"fmt"
	"testing"

	"github.com/ApicaSystem/casbin/v2/util"
)

func rawEnforce(sub string, obj string, act string) bool {
	policy := [2][3]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}}
	for _, rule := range policy {
		if sub == rule[0] && obj == rule[1] && act == rule[2] {
			return true
		}
	}
	return false
}

func BenchmarkRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawEnforce("alice", "data1", "read")
	}
}

func BenchmarkBasicModel(b *testing.B) {
	e, _ := NewEnforcer("examples/basic_model.conf", "examples/basic_policy.csv", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "data1", "read")
	}
}

func BenchmarkRBACModel(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "data2", "read")
	}
}

func BenchmarkRBACModelSizes(b *testing.B) {
	cases := []struct {
		name      string
		roles     int
		resources int
		users     int
	}{
		{name: "small", roles: 100, resources: 10, users: 1000},
		{name: "medium", roles: 1000, resources: 100, users: 10000},
		{name: "large", roles: 10000, resources: 1000, users: 100000},
	}
	for _, c := range cases {
		c := c

		e, err := NewEnforcer("examples/rbac_model.conf", false)
		if err != nil {
			b.Fatal(err)
		}

		pPolicies := make([][]string, c.roles)
		for i := range pPolicies {
			pPolicies[i] = []string{
				fmt.Sprintf("group-has-a-very-long-name-%d", i),
				fmt.Sprintf("data-has-a-very-long-name-%d", i%c.resources),
				"read",
			}
		}
		if _, err := e.AddPolicies(pPolicies); err != nil {
			b.Fatal(err)
		}

		gPolicies := make([][]string, c.users)
		for i := range gPolicies {
			gPolicies[i] = []string{
				fmt.Sprintf("user-has-a-very-long-name-%d", i),
				fmt.Sprintf("group-has-a-very-long-name-%d", i%c.roles),
			}
		}
		if _, err := e.AddGroupingPolicies(gPolicies); err != nil {
			b.Fatal(err)
		}

		// Set up enforcements, alternating between things a user can access
		// and things they cannot. Use 17 tests so that we get a variety of users
		// and roles rather than always landing on a multiple of 2/10/whatever.
		enforcements := make([][]interface{}, 17)
		for i := range enforcements {
			userNum := (c.users / len(enforcements)) * i
			roleNum := userNum % c.roles
			resourceNum := roleNum % c.resources
			if i%2 == 0 {
				resourceNum += 1
				resourceNum %= c.resources
			}
			enforcements[i] = []interface{}{
				fmt.Sprintf("user-has-a-very-long-name-%d", userNum),
				fmt.Sprintf("data-has-a-very-long-name-%d", resourceNum),
				"read",
			}
		}

		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = e.Enforce(enforcements[i%len(enforcements)]...)
			}
		})
	}
}

func BenchmarkRBACModelSmall(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_model.conf", false)

	// 100 roles, 10 resources.
	for i := 0; i < 100; i++ {
		_, err := e.AddPolicy(fmt.Sprintf("group%d", i), fmt.Sprintf("data%d", i/10), "read")
		if err != nil {
			b.Fatal(err)
		}
	}

	// 1000 users.
	for i := 0; i < 1000; i++ {
		_, err := e.AddGroupingPolicy(fmt.Sprintf("user%d", i), fmt.Sprintf("group%d", i/10))
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("user501", "data9", "read")
	}
}

func BenchmarkRBACModelMedium(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_model.conf", false)

	// 1000 roles, 100 resources.
	pPolicies := make([][]string, 0)
	for i := 0; i < 1000; i++ {
		pPolicies = append(pPolicies, []string{fmt.Sprintf("group%d", i), fmt.Sprintf("data%d", i/10), "read"})
	}

	_, err := e.AddPolicies(pPolicies)
	if err != nil {
		b.Fatal(err)
	}

	// 10000 users.
	gPolicies := make([][]string, 0)
	for i := 0; i < 10000; i++ {
		gPolicies = append(gPolicies, []string{fmt.Sprintf("user%d", i), fmt.Sprintf("group%d", i/10)})
	}

	_, err = e.AddGroupingPolicies(gPolicies)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("user5001", "data99", "read")
	}
}

func BenchmarkRBACModelLarge(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_model.conf", false)

	// 10000 roles, 1000 resources.
	pPolicies := make([][]string, 0)
	for i := 0; i < 10000; i++ {
		pPolicies = append(pPolicies, []string{fmt.Sprintf("group%d", i), fmt.Sprintf("data%d", i/10), "read"})
	}

	_, err := e.AddPolicies(pPolicies)
	if err != nil {
		b.Fatal(err)
	}

	// 100000 users.
	gPolicies := make([][]string, 0)
	for i := 0; i < 100000; i++ {
		gPolicies = append(gPolicies, []string{fmt.Sprintf("user%d", i), fmt.Sprintf("group%d", i/10)})
	}

	_, err = e.AddGroupingPolicies(gPolicies)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("user50001", "data999", "read")
	}
}

func BenchmarkRBACModelWithResourceRoles(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_with_resource_roles_model.conf", "examples/rbac_with_resource_roles_policy.csv", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "data1", "read")
	}
}

func BenchmarkRBACModelWithDomains(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_with_domains_model.conf", "examples/rbac_with_domains_policy.csv", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "domain1", "data1", "read")
	}
}

func BenchmarkABACModel(b *testing.B) {
	e, _ := NewEnforcer("examples/abac_model.conf", false)
	data1 := newTestResource("data1", "alice")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", data1, "read")
	}
}

func BenchmarkABACRuleModel(b *testing.B) {
	e, _ := NewEnforcer("examples/abac_rule_model.conf", false)
	sub := newTestSubject("alice", 18)

	for i := 0; i < 1000; i++ {
		_, _ = e.AddPolicy("r.sub.Age > 20", fmt.Sprintf("data%d", i), "read")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce(sub, "data100", "read")
	}
}

func BenchmarkKeyMatchModel(b *testing.B) {
	e, _ := NewEnforcer("examples/keymatch_model.conf", "examples/keymatch_policy.csv", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "/alice_data/resource1", "GET")
	}
}

func BenchmarkRBACModelWithDeny(b *testing.B) {
	e, _ := NewEnforcer("examples/rbac_with_deny_model.conf", "examples/rbac_with_deny_policy.csv")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "data1", "read")
	}
}

func BenchmarkPriorityModel(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model.conf", "examples/priority_policy.csv")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("alice", "data1", "read")
	}
}

func BenchmarkRBACModelWithDomainPatternLarge(b *testing.B) {
	e, _ := NewEnforcer("examples/performance/rbac_with_pattern_large_scale_model.conf", "examples/performance/rbac_with_pattern_large_scale_policy.csv")
	e.AddNamedDomainMatchingFunc("g", "", util.KeyMatch4)
	_ = e.BuildRoleLinks()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Enforce("staffUser1001", "/orgs/1/sites/site001", "App001.Module001.Action1001")
	}
}
