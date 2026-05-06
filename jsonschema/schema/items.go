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

package schema

import (
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

/*
"items": {
	"anyOf": [
		{ "$ref": "#" },
		{ "$ref": "#/definitions/schemaArray" }
	],
	"default": {}
},
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor37
//  The value of "items" MUST be either an object or an array. If it is an object, this object MUST be a valid JSON Schema. If it is an array, items of this array MUST be objects, and each of these objects MUST be a valid JSON Schema.
type Items struct {
	Object *Schema
	Array  []*Schema
}

func (this *Items) UnmarshalJSON(buf []byte) error {
	var s *Schema
	if err := std.UnmarshalJSON(buf, &s); err == nil {
		*this = Items{Object: s}
		return nil
	}
	schemas := []*Schema{}
	if err := std.UnmarshalJSON(buf, &schemas); err != nil {
		return err
	}
	*this = Items{Array: schemas}
	return nil
}
