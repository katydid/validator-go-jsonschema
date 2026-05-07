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

/*
	{
	    "id": "http://json-schema.org/draft-04/schema#",
	    "$schema": "http://json-schema.org/draft-04/schema#",
	    "description": "Core schema meta-schema",
	    "definitions": {
	        "schemaArray": {
	            "type": "array",
	            "minItems": 1,
	            "items": { "$ref": "#" }
	        },
	        "positiveInteger": {
	            "type": "integer",
	            "minimum": 0
	        },
	        "positiveIntegerDefault0": {
	            "allOf": [ { "$ref": "#/definitions/positiveInteger" }, { "default": 0 } ]
	        },
	        "simpleTypes": {
	            "enum": [ "array", "boolean", "integer", "null", "number", "object", "string" ]
	        },
	        "stringArray": {
	            "type": "array",
	            "items": { "type": "string" },
	            "minItems": 1,
	            "uniqueItems": true
	        }
	    },
	    "type": "object",
	    "properties": {
	        "id": {
	            "type": "string",
	            "format": "uri"
	        },
	        "$schema": {
	            "type": "string",
	            "format": "uri"
	        },
	        "title": {
	            "type": "string"
	        },
	        "description": {
	            "type": "string"
	        },
	        "default": {},
	        "multipleOf": {
	            "type": "number",
	            "minimum": 0,
	            "exclusiveMinimum": true
	        },
	        "maximum": {
	            "type": "number"
	        },
	        "exclusiveMaximum": {
	            "type": "boolean",
	            "default": false
	        },
	        "minimum": {
	            "type": "number"
	        },
	        "exclusiveMinimum": {
	            "type": "boolean",
	            "default": false
	        },
	        "maxLength": { "$ref": "#/definitions/positiveInteger" },
	        "minLength": { "$ref": "#/definitions/positiveIntegerDefault0" },
	        "pattern": {
	            "type": "string",
	            "format": "regex"
	        },
	        "additionalItems": {
	            "anyOf": [
	                { "type": "boolean" },
	                { "$ref": "#" }
	            ],
	            "default": {}
	        },
	        "items": {
	            "anyOf": [
	                { "$ref": "#" },
	                { "$ref": "#/definitions/schemaArray" }
	            ],
	            "default": {}
	        },
	        "maxItems": { "$ref": "#/definitions/positiveInteger" },
	        "minItems": { "$ref": "#/definitions/positiveIntegerDefault0" },
	        "uniqueItems": {
	            "type": "boolean",
	            "default": false
	        },
	        "maxProperties": { "$ref": "#/definitions/positiveInteger" },
	        "minProperties": { "$ref": "#/definitions/positiveIntegerDefault0" },
	        "required": { "$ref": "#/definitions/stringArray" },
	        "additionalProperties": {
	            "anyOf": [
	                { "type": "boolean" },
	                { "$ref": "#" }
	            ],
	            "default": {}
	        },
	        "definitions": {
	            "type": "object",
	            "additionalProperties": { "$ref": "#" },
	            "default": {}
	        },
	        "properties": {
	            "type": "object",
	            "additionalProperties": { "$ref": "#" },
	            "default": {}
	        },
	        "patternProperties": {
	            "type": "object",
	            "additionalProperties": { "$ref": "#" },
	            "default": {}
	        },
	        "dependencies": {
	            "type": "object",
	            "additionalProperties": {
	                "anyOf": [
	                    { "$ref": "#" },
	                    { "$ref": "#/definitions/stringArray" }
	                ]
	            }
	        },
	        "enum": {
	            "type": "array",
	            "minItems": 1,
	            "uniqueItems": true
	        },
	        "type": {
	            "anyOf": [
	                { "$ref": "#/definitions/simpleTypes" },
	                {
	                    "type": "array",
	                    "items": { "$ref": "#/definitions/simpleTypes" },
	                    "minItems": 1,
	                    "uniqueItems": true
	                }
	            ]
	        },
	        "allOf": { "$ref": "#/definitions/schemaArray" },
	        "anyOf": { "$ref": "#/definitions/schemaArray" },
	        "oneOf": { "$ref": "#/definitions/schemaArray" },
	        "not": { "$ref": "#" }
	    },
	    "dependencies": {
	        "exclusiveMaximum": [ "maximum" ],
	        "exclusiveMinimum": [ "minimum" ]
	    },
	    "default": {}
	}
*/
type Schema struct {
	Id          string             `json:"id,omitempty"`
	Schema      string             `json:"$schema,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Default     any                `json:"default,omitempty"`
	Definitions map[string]*Schema `json:"definitions,omitempty"`
	Numeric
	String
	Array
	Object
	Operators
	Type *Type `json:"type,omitempty"`

	Ref    string `json:"$ref,omitempty"`
	Format string `json:"format,omitempty"`
}

func (this Schema) GetType() []SimpleType {
	if this.Type != nil {
		return *this.Type
	}
	return nil
}
