//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package schema

import (
	"encoding/json"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

func ParseSchema(jsonSchema []byte) (*Schema, error) {
	schema := &Schema{}
	if err := std.UnmarshalJSON(jsonSchema, schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func (this *Schema) JsonString() string {
	data, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Schema struct {
	Id          string `json:"id,omitempty"`
	Anchor      string `json:"$anchor,omitempty"`
	Schema      string `json:"$schema,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`

	//  This keyword's value MUST be an object. Each member value of this object MUST be a valid JSON Schema.
	Definitions map[string]*Schema `json:"definitions,omitempty"`

	Numeric
	String
	Array
	Object
	Operators
	Type *Type `json:"type,omitempty"`
	// Const is *any because a JSON null (Go nil) is a valid value.
	Const *any `json:"const,omitempty"`

	Ref string `json:"$ref,omitempty"`
}

func (this Schema) GetType() []SimpleType {
	if this.Type != nil {
		return *this.Type
	}
	return nil
}

func (this Schema) GetVersion() Version {
	return detectVersion(this.Schema)
}

func (this *Schema) SetDefaultVersion(defaultVersion Version) {
	setDefaultVersion(this, defaultVersion)
}
