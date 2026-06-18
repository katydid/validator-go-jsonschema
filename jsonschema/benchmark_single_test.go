// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonschema

import (
	"strings"
	"testing"
)

func TestBenchmarkSuiteSingle(t *testing.T) {
	filename := "lerna"
	suites, err := getBenchmarks()
	if err != nil {
		t.Fatal(err)
	}
	for _, suite := range suites {
		if filename != suite.name {
			continue
		}
		g, err := newGrammar(suite.schema)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("translated to: %v", g.String())
		matcher, err := Compile(suite.schema)
		if err != nil {
			t.Fatal(err)
		}
		want := !strings.Contains(suite.name, "-invalid")
		for i, data := range suite.datas {
			got, err := matcher.MatchBytes(data)
			if err != nil {
				t.Fatalf("at %d error: %v, given: %q", i, err, string(data))
			}
			if want != got {
				t.Fatalf("at %d want %v got %v, given: %q", i, want, got, string(data))
			}
		}
	}
}
