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

import "testing"

func TestTypeList(t *testing.T) {
	schema := `
    {
      "properties": {
        "field": {
          "type": [
            "string",
            "array"
          ],
          "items": {
            "type": "string"
          }
        }
      }
    }`
	tests := map[string]bool{
		`{"field": "astring"}`:                    true,
		`{"field": 123}`:                          false,
		`{"field": ["anarray", "of", "strings"]}`: true,
		`{"field": [123, 456]}`:                   false,
	}
	for test, want := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := MatchBytes([]byte(schema), []byte(test))
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Error()
			}
		})
	}
}
