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
*/
type Type []SimpleType

func (this *Type) HasArray() bool {
	return this.has(TypeArray)
}

func (this *Type) HasNumeric() bool {
	return this.has(TypeInteger) || this.has(TypeNumber)
}

func (this *Type) HasString() bool {
	return this.has(TypeString)
}

func (this *Type) has(s SimpleType) bool {
	if this == nil {
		return false
	}
	for _, t := range *this {
		if t == s {
			return true
		}
	}
	return false
}

func (this *Type) Single() bool {
	if this == nil {
		return false
	}
	return len(*this) == 1
}

func (this *Type) UnmarshalJSON(buf []byte) error {
	t := []SimpleType{}
	var s string
	if err := std.UnmarshalJSON(buf, &s); err == nil {
		simpleType, err := newSimpleType(s)
		if err != nil {
			return err
		}
		t = append(t, simpleType)
		*this = t
		return nil
	}
	var ss []string
	if err := std.UnmarshalJSON(buf, &ss); err != nil {
		return err
	}
	simpleTypes := make(map[string]struct{})
	for _, s := range ss {
		simpleType, err := newSimpleType(s)
		if err != nil {
			return err
		}
		if _, ok := simpleTypes[s]; ok {
			err := fmt.Errorf("type alternatives are not unique, duplicate %s found", simpleType)
			return err
		}
		simpleTypes[s] = struct{}{}
		t = append(t, simpleType)
	}
	*this = t
	return nil
}

type SimpleType string

const (
	TypeUnknown = SimpleType("unknown")
	TypeArray   = SimpleType("array")
	TypeBoolean = SimpleType("boolean")
	TypeInteger = SimpleType("integer")
	TypeNull    = SimpleType("null")
	TypeNumber  = SimpleType("number")
	TypeObject  = SimpleType("object")
	TypeString  = SimpleType("string")
)

func newSimpleType(s string) (SimpleType, error) {
	switch s {
	case "array":
		return TypeArray, nil
	case "boolean":
		return TypeBoolean, nil
	case "integer":
		return TypeInteger, nil
	case "null":
		return TypeNull, nil
	case "number":
		return TypeNumber, nil
	case "object":
		return TypeObject, nil
	case "string":
		return TypeString, nil
	}
	err := fmt.Errorf("unknown simpletype %s", s)
	return TypeUnknown, err
}
