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
	"sort"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go/validator/ast"
)

func translateObject(s *schema.Schema) (*ast.Pattern, error) {
	if s.MaxProperties != nil {
		return nil, fmt.Errorf("maxProperties not supported")
	}
	if s.MinProperties > 0 {
		return nil, fmt.Errorf("minProperties not supported")
	}
	required := make(map[string]struct{})
	for _, req := range s.Required {
		required[req] = struct{}{}
	}
	requiredIf := make(map[string][]string)
	moreProperties := make(map[string]*schema.Schema)
	if s.Dependencies != nil {
		deps := *s.Dependencies
		for name, dep := range deps {
			if len(dep.RequiredProperty) > 0 {
				requiredIf[name] = deps[name].RequiredProperty
			} else {
				moreProperties[name] = deps[name].Schema
			}
		}
	}

	names := []string{}
	for name := range s.Properties {
		names = append(names, name)
	}
	sort.Strings(names)

	patternNames := []string{}
	for name := range s.PatternProperties {
		patternNames = append(patternNames, name)
	}
	sort.Strings(patternNames)

	additional := ast.NewZAny()
	if len(names) > 0 || len(patternNames) > 0 {
		nameExprs := make([]*ast.NameExpr, len(names)+len(patternNames))
		for i, name := range names {
			nameExprs[i] = ast.NewStringName(name)
		}
		for i, name := range patternNames {
			nameExprs[i+len(names)] = ast.NewRegexName(name)
		}
		additional = ast.NewZeroOrMore(
			ast.NewTreeNode(ast.NewAnyNameExcept(
				ast.NewNameChoice(nameExprs...),
			), ast.NewZAny()),
		)
	}
	if s.AdditionalProperties != nil {
		if s.AdditionalProperties.Bool != nil && !(*s.AdditionalProperties.Bool) {
			additional = ast.NewEmpty()
		} else if s.AdditionalProperties.Type != schema.TypeUnknown {
			typ, err := translateType(s.AdditionalProperties.Type)
			if err != nil {
				return nil, err
			}
			additional = ast.NewZeroOrMore(
				ast.NewTreeNode(ast.NewAnyName(), typ),
			)
		}
	}
	patterns := make(map[string]*ast.Pattern)
	for _, name := range names {
		child, err := translate(s.Properties[name])
		if err != nil {
			return nil, err
		}
		patterns[name] = ast.NewTreeNode(ast.NewStringName(name), child)
	}
	for _, name := range patternNames {
		child, err := translate(s.PatternProperties[name])
		if err != nil {
			return nil, err
		}
		patterns[name] = ast.NewTreeNode(ast.NewRegexName(name), child)
	}
	for _, name := range names {
		if _, ok := requiredIf[name]; ok {
			return nil, fmt.Errorf("dependencies are not supported")
		}
		if _, ok := moreProperties[name]; ok {
			return nil, fmt.Errorf("dependencies are not supported")
		}
		if _, ok := required[name]; !ok {
			patterns[name] = ast.NewOptional(patterns[name])
		}
	}

	patternList := make([]*ast.Pattern, 0, len(patterns))
	for _, name := range names {
		patternList = append(patternList, patterns[name])
	}
	patternList = append(patternList, additional)
	// TODO: Be more specific and create ast.NewTagName for "object"
	return ast.NewTreeNode(ast.NewStringName("object"), ast.NewInterleave(patternList...)), nil
}
