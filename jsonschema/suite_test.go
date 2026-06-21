//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

func getFileNames(testPath string) []string {
	files := []string{}
	filepath.Walk(testPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == ".DS_Store" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if len(files) == 0 {
		panic("no test files found at " + testPath)
	}
	return files
}

type SchemaTest struct {
	Description string
	Schema      any
	Tests       []*SchemaTesty
}

type SchemaTesty struct {
	Description string
	Data        any
	Valid       bool
}

type Test struct {
	Filename    string
	Description string
	Schema      []byte
	Data        []byte
	Valid       bool
}

func (this Test) String() string {
	return this.Filename + ":" + this.Description
}

type Supported struct {
	// Files where all the tests must pass or the test actually fails.
	passingFiles  map[string]bool
	skippingFiles map[string]bool
	passingTests  map[string]bool
	skippingTests map[string]bool
}

func buildTests(t *testing.T, testPath string) []Test {
	tests := []Test{}
	filenames := getFileNames(testPath)
	if len(filenames) == 0 {
		t.Fatalf("expected test files, but found none")
	}
	for _, filename := range filenames {
		content, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		var schemaTests []*SchemaTest
		if err := std.UnmarshalJSON(content, &schemaTests); err != nil {
			panic(filename + ":" + err.Error())
		}
		for i := range schemaTests {
			schemaStr, err := json.Marshal(schemaTests[i].Schema)
			if err != nil {
				panic(err)
			}
			schemaDesc := schemaTests[i].Description
			for j := range schemaTests[i].Tests {
				dataStr, err := json.Marshal(schemaTests[i].Tests[j].Data)
				if err != nil {
					panic(err)
				}
				if len(filename) > len(testPath) {
					filename = filename[len(testPath):]
				}
				tests = append(tests, Test{
					Filename:    filename,
					Description: schemaDesc + ":" + schemaTests[i].Tests[j].Description,
					Schema:      schemaStr,
					Data:        dataStr,
					Valid:       schemaTests[i].Tests[j].Valid,
				})
			}
		}
	}
	if len(tests) == 0 {
		t.Fatalf("expected tests, but found none")
	}
	return tests
}

func runTests(t *testing.T, testPath string, supported *Supported, opts ...Option) {
	tests := buildTests(t, testPath)
	t.Logf("total number of tests: %d", len(tests))

	checkFilesExists(supported.passingFiles, tests)
	checkFilesExists(supported.skippingFiles, tests)
	checkTestsExists(supported.passingTests, tests)
	checkTestsExists(supported.skippingTests, tests)

	passed := 0
	skippedTests := 0
	failedTests := 0

	for _, test := range tests {
		if supported.skippingFiles[test.Filename] {
			t.Logf("skip: %v", test)
			skippedTests++
			continue
		}
		if supported.skippingTests[test.String()] {
			t.Logf("skip: %v", test)
			skippedTests++
			continue
		}
		t.Logf("## RUN: %v", test)
		valid, err := MatchBytes(test.Schema, test.Data, opts...)
		if err != nil || valid != test.Valid {
			if supported.passingFiles[test.Filename] || supported.passingTests[test.String()] {
				if err != nil {
					t.Errorf("UNEXPECTED ERROR: %v: Interpret error %v", test, err)
				} else {
					t.Errorf("UNEXPECTED FAILURE: %v: expected %v got %v", test, test.Valid, valid)
				}
			} else {
				if err != nil {
					t.Logf("ERROR: %v: Interpret error %v", test, err)
				} else {
					t.Logf("TODO: %v: expected %v got %v", test, test.Valid, valid)
				}
			}
			failedTests++
		} else {
			t.Logf("PASSED: %v", test)
			passed++
		}
	}
	t.Logf("number of tests passing: %d, skippedTests: %d, failedTests: %d", passed, skippedTests, failedTests)
}

// check that files specified in the skip/pass sets actually exist.
func checkFilesExists(spec map[string]bool, tests []Test) {
	for name := range spec {
		found := false
		for _, test := range tests {
			if test.Filename == name {
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("given %v file not found %s", tests, name))
		}
	}
}

func checkTestsExists(spec map[string]bool, tests []Test) {
	for name := range spec {
		found := false
		for _, test := range tests {
			if test.String() == name {
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("test not found %s", name))
		}
	}
}
