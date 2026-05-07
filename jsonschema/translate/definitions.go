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

import (
	"fmt"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func findMainDefinitions(s *schema.Schema) (map[string]*schema.Schema, error) {
	defs := make(map[string]*schema.Schema)
	err := findSchemaDefinitions(s, defs)
	if err != nil {
		return nil, err
	}

	if _, ok := defs["main"]; ok {
		return nil, fmt.Errorf("main is a reserved definition name for katydid")
	}
	// katydid starts with the main pattern
	if len(s.Id) > 0 {
		defs["main"] = &schema.Schema{Ref: s.Id}
		defs[s.Id] = s
	} else {
		defs["main"] = s
	}
	return defs, nil
}

func findSchemaDefinitions(s *schema.Schema, res map[string]*schema.Schema) error {
	for name, sch := range s.Definitions {
		realname := name
		if len(sch.Id) > 0 {
			realname = sch.Id
		}
		if _, ok := res[realname]; ok {
			return fmt.Errorf("duplicate definition name: %s", realname)
		}
		res[realname] = sch
	}
	for _, sch := range s.Definitions {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	// TODO: s.Array.Additional. Right now it does not a Schema inside, but it will probably in future.
	if sch := s.Array.GetItems().GetObject(); sch != nil {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.Array.GetItems().GetArray() {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	// TODO s.Object.AdditionalProperties. Right now it does not a Schema inside, but it will probably in future.
	for _, sch := range s.Object.Properties {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.Object.PatternProperties {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	if s.Object.Dependencies != nil {
		for _, dep := range *s.Object.Dependencies {
			if sch := dep.Schema; sch != nil {
				if err := findSchemaDefinitions(sch, res); err != nil {
					return err
				}
			}
		}
	}
	for _, sch := range s.AllOf {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.AnyOf {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.OneOf {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}
	if sch := s.Not; sch != nil {
		if err := findSchemaDefinitions(sch, res); err != nil {
			return err
		}
	}

	return nil
}

func translateDefinitions(s *schema.Schema) (map[string]*ast.Pattern, error) {
	refs := make(map[string]*ast.Pattern)
	defs, err := findMainDefinitions(s)
	if err != nil {
		return nil, err
	}
	names := std.SortedKeys(defs)
	for _, name := range names {
		p, err := translate(defs[name])
		if err != nil {
			return nil, err
		}
		refs[name] = p
	}
	return refs, nil
}
