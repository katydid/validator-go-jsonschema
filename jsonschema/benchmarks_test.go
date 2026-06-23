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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const pathBenchmarks = "../../validator-jsonschema-benchmarks/schemas/"

type benchsuite struct {
	name   string
	schema []byte
	datas  [][]byte
}

func getBenchmarks() ([]*benchsuite, error) {
	res := []*benchsuite{}
	entries, err := os.ReadDir(pathBenchmarks)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(pathBenchmarks, entry.Name())
		schemaName := filepath.Join(path, "schema.json")
		schemaData, err := os.ReadFile(schemaName)
		if err != nil {
			return nil, err
		}

		instancesName := filepath.Join(path, "instances.jsonl")
		instancesData, err := os.ReadFile(instancesName)
		if err != nil {
			return nil, err
		}

		lines := bytes.Split(instancesData, []byte("\n"))
		if len(lines[len(lines)-1]) == 0 {
			lines = lines[:len(lines)-1]
		}
		res = append(res, &benchsuite{name: entry.Name(), schema: schemaData, datas: lines})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("couldn't find benchmarks at %s", pathBenchmarks)
	}
	return res, nil
}

func TestBenchmarkSuite(t *testing.T) {
	notSupported := map[string]string{
		"ajv-cosmicrealms-invalid": "uniqueItems not supported",
		"ajv-cosmicrealms-valid":   "uniqueItems not supported",
		"cspell":                   "uniqueItems not supported",
		"cql2":                     "dynamicRef not supported",
		"deno":                     "uniqueItems not supported",
		"draft-04":                 "uniqueItems not supported",
		"jsconfig":                 "uniqueItems not supported",
		"krakend":                  "uniqueItems not supported",
		"lazygit":                  "uniqueItems not supported",
		"openapi":                  "dynamicRef not supported",
		"stylecop":                 "uniqueItems not supported",
		"ui5-manifest":             "uniqueItems not supported",
		"unreal-engine-uproject":   "uniqueItems not supported",
		"zschema-basic-invalid":    "uniqueItems not supported",
		"zschema-basic-valid":      "uniqueItems not supported",
		"zschema-advanced-invalid": "uniqueItems not supported",
		"zschema-advanced-valid":   "uniqueItems not supported",
	}
	unsupportedByOthers := map[string]string{
		"krakend-rmUniqueItems":      "not supported by ajv, ajv-bun, hyperjump, networknt. Features: ref with weird syntax: #/definitions/https%3A~1~1www.krakend.io~1schema~1v2.7~1timeunits.json/definitions/timeunit, default, const, patternProperties, title, pattern, not, anyOf, deprecated",
		"ui5-manifest-rmUniqueItems": "not supported by ajv, ajv-bun, boon, go-kaptinlin, go-santhosh-tekuri, hyperjump, networknt",
		"cspell-rmUniqueItems":       "not supported by boon, go-kaptinlin, go-santhosh-tekuri, json_schemer and kmp",
	}
	skippingBecauseSlow := map[string]string{
		"geojson-invalid": "just slow to compile",
		"geojson":         "just slow to compile",
	}
	notMatchingYet := map[string]string{}
	suites, err := getBenchmarks()
	if err != nil {
		t.Fatal(err)
	}
	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			if reason, ok := notSupported[suite.name]; ok {
				t.Skipf("skipping unsupported, because %v", reason)
			}
			if reason, ok := unsupportedByOthers[suite.name]; ok {
				t.Skipf("skipping unsupported by others, because %v", reason)
			}
			if reason, ok := skippingBecauseSlow[suite.name]; ok {
				t.Skipf("skipping temporarily, because %v", reason)
			}
			matcher, err := Compile(suite.schema)
			if err != nil {
				t.Fatal(err)
			}
			if reason, ok := notMatchingYet[suite.name]; ok {
				t.Skipf("temporary skipping matching, because %v", reason)
			}
			want := !strings.Contains(suite.name, "-invalid")
			for _, data := range suite.datas {
				got, err := matcher.MatchBytes(data)
				if err != nil {
					t.Fatalf("error: %v, given: %q", err, string(data))
				}
				if want != got {
					t.Errorf("want %v got %v, given: %q", want, got, string(data))
				}
			}
		})
	}
}
