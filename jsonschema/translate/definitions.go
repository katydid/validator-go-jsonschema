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
	"maps"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func findDefinitions(s *schema.Schema) (map[string]*schema.Schema, error) {
	defs := make(map[string]*schema.Schema)
	maps.Copy(defs, s.Definitions)
	if _, ok := defs["main"]; ok {
		return nil, fmt.Errorf("main is a reserved definition name for katydid")
	}
	// katydid starts with the main pattern
	if len(s.Id) > 0 {
		defs[s.Id] = s
		defs["main"] = &schema.Schema{Ref: s.Id}
	} else {
		defs["main"] = s
	}
	return defs, nil
}

func translateDefinitions(s *schema.Schema) (map[string]*ast.Pattern, error) {
	refs := make(map[string]*ast.Pattern)
	defs, err := findDefinitions(s)
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
