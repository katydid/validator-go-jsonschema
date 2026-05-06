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
"dependencies": {
	"type": "object",
	"additionalProperties": {
		"anyOf": [
			{ "$ref": "#" },
			{ "$ref": "#/definitions/stringArray" }
		]
	}
},
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor70
//  This keyword's value MUST be an object. Each value of this object MUST be either an object or an array.
//  If the value is an object, it MUST be a valid JSON Schema. This is called a schema dependency.
//  If the value is an array, it MUST have at least one element. Each element MUST be a string, and elements in the array MUST be unique. This is called a property dependency.
//http://spacetelescope.github.io/understanding-json-schema/reference/object.html#dependencies
type Dependencies map[string]*Dependency

type Dependency struct {
	Schema           *Schema
	RequiredProperty []string
}

func (this *Dependency) UnmarshalJSON(buf []byte) error {
	var s *Schema
	if err := std.UnmarshalJSON(buf, &s); err == nil {
		*this = Dependency{Schema: s}
		return nil
	}
	var ss []string
	if err := std.UnmarshalJSON(buf, &ss); err != nil {
		return err
	}
	*this = Dependency{RequiredProperty: ss}
	checkUnique := make(map[string]struct{})
	for _, s := range this.RequiredProperty {
		if _, ok := checkUnique[s]; ok {
			return fmt.Errorf("duplicate found in property dependency list %s", s)
		}
		checkUnique[s] = struct{}{}
	}
	return nil
}
