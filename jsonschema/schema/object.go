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
	"encoding/json"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

// http://json-schema.org/latest/json-schema-validation.html#anchor53
type Object struct {
	MaxProperties        *uint64     `json:"maxProperties,omitempty"`
	MinProperties        uint64      `json:"minProperties,omitempty"`
	Required             []string    `json:"required,omitempty"`
	AdditionalProperties *Additional `json:"additionalProperties,omitempty"`

	Properties *Properties `json:"properties,omitempty"`
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	//http://json-schema.org/latest/json-schema-validation.html#anchor64
	//  The value of "patternProperties" MUST be an object. Each property name of this object SHOULD be a valid regular expression, according to the ECMA 262 regular expression dialect. Each property value of this object MUST be an object, and each object MUST be a valid JSON Schema.
	PatternProperties map[string]*Schema  `json:"patternProperties,omitempty"`
	Dependencies      *Dependencies       `json:"dependencies,omitempty"` // Kept for compatibliy with versions before Draft 2019-09. Now this is two seperate fields dependentRequired and dependentSchemas
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"`
	DependentSchemas  map[string]*Schema  `json:"dependentSchemas,omitempty"`
}

func (this Object) HasObjectConstraints() bool {
	return this.MaxProperties != nil || this.MinProperties > 0 ||
		this.Required != nil || this.AdditionalProperties != nil ||
		this.Properties != nil || this.PatternProperties != nil ||
		this.Dependencies != nil || this.DependentRequired != nil || this.DependentSchemas != nil
}

type Properties map[string]*Schema

/*
   "type": "object",
   "additionalProperties": { "$ref": "#" },
   "default": {}
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor64
//  The value of "properties" MUST be an object. Each value of this object MUST be an object, and each object MUST be a valid JSON Schema.
// But sometimes it isn't and then we can ignore those values.
func (this *Properties) UnmarshalJSON(data []byte) error {
	var objmap map[string]json.RawMessage
	if err := std.UnmarshalJSON(data, &objmap); err != nil {
		return err
	}
	props := map[string]*Schema{}
	for k := range objmap {
		var s *Schema
		if err := std.UnmarshalJSON(objmap[k], &s); err != nil {
			// ignore non schema values
			continue
		}
		props[k] = s
	}
	*this = props
	return nil
}

func (this *Object) GetProperties() map[string]*Schema {
	if this == nil {
		return map[string]*Schema{}
	}
	if this.Properties == nil {
		return map[string]*Schema{}
	}
	return *this.Properties
}
