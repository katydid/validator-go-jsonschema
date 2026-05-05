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
	"testing"
)

// Files where all the tests must pass or the test actually fails.
var passingFile = map[string]bool{
	"minLength.json": true,
	"maxLength.json": true,
	"date-time.json": true,
	"email.json":     true,
	"hostname.json":  true,
	"ipv4.json":      true,
	"ipv6.json":      true,
	"uri.json":       true,
	"unknonw.json":   true,
}

var skippingFile = map[string]bool{
	"uniqueItems.json": true, // We do not support uniqueItems, see https://github.com/katydid/validator-go-jsonschema/blob/main/decisions/uniqueItems.md
}

var skippingTest = map[string]bool{
	"ecmascript-regex.json:patterns always use unicode semantics with pattern:ascii character in json string":                      true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with patternProperties:ascii character in json string":            true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with pattern:literal unicode character in json string":            true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with patternProperties:literal unicode character in json string":  true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with pattern:unicode character in hex format in string":           true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with patternProperties:unicode character in hex format in string": true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with pattern:unicode matching is case-sensitive":                  true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patterns always use unicode semantics with patternProperties:unicode matching is case-sensitive":        true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:pattern with non-ASCII digits:ascii digits":                                                             true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patternProperties with non-ASCII digits:ascii digits":                                                   true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:pattern with non-ASCII digits:ascii non-digits":                                                         true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patternProperties with non-ASCII digits:ascii non-digits":                                               true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:pattern with non-ASCII digits:non-ascii digits (BENGALI DIGIT FOUR, BENGALI DIGIT TWO)":                 true, // https://github.com/dlclark/regexp2/issues/101
	"ecmascript-regex.json:patternProperties with non-ASCII digits:non-ascii digits (BENGALI DIGIT FOUR, BENGALI DIGIT TWO)":       true, // https://github.com/dlclark/regexp2/issues/101
}

func TestSuiteDraft4(t *testing.T) {
	tests := buildTests(t)
	t.Logf("skipping files: %d", len(skippingFile))
	t.Logf("total number of tests: %d", len(tests))
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
			if passingFile[test.Filename] {
				if err != nil {
					t.Errorf("--- FAIL: %v: Interpret error %v", test, err)
				} else {
					t.Errorf("--- FAIL: %v: expected %v got %v", test, test.Valid, valid)
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
