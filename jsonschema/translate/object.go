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

func translateObject(s *schema.Schema) (*ast.Pattern, error) {
	if s.MaxProperties != nil {
		return nil, fmt.Errorf("TODO: maxProperties not supported")
	}
	if s.MinProperties > 0 {
		return nil, fmt.Errorf("TODO: minProperties not supported")
	}

	// TODO: Do some with dependencies
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

	names := std.SortedKeys(s.Properties)
	patternNames := std.SortedKeys(s.PatternProperties)

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
			// TODO: Investigate whether this is correct.
			additional = ast.NewZeroOrMore(
				ast.NewTreeNode(ast.NewAnyName(), typ),
			)
		}
	}

	props, err := newProperties(s)
	if err != nil {
		return nil, err
	}

	p, err := translateProps(props)
	if err != nil {
		return nil, err
	}

	required, err := translateRequired(s.Required)
	if err != nil {
		return nil, err
	}

	if len(props) == 0 {
		return NewObjectNode(ast.NewAnd(additional, required)), nil
	}

	return NewObjectNode(ast.NewAnd(ast.NewInterleave(p, additional), required)), nil
}

type property struct {
	key   string
	name  *ast.NameExpr
	child *ast.Pattern
}

func newProperties(s *schema.Schema) ([]*property, error) {
	names := std.SortedKeys(s.Properties)
	patternNames := std.SortedKeys(s.PatternProperties)
	props := make([]*property, 0, len(names)+len(patternNames))
	for _, name := range names {
		p, err := newProperty(name, s.Properties[name])
		if err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	for _, name := range patternNames {
		p, err := newPatternProperty(name, s.PatternProperties[name])
		if err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	return props, nil
}

func newProperty(name string, s *schema.Schema) (*property, error) {
	child, err := translate(s)
	if err != nil {
		return nil, err
	}
	return &property{
		key:   name,
		name:  ast.NewStringName(name),
		child: child,
	}, nil
}

func newPatternProperty(name string, s *schema.Schema) (*property, error) {
	child, err := translate(s)
	if err != nil {
		return nil, err
	}
	return &property{
		key:   name,
		name:  ast.NewRegexName(name),
		child: child,
	}, nil
}

// patternProperties has a pattern as a name and a child schema, for example:
//
//	"patternProperties": {
//		  "aaa*": {"maximum": 20}
//	}
//
// This is support by a regular expression in the name and the normal schema matching operators for the child schema.
//
// Things get more complicated when matching multiple patternProperties.
// If only one name matches then only that child schema is taken into account and if that schema matches it is match.
// But if both names matches then both schemas have to match.
// Take this example from "multiple simultaneous patternProperties are validated" in patternProperties.json in the draft4 testsuite:
//
//	"patternProperties": {
//	    "a*": {"type": "integer"},
//	    "aaa*": {"maximum": 20}
//	}
//
// {"aaaa": 31} => false, because both names match, so both schemas have to match.
// {"a": 21} => true, because only one names matches, only that schema has to match.
//
// That means we have to take into account all combinations and translate it to:
// ("a*"&"aaa*"):{"type": "integer"}&{"maximum": 20}
// | ("a*"&!"aaa*"):{"maximum": 20}
// | (!"a*"&!"aaa*"):{"type": "integer"}
// We calculate these combinations as complementary subsets.
func translateProps(props []*property) (*ast.Pattern, error) {
	var res []*ast.Pattern

	propCompSubsets := std.ComplementarySubsets(props)
	for _, propCompSubset := range propCompSubsets {
		names := make([]*ast.NameExpr, 0, len(propCompSubset.Left)+len(propCompSubset.Right))
		for _, prop := range propCompSubset.Left {
			names = append(names, prop.name)
		}
		for _, prop := range propCompSubset.Right {
			names = append(names, ast.NewAnyNameExcept(prop.name))
		}
		name := ast.NewNameConj(names...)

		children := make([]*ast.Pattern, 0, len(propCompSubset.Left))
		for _, prop := range propCompSubset.Left {
			children = append(children, prop.child.Clone())
		}

		r := ast.NewTreeNode(name, ast.NewAnd(children...))
		res = append(res, r)
	}

	return ast.NewZeroOrMore(ast.NewOr(res...)), nil
}

func translateRequired(required []string) (*ast.Pattern, error) {
	res := []*ast.Pattern{}
	for _, req := range required {
		res = append(res, ast.NewTreeNode(ast.NewStringName(req), ast.NewZAny()))
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewInterleave(res...), nil
}
