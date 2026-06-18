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

// http://json-schema.org/latest/json-schema-validation.html#anchor75
type Operators struct {
	/*
	   "type": "array",
	   "minItems": 1,
	   "uniqueItems": true
	*/
	Enum  []any     `json:"enum,omitempty"`
	AllOf []*Schema `json:"allOf,omitempty"`
	AnyOf []*Schema `json:"anyOf,omitempty"`
	OneOf []*Schema `json:"oneOf,omitempty"`
	Not   *Schema   `json:"not,omitempty"`

	If   *Schema `json:"if,omitempty"`
	Then *Schema `json:"then,omitempty"`
	Else *Schema `json:"else,omitempty"`
}

func (this Operators) HasOperatorConstraints() bool {
	return this.Enum != nil ||
		this.AllOf != nil || this.AnyOf != nil ||
		this.OneOf != nil || this.Not != nil ||
		this.If != nil
}
