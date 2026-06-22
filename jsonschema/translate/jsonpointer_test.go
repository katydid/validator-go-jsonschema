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

	"github.com/katydid/validator-go-jsonschema/jsonschema/std/testutil"
)

func TestJSONPointer(t *testing.T) {
	expect := func(input string, want []string) {
		t.Helper()
		got, err := parsePointer(input)
		if err != nil {
			t.Errorf("given input %s error: %v", input, err)
		} else {
			testutil.ExpectErr(t, input, want, got)
		}
	}
	expect(`/definitions/tilde~0field`, []string{"definitions", "tilde~field"})
	expect(`/definitions/slash~1field`, []string{"definitions", "slash/field"})
	expect(`/definitions/percent%25field`, []string{"definitions", "percent%field"})
	expect(`#/definitions/percent%25field`, []string{"definitions", "percent%field"})
	expect(`#/definitions/percent%field`, []string{"definitions", "percent%field"})
	expect(`#/definitions//definitions/`, []string{"definitions", reservedWordForEmpty, "definitions", reservedWordForEmpty})
	expect(`#`, []string{})
	expect("http://json-schema.org/draft-04/schema", []string{"http://json-schema.org/draft-04/schema"})
	expect("http://json-schema.org/draft-04/schema/", []string{"http://json-schema.org/draft-04/schema/"})
	expect("http://json-schema.org/draft-04/schema#", []string{"http://json-schema.org/draft-04/schema"})
	expect("http://json-schema.org/draft-04/schema#abc", []string{"http://json-schema.org/draft-04/schema", "abc"})
}
