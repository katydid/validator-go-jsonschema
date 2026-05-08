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

	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parse/debug"
	"github.com/katydid/validator-go/validator/intern"
)

func TestDebug(t *testing.T) {
	tests := buildTests(t)
	for _, test := range tests {
		if !strings.Contains(test.String(), "items.json:an array of schemas for items:wrong types") {
			continue
		}
		testDebug(t, test)
		return
	}
	t.Fatal("test not found")
}

func testDebug(t *testing.T, test Test) {
	g, err := newGrammar(test.Schema)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("test.Data: %s", test.Data)
	t.Logf("translated to: %v", g.String())
	// gStr := `
	// 	tag(object):{
	// 		(bar:*&{(bar:->integer())*;(!(bar):*)*})
	// 		;(foo:*&{(foo:->stringType())*;(!(foo):*)*})
	// 	}`
	// g, err = parser.NewParser().ParseGrammar(gStr)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	jsonp := json.NewJSONSchemaParser()
	p := debug.NewLogger(jsonp, debug.NewLineLogger())
	p.Init(test.Data)

	valid, err := intern.Interpret(g, true, p)
	if err != nil {
		t.Fatalf("Interpret error %v", err)
	} else if valid != test.Valid {
		t.Fatalf("expected %v got %v", test.Valid, valid)
	}
}
