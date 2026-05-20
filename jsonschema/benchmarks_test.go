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

		instancesDatas := bytes.Split(instancesData, []byte("\n"))
		res = append(res, &benchsuite{name: entry.Name(), schema: schemaData, datas: instancesDatas})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("couldn't find benchmarks at %s", pathBenchmarks)
	}
	return res, nil
}

func TestBenchmarkSuite(t *testing.T) {
	notSupported := map[string]string{
		"ansible-meta":           "json: cannot unmarshal bool into Go struct field Schema.Object.properties of type schema.Schema",
		"cmake-presets":          "just takes long",
		"cql2":                   "could not find schema for #/$defs/andOrExpression",
		"cspell":                 "uniqueItems not supported",
		"deno":                   "uniqueItems not supported",
		"draft-04":               "json: cannot unmarshal bool into Go struct field Schema.Object.properties of type schema.Schema",
		"unreal-engine-uproject": "uniqueItems not supported",
		"geojson":                "timed out",
		"jsconfig":               "uniqueItems not supported",
		"krakend":                "uniqueItems not supported",
		"lazygit":                "uniqueItems are not supported",
		"openapi":                "could not find schema for #/$defs/server",
		"stylecop":               "uniqueItems are not supported",
		"ui5-manifest":           "json: cannot unmarshal bool into Go struct field Schema.definitions.Object.properties.Array.items of type []*schema.Schema",
	}
	suites, err := getBenchmarks()
	if err != nil {
		t.Fatal(err)
	}
	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			if reason, ok := notSupported[suite.name]; ok {
				t.Skipf("skipping, because %v", reason)
			}
			matcher, err := Compile(suite.schema)
			if err != nil {
				t.Fatal(err)
			}
			for _, data := range suite.datas {
				_, err := matcher.MatchBytes(data)
				if err != nil {
					t.Fatalf("error: %v, given: %q", err, string(data))
				}
			}
		})
	}
}
