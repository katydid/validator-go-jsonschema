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
	"fmt"

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
	Type SimpleType
}

type aSchema struct {
	Type *Type `json:"type"`
}

func (this *Additional) UnmarshalJSON(buf []byte) error {
	var b bool
	if err := std.UnmarshalJSON(buf, &b); err == nil {
		*this = Additional{Bool: &b}
		return nil
	}
	s := &aSchema{}
	if err := std.UnmarshalJSON(buf, s); err != nil {
		return err
	}
	if s.Type == nil {
		return fmt.Errorf("the additional(Items|Properties) field is empty")
	}
	typ := *s.Type
	if len(typ) > 1 {
		return fmt.Errorf("the additional(Items|Properties) field's type field has more than one element")
	}
	if len(typ) == 0 {
		panic(fmt.Errorf("%#v buf = %s", s.Type, string(buf)))
	}
	*this = Additional{Type: typ[0]}
	return nil
}
