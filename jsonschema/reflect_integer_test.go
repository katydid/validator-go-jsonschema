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
	goreflect "reflect"
	"testing"

	"github.com/katydid/parser-go-reflect/reflect"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

func TestReflectInteger(t *testing.T) {
	sch := `{
		"properties": {
			"age": {
				"type": "integer"
			}
		}
	}`
	g, err := newGrammar([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())
	p := reflect.NewJSONSchemaParser()
	input := `
	{
	  "age":123
	}`
	var v any
	if err := std.UnmarshalJSON([]byte(input), &v); err != nil {
		t.Fatal(err)
	}
	p.Init(goreflect.ValueOf(v))
	m, err := MatchParser([]byte(sch), p)
	if err != nil {
		t.Fatal(err)
	}
	if !m {
		t.Errorf("expected true, but got no match for %s", input)
	}
}
