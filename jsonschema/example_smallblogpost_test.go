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
	"testing"

	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parse"
)

const SchemaSmallBlogPostExample = `{ "title": "small jsonschema for a blogpost",
  "type":"object", "additionalProperties":false, "required": ["content"],
  "properties": {
    "content": { "type":"string" },
    "author": { "$ref":"#/definitions/user-profile" } },
  "definitions": { "user-profile": {
    "type": "object", "additionalProperties":false, "required": ["username"], 
    "properties": {
      "username": { "type":"string" },
      "email": { "type":"string", "format":"email" } } } } }`

func TestSchemaSmallBlogPostExample(t *testing.T) {
	sch := SchemaSmallBlogPostExample
	passes := []string{
		`{"content": "Dragons"}`,
		`{"content": "Dragons", "author": {"username": "Khaleesi"}}`,
	}
	var p parse.ParserWithInit = json.NewJSONSchemaParser()

	g, err := newGrammar([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())
	for _, input := range passes {
		p.Init([]byte(input))
		m, err := MatchParser([]byte(sch), p)
		if err != nil {
			t.Fatal(err)
		}
		if !m {
			t.Errorf("expected true, but got no match for %s", input)
		}
	}
}
