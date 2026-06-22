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
	"strconv"
	"strings"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func findDefinitions(s *schema.Schema) (map[string]*schema.Schema, error) {
	defs := make(map[string]*schema.Schema)
	err := findSchemaDefinitions(s, "", s, defs)
	if err != nil {
		return nil, err
	}
	return defs, nil
}

func findSchemaDefinitions(root *schema.Schema, prefix string, s *schema.Schema, res map[string]*schema.Schema) error {
	for _, name := range std.SortedKeys(s.Definitions) {
		sch := s.Definitions[name]
		defname, err := definitionToDefName(prefix, sch.Id, sch.Anchor, s.Definitions[name].Id, name)
		if err != nil {
			return err
		}
		if _, ok := res[defname]; ok {
			return fmt.Errorf("duplicate definition name: %s", defname)
		}
		res[defname] = sch
	}
	for _, name := range std.SortedKeys(s.Definitions) {
		sch := s.Definitions[name]
		newprefix := definitionToPrefix(prefix, sch.Id, name)
		if err := findSchemaDefinitions(root, newprefix, sch, res); err != nil {
			return err
		}
	}
	if sch := s.Array.AdditionalItems.GetSchema(); sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/additionalItems", sch, res); err != nil {
			return err
		}
	}
	if sch := s.Array.GetItems().GetObject(); sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/items", sch, res); err != nil {
			return err
		}
	}
	for i, sch := range s.Array.GetItems().GetArray() {
		if err := findSchemaDefinitions(root, prefix+"/items/"+strconv.Itoa(i), sch, res); err != nil {
			return err
		}
	}
	if sch := s.Object.AdditionalProperties.GetSchema(); sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/additionalProperties", sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.Object.GetProperties() {
		if err := findSchemaDefinitions(root, prefix+"/properties", sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.Object.PatternProperties {
		if err := findSchemaDefinitions(root, prefix+"/patternProperties", sch, res); err != nil {
			return err
		}
	}
	if s.Operators.Dependencies != nil {
		for name, dep := range *s.Operators.Dependencies {
			if sch := dep.Schema; sch != nil {
				if err := findSchemaDefinitions(root, prefix+"/dependencies/"+name, sch, res); err != nil {
					return err
				}
			}
		}
	}
	for _, sch := range s.AllOf {
		if err := findSchemaDefinitions(root, prefix+"/allOf", sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.AnyOf {
		if err := findSchemaDefinitions(root, prefix+"/anyOf", sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.OneOf {
		if err := findSchemaDefinitions(root, prefix+"/oneOf", sch, res); err != nil {
			return err
		}
	}
	if sch := s.Not; sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/not", sch, res); err != nil {
			return err
		}
	}
	if sch := s.If; sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/if", sch, res); err != nil {
			return err
		}
	}
	if sch := s.Then; sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/then", sch, res); err != nil {
			return err
		}
	}
	if sch := s.Else; sch != nil {
		if err := findSchemaDefinitions(root, prefix+"/else", sch, res); err != nil {
			return err
		}
	}
	if len(s.Ref) > 0 {
		if s.Ref == "#" {
			// main reference is already added, so nothing to do here.
		} else if strings.HasPrefix(s.Ref, "#/definitions/") {
			// other definitions are already added, so nothing to do there.
		} else if strings.HasPrefix(s.Ref, "#/") {
			pointer, err := parsePointer(s.Ref)
			if err != nil {
				return err
			}
			sch := findSchema(pointer, root)
			if sch == nil {
				return fmt.Errorf("could not find schema for %s", s.Ref)
			}
			defName, err := refToDefName(s.Id, s.Ref)
			if err != nil {
				return err
			}
			res[defName] = sch
		} else if strings.HasPrefix(s.Ref, "http") {
			defName, err := refToDefName(s.Id, s.Ref)
			if err != nil {
				return err
			}
			switch defName {
			case "http://json-schema.org/draft-04/schema":
				if _, ok := res[defName]; ok {
					return nil
				}
				// handle meta schema validation with unique items
				s, err := schema.ParseSchema([]byte(schema.SchemaDraft4ExlcudeUniqueItems))
				if err != nil {
					return err
				}
				s.SetDefaultVersion(schema.VersionDraft4)
				res[defName] = s
				if err := findSchemaDefinitions(s, "", s, res); err != nil {
					return err
				}
			}
			// TODO if it has a local part that goes deeper into the schema and does not just reference it, then there is more work to do here.
		} else if strings.HasPrefix(s.Ref, "file") {
			return fmt.Errorf("file ref not supported")
		} else {
			// TODO if it has a local part that goes deeper into the schema and does not just reference it, then there is more work to do here.
		}
	}

	return nil
}

func findSchema(pointer []string, s *schema.Schema) *schema.Schema {
	if len(pointer) == 0 {
		return nil
	}
	name := pointer[0]
	switch name {
	case "properties":
		if len(pointer) < 2 {
			return nil
		}
		sch, ok := s.GetProperties()[pointer[1]]
		if !ok {
			return nil
		}
		if len(pointer) > 2 {
			return findSchema(pointer[2:], sch)
		}
		return sch
	case "definitions":
		if s.Definitions == nil {
			return nil
		}
		sch, ok := s.Definitions[pointer[1]]
		if !ok {
			return nil
		}
		if len(pointer) > 2 {
			return findSchema(pointer[2:], sch)
		}
		return sch
	case "if":
		sch := s.If
		if sch == nil {
			return nil
		}
		if len(pointer) > 1 {
			return findSchema(pointer[1:], sch)
		}
		return sch
	case "else":
		sch := s.Else
		if sch == nil {
			return nil
		}
		if len(pointer) > 1 {
			return findSchema(pointer[1:], sch)
		}
		return sch
	case "then":
		sch := s.Then
		if sch == nil {
			return nil
		}
		if len(pointer) > 1 {
			return findSchema(pointer[1:], sch)
		}
		return sch
	case "items":
		if sch := s.Items.GetObject(); sch != nil {
			if len(pointer) == 1 {
				return sch
			}
			return findSchema(pointer[1:], sch)
		}
		if len(pointer) < 2 {
			return nil
		}
		idx, err := strconv.Atoi(pointer[1])
		if err != nil {
			return nil
		}
		if idx >= len(s.Items.GetArray()) {
			return nil
		}
		sch := s.Items.GetArray()[idx]
		if len(pointer) == 2 {
			return sch
		}
		return findSchema(pointer[2:], sch)
	default:
		return nil
	}
}

func translateDefinitions(s *schema.Schema) (map[string]*ast.Pattern, error) {
	refs := make(map[string]*ast.Pattern)
	defs, err := findDefinitions(s)
	if err != nil {
		return nil, err
	}
	if _, ok := defs["main"]; ok {
		return nil, fmt.Errorf("main is a reserved definition name for katydid")
	}
	// katydid starts with the main pattern
	defs["main"] = s
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
