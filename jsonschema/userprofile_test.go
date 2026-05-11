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

const SchemaJSONSchemaExampleUserProfile = `{
  "$id": "https://example.com/user-profile.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a user profile",
  "type": "object",
  "required": ["Username", "Email"],
  "properties": {
    "Username": {
      "type": "string"
    },
    "Email": {
      "type": "string",
      "format": "email"
    },
    "FullName": {
      "type": "string"
    },
    "Age": {
      "type": "integer",
      "minimum": 0
    },
    "Location": {
      "type": "string"
    },
    "Interests": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  }
}`

var userProfilePasses = []string{
	`{"Username":"OeLib","Email":"tzgwsgm@hotmail.com","FullName":"LSn3Er3","Age":3454797380135547817,"Location":"aSOafQQ"}`,
}

func TestUserProfileJSON(t *testing.T) {
	sch := SchemaJSONSchemaExampleUserProfile
	passes := userProfilePasses
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

func TestUserProfileReflect(t *testing.T) {
	sch := SchemaJSONSchemaExampleUserProfile
	passes := userProfilePasses
	var p parse.ParserWithInit = newReflectParser()

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
