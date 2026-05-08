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
"additionalItems": {
	"anyOf": [
		{ "type": "boolean" },
		{ "$ref": "#" }
	],
	"default": {}
},

"additionalProperties": {
	"anyOf": [
		{ "type": "boolean" },
		{ "$ref": "#" }
	],
	"default": {}
},

"anyOf": [
	{ "type": "boolean" },
	{ "$ref": "#" }
],
"default": {}
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor37
//  The value of "additionalItems" MUST be either a boolean or an object. If it is an object, this object MUST be a valid JSON Schema.
//http://json-schema.org/latest/json-schema-validation.html#anchor49
//  The value of "additionalProperties" MUST be a boolean or an object. If it is an object, it MUST also be a valid JSON Schema.
type Additional struct {
	Bool *bool
	//Typically only the type field of the jsonschema is set.
	Schema *Schema
}

func (this *Additional) UnmarshalJSON(buf []byte) error {
	var b bool
	if err := std.UnmarshalJSON(buf, &b); err == nil {
		*this = Additional{Bool: &b}
		return nil
	}
	s := &Schema{}
	if err := std.UnmarshalJSON(buf, s); err != nil {
		return err
	}
	*this = Additional{Schema: s}
	return nil
}

func (this *Additional) GetSchema() *Schema {
	if this == nil {
		return nil
	}
	return this.Schema
}
