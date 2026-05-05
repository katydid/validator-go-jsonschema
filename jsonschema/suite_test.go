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
	"strings"
	"testing"

	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parse/debug"
)

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
		if err != nil {
			t.Logf("--- FAIL: %v: Interpret error %v", test, err)
			failedTests++
		} else if valid != test.Valid {
			t.Logf("--- FAIL: %v: expected %v got %v", test, test.Valid, valid)
			failedTests++
		} else {
			t.Logf("--- PASS: %v", test)
			passed++
		}
	}
	t.Logf("number of tests passing: %d, skippedTests: %d, failedTests: %d", passed, skippedTests, failedTests)
}

func testDebug(t *testing.T, test Test) {
	g, err := newGrammar(test.Schema)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())

	jsonp := json.NewJSONSchemaParser()
	p := debug.NewLogger(jsonp, debug.NewLineLogger())
	p.Init(test.Data)
	valid, err := ValidateParser(test.Schema, p)
	if err != nil {
		t.Fatalf("Interpret error %v", err)
	} else if valid != test.Valid {
		t.Fatalf("expected %v got %v", test.Valid, valid)
	}
}

func TestDebug(t *testing.T) {
	tests := buildTests(t)
	for _, test := range tests {
		if !strings.Contains(test.String(), "properties.json:object properties validation:doesn't invalidate other properties") {
			continue
		}
		testDebug(t, test)
	}
}
