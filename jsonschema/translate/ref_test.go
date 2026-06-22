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

import "testing"

// # Draft 4 test case:
//
//	{
//		"description": "$ref prevents a sibling id from changing the base uri",
//		"schema": {
//			"id": "http://localhost:1234/sibling_id/base/",
//			"definitions": {
//				"foo": {
//					"id": "http://localhost:1234/sibling_id/foo.json",
//					"type": "string"
//				},
//				"base_foo": {
//					"$comment": "this canonical uri is http://localhost:1234/sibling_id/base/foo.json",
//					"id": "foo.json",
//					"type": "number"
//				}
//			},
//			"allOf": [
//				{
//					"$comment": "$ref resolves to http://localhost:1234/sibling_id/base/foo.json, not http://localhost:1234/sibling_id/foo.json",
//					"id": "http://localhost:1234/sibling_id/",
//					"$ref": "foo.json"
//				}
//			]
//		},
//		"tests": [
//			{
//				"description": "$ref resolves to /definitions/base_foo, data does not validate",
//				"data": "a",
//				"valid": false
//			},
//			{
//				"description": "$ref resolves to /definitions/base_foo, data validates",
//				"data": 1,
//				"valid": true
//			}
//		]
//	},
func TestRefPreventsSibling(t *testing.T) {
	// prefix, parentId, name, id, anchor
	defName, err := definitionToDefName("", "http://localhost:1234/sibling_id/base/", "base_foo", "foo.json", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("http://localhost:1234/sibling_id/base/", "foo.json")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
}
