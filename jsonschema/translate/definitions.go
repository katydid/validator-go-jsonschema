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

func findMainDefinitions(s *schema.Schema) (map[string]*schema.Schema, error) {
	defs := make(map[string]*schema.Schema)
	err := findSchemaDefinitions(s, s, defs)
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

func findSchemaDefinitions(root *schema.Schema, s *schema.Schema, res map[string]*schema.Schema) error {
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
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	// TODO: s.Array.Additional. Right now it does not a Schema inside, but it will probably in future.
	if sch := s.Array.GetItems().GetObject(); sch != nil {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.Array.GetItems().GetArray() {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	// TODO s.Object.AdditionalProperties. Right now it does not a Schema inside, but it will probably in future.
	for _, sch := range s.Object.Properties {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.Object.PatternProperties {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	if s.Object.Dependencies != nil {
		for _, dep := range *s.Object.Dependencies {
			if sch := dep.Schema; sch != nil {
				if err := findSchemaDefinitions(root, sch, res); err != nil {
					return err
				}
			}
		}
	}
	for _, sch := range s.AllOf {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.AnyOf {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	for _, sch := range s.OneOf {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	if sch := s.Not; sch != nil {
		if err := findSchemaDefinitions(root, sch, res); err != nil {
			return err
		}
	}
	if len(s.Ref) > 0 {
		if s.Ref == "#" {
			// main reference is already added, so nothing to do here.
		} else if strings.HasPrefix(s.Ref, "#/definitions/") {
			// other definitions are already added, so nothing to do there.
		} else if strings.HasPrefix(s.Ref, "#/") {
			pointer := strings.Split(s.Ref[2:], "/")
			sch := findSchema(pointer, root)
			if sch == nil {
				return fmt.Errorf("could not find schema for %s", s.Ref)
			}
			res[s.Ref] = sch
		} else if strings.HasPrefix(s.Ref, "http") {
			return fmt.Errorf("remote ref not supported")
		} else if strings.HasPrefix(s.Ref, "file") {
			return fmt.Errorf("file ref not supported")
		} else {
			return fmt.Errorf("unknonw ref type not supported")
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
		sch, ok := s.Properties[pointer[1]]
		if !ok {
			return nil
		}
		if len(pointer) > 2 {
			return findSchema(pointer[2:], sch)
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
		// TODO: support more relative pointers
		return nil
	}
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
