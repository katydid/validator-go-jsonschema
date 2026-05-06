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

package translate

import (
	"testing"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	. "github.com/katydid/validator-go-jsonschema/jsonschema/std/testutil"
)

func TestType(t *testing.T) {
	inputs := map[string]TagType{
		`{"default": {}}`:        ObjectTag,
		`{"default": []}`:        ArrayTag,
		`{"default": [1]}`:       ArrayTag,
		`{"default": {"a":"b"}}`: ObjectTag,
		`{"default": 0}`:         FieldTag,
		`{"default": 1.0}`:       FieldTag,
		`{"default": 1.1}`:       FieldTag,
		`{"default": ""}`:        FieldTag,
		`{"default": "a"}`:       FieldTag,
		`{"default": true}`:      FieldTag,
		`{"default": false}`:     FieldTag,
		`{"default": null}`:      FieldTag,
	}
	for input, want := range inputs {
		s := Must(schema.ParseSchema([]byte(input)))(t)
		ExpectErr(t, input, want, guessDefault(s))
	}
}
