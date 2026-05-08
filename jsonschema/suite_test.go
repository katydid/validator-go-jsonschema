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
	"fmt"
	"testing"
)

// Files where all the tests must pass or the test actually fails.
var passingFile = map[string]bool{
	// "additionalItems.json": true,
	"additionalProperties.json": true,
	// "allOf.json": true,
	// "anyOf.json":   true,
	"default.json": true,
	// "dependencies.json": true,
	"enum.json":   true,
	"format.json": true,
	// "infinite-loop-detection.json": true,
	// "items.json":             true,
	"maximum.json":           true,
	"maxItems.json":          true,
	"maxLength.json":         true,
	"maxProperties.json":     true,
	"minimum.json":           true,
	"minItems.json":          true,
	"minLength.json":         true,
	"minProperties.json":     true,
	"multipleOf.json":        true,
	"not.json":               true,
	"oneOf.json":             true,
	"pattern.json":           true,
	"patternProperties.json": true,
	"properties.json":        true,
	// "ref.json": true,
	"required.json": true,
	"type.json":     true,

	// optional/format
	"date-time.json": true,
	"email.json":     true,
	"hostname.json":  true,
	"ipv4.json":      true,
	"ipv6.json":      true,
	"unknown.json":   true,
	"uri.json":       true,

	// optional
	"ecmascript-regex.json":     true,
	"non-bmp-regex.json":        true,
	"zeroTerminatedFloats.json": true,
}

var skippingFile = map[string]bool{
	"uniqueItems.json": true, // We do not support uniqueItems, see https://github.com/katydid/validator-go-jsonschema/blob/main/decisions/uniqueItems.md
	"refRemote.json":   true, // remote and file ref support should be relatively easy to add, but is just not of theoretical importance at this stage.
	"definitions.json": true, // remote and file ref support should be relatively easy to add, but is just not of theoretical importance at this stage.
	// optional
	"id.json":             true, // remote and file ref support should be relatively easy to add, but is just not of theoretical importance at this stage.
	"bignum.json":         true, // Need better decimal support in at least maximum, integer, number, minimum
	"float-overflow.json": true, // Need better checking for float overflow to convert to decimal in the json parser and we need to support decimal in multipleOf
}

var passingTest = map[string]bool{}

var skippingTest = map[string]bool{
	// optional/format
	"uri.json:validation of URIs:unescaped non US-ASCII characters": true, // need a better URI library
	"uri.json:validation of URIs:invalid backslash character":       true, // need a better URI library
	"uri.json:validation of URIs:invalid \" character":              true, // need a better URI library
	"uri.json:validation of URIs:invalid <> characters":             true, // need a better URI library
	"uri.json:validation of URIs:invalid {} characters":             true, // need a better URI library
	"uri.json:validation of URIs:invalid ^ character":               true, // need a better URI library
	"uri.json:validation of URIs:invalid ` character":               true, // need a better URI library
	"uri.json:validation of URIs:invalid SPACE character":           true, // need a better URI library
	"uri.json:validation of URIs:invalid | character":               true, // need a better URI library
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
			panic(fmt.Sprintf("file not found %s", name))
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

func TestSuiteDraft4(t *testing.T) {
	tests := buildTests(t)
	t.Logf("skipping files: %d", len(skippingFile))
	t.Logf("total number of tests: %d", len(tests))

	checkFilesExists(passingFile, tests)
	checkFilesExists(skippingFile, tests)
	checkTestsExists(skippingTest, tests)
	checkTestsExists(passingTest, tests)

	passed := 0
	skippedTests := 0
	failedTests := 0

	for _, test := range tests {
		if skippingFile[test.Filename] {
			t.Logf("--- SKIP: %v", test)
			skippedTests++
			continue
		}
		if skippingTest[test.String()] {
			t.Logf("--- SKIP: %v", test)
			skippedTests++
			continue
		}
		t.Logf("--- RUN: %v", test)
		valid, err := Validate(test.Schema, test.Data)
		if err != nil || valid != test.Valid {
			if passingFile[test.Filename] || passingTest[test.String()] {
				if err != nil {
					t.Errorf("UNEXPECTED FAILURE: %v: Interpret error %v", test, err)
				} else {
					t.Errorf("UNEXPECTED FAILURE: %v: expected %v got %v", test, test.Valid, valid)
				}
			} else {
				if err != nil {
					t.Logf("--- FAIL: %v: Interpret error %v", test, err)
				} else {
					t.Logf("--- FAIL: %v: expected %v got %v", test, test.Valid, valid)
				}
			}
			failedTests++
		} else {
			t.Logf("--- PASS: %v", test)
			passed++
		}
	}
	t.Logf("number of tests passing: %d, skippedTests: %d, failedTests: %d", passed, skippedTests, failedTests)
}
