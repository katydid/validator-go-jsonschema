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
func TestDraft4RefPreventsSibling(t *testing.T) {
	want := "http://localhost:1234/sibling_id/base/foo.json"
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
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}

// # Draft 4 test case:
//
//		{
//		"description": "Recursive references between schemas",
//		"schema": {
//			"id": "http://localhost:1234/tree",
//			"description": "tree of nodes",
//			"type": "object",
//			"properties": {
//				"meta": {"type": "string"},
//				"nodes": {
//					"type": "array",
//					"items": {"$ref": "node"}
//				}
//			},
//			"required": ["meta", "nodes"],
//			"definitions": {
//				"node": {
//					"id": "http://localhost:1234/node",
//					"description": "node",
//					"type": "object",
//					"properties": {
//						"value": {"type": "number"},
//						"subtree": {"$ref": "tree"}
//					},
//					"required": ["value"]
//				}
//			}
//		},
//	}
func TestDraft4RecursiveReferences1(t *testing.T) {
	want := "http://localhost:1234/node"
	// prefix, parentId, name, id, anchor
	// is the prefix  "http://localhost:1234/tree" ?
	defName, err := definitionToDefName("", "http://localhost:1234/tree", "node", "http://localhost:1234/node", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("http://localhost:1234/tree", "node")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}

func TestDraft4RecursiveReferences2(t *testing.T) {
	want := "http://localhost:1234/tree"
	// prefix, parentId, name, id, anchor
	defName, err := definitionToDefName("", "", "", "http://localhost:1234/tree", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("http://localhost:1234/tree", "tree")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}

// # Draft 4 test case:
//
//	{
//		"description": "id inside an enum is not a real identifier",
//		"comment": "the implementation must not be confused by an id buried in the enum",
//		"schema": {
//			"definitions": {
//				"id_in_enum": {
//					"enum": [
//						{
//							"id": "https://localhost:1234/my_identifier.json",
//							"type": "null"
//						}
//					]
//				},
//				"real_id_in_schema": {
//					"id": "https://localhost:1234/my_identifier.json",
//					"type": "string"
//				},
//				"zzz_id_in_const": {
//					"const": {
//						"id": "https://localhost:1234/my_identifier.json",
//						"type": "null"
//					}
//				}
//			},
//			"anyOf": [
//				{ "$ref": "#/definitions/id_in_enum" },
//				{ "$ref": "https://localhost:1234/my_identifier.json" }
//			]
//		},
//		"tests": [
//			{
//				"description": "exact match to enum, and type matches",
//				"data": {
//					"id": "https://localhost:1234/my_identifier.json",
//					"type": "null"
//				},
//				"valid": true
//			},
//			{
//				"description": "match $ref to id",
//				"data": "a string to match #/definitions/id_in_enum",
//				"valid": true
//			},
//			{
//				"description": "no match on enum or $ref to id",
//				"data": 1,
//				"valid": false
//			}
//		]
//	}
func TestDraft4Id1(t *testing.T) {
	want := "definitions/id_in_enum"
	// prefix, parentId, name, id, anchor
	defName, err := definitionToDefName("", "", "id_in_enum", "", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("", "#/definitions/id_in_enum")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}

func TestDraft4Id2(t *testing.T) {
	want := "https://localhost:1234/my_identifier.json"
	// prefix, parentId, name, id, anchor
	defName, err := definitionToDefName("", "", "real_id_in_schema", "https://localhost:1234/my_identifier.json", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("", "https://localhost:1234/my_identifier.json")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}

// # Draft 4 test case:
//
//	{
//		"description": "id with file URI still resolves pointers - *nix",
//		"schema": {
//			"id": "file:///folder/file.json",
//			"definitions": {
//				"foo": {
//					"type": "number"
//				}
//			},
//			"allOf": [
//				{
//					"$ref": "#/definitions/foo"
//				}
//			]
//		},
//		"tests": [
//			{
//				"description": "number is valid",
//				"data": 1,
//				"valid": true
//			},
//			{
//				"description": "non-number is invalid",
//				"data": "a",
//				"valid": false
//			}
//		]
//	},
func TestDraft4File(t *testing.T) {
	want := "file:///folder/file.json/definitions/foo"
	// prefix, parentId, name, id, anchor
	defName, err := definitionToDefName("", "file:///folder/file.json", "foo", "", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("file:///folder/file.json", "#/definitions/foo")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}

// # Draft 4 Test
//
//	{
//		"description": "Location-independent identifier with base URI change in subschema",
//		"schema": {
//			"id": "http://localhost:1234/root",
//			"allOf": [{
//				"$ref": "http://localhost:1234/nested.json#foo"
//			}],
//			"definitions": {
//				"A": {
//					"id": "nested.json",
//					"definitions": {
//						"B": {
//							"id": "#foo",
//							"type": "integer"
//						}
//					}
//				}
//			}
//		},
//		"tests": [
//			{
//				"data": 1,
//				"description": "match",
//				"valid": true
//			},
//			{
//				"data": "a",
//				"description": "mismatch",
//				"valid": false
//			}
//		]
//	},
func TestDraft4LocationIndependent(t *testing.T) {
	want := "http://localhost:1234/nested.json/foo"
	// prefix, parentId, name, id, anchor
	defName, err := definitionToDefName("/definitions/A/definitions/B", "http://localhost:1234/nested.json", "B", "#foo", "")
	if err != nil {
		t.Fatal(err)
	}
	// parentId, name
	refName, err := refToDefName("http://localhost:1234/root", "http://localhost:1234/nested.json#foo")
	if err != nil {
		t.Fatal(err)
	}
	if defName != refName {
		t.Fatalf("definitionToDefName = %s, but refToDefName = %s", defName, refName)
	}
	if defName != want {
		t.Fatalf("got %s want %s", defName, want)
	}
}
