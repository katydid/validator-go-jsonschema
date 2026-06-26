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
	"errors"
	"regexp"
	"slices"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func translateObject(parentId string, s *schema.Schema) (*ast.Pattern, error) {
	var constraints []*ast.Pattern
	if s.MaxProperties != nil {
		constraints = append(constraints, maxProperties(int(*s.MaxProperties)))
	}
	if s.MinProperties > 0 {
		constraints = append(constraints, minProperties(int(s.MinProperties)))
	}

	props, err := newProperties(parentId, s)
	if err != nil {
		return nil, err
	}
	if len(s.Required) > 0 {
		required, err := translateRequired(props)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, required)
	}

	additional, err := translateAdditionalProperties(parentId, s)
	if err != nil {
		return nil, err
	}
	if s.PropertyNames != nil {
		// currently only pattern property names are supported.
		if s.PropertyNames.Pattern != nil {
			if s.PatternProperties == nil {
				s.PatternProperties = make(map[string]*schema.Schema)
			}
			s.PatternProperties[*s.PropertyNames.Pattern] = &schema.Schema{}
			if s.AdditionalProperties != nil {
				if s.AdditionalProperties.Schema != nil {
					// we handle additional properties in propertyNames, so we can ignore it later.
					additional = ast.NewEmpty()
					s.PatternProperties[*s.PropertyNames.Pattern] = s.AdditionalProperties.Schema
				}
			}
		} else {
			return nil, errors.New("propertyNames support is limited to pattern")
		}
	}

	props, err = newProperties(parentId, s)
	if err != nil {
		return nil, err
	}
	p, err := translateProps(props)
	if err != nil {
		return nil, err
	}

	if len(props) == 0 {
		constraints = append(constraints, additional)
	} else {
		constraints = append(constraints, ast.NewInterleave(p, additional))
	}

	return newAnd(constraints...), nil
}

type property struct {
	key      string
	pattern  bool
	name     *ast.NameExpr
	child    *ast.Pattern
	required bool
}

func newProperties(parentId string, s *schema.Schema) ([]*property, error) {
	names := std.SortedKeys(s.GetProperties())
	patternNames := std.SortedKeys(s.PatternProperties)
	props := make([]*property, 0, len(names)+len(patternNames))
	requires := make([]string, len(s.Required))
	copy(requires, s.Required)
	for _, name := range names {
		index := slices.Index(requires, name)
		required := index != -1
		if required {
			requires = slices.Delete(requires, index, index+1)
		}
		p, err := newProperty(getId(parentId, s), name, s.GetProperties()[name], required)
		if err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	for _, name := range requires {
		p, err := newProperty(getId(parentId, s), name, &schema.Schema{}, true)
		if err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	for _, name := range patternNames {
		p, err := newPatternProperty(getId(parentId, s), name, s.PatternProperties[name])
		if err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	return props, nil
}

func newProperty(parentId string, name string, s *schema.Schema, required bool) (*property, error) {
	child, err := translate(parentId, s)
	if err != nil {
		return nil, err
	}
	return &property{
		key:      name,
		name:     ast.NewStringName(name),
		child:    child,
		required: required,
	}, nil
}

func newPatternProperty(parentId string, name string, s *schema.Schema) (*property, error) {
	child, err := translate(parentId, s)
	if err != nil {
		return nil, err
	}
	return &property{
		key:     name,
		pattern: true,
		name:    ast.NewRegexName(name),
		child:   child,
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
	overlapping, nonOverlapping, err := findOverlapping(props)
	if err != nil {
		return nil, err
	}

	for _, prop := range nonOverlapping {
		r := ast.NewTreeNode(prop.name, prop.child)
		res = append(res, r)
	}

	propCompSubsets := std.ComplementarySubsets(overlapping)
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

		r := ast.NewTreeNode(name, newAnd(children...))
		res = append(res, r)
	}

	return ast.NewZeroOrMore(newOr(res...)), nil
}

func findOverlapping(props []*property) ([]*property, []*property, error) {
	patternProps := []*property{}
	for _, prop := range props {
		if prop.pattern {
			patternProps = append(patternProps, prop)
		}
	}
	overlapping := patternProps
	others := []*property{}
	for i, prop := range props {
		if !prop.pattern {
			overlaps := false
			for _, patternProp := range patternProps {
				m, err := regexp.MatchString(patternProp.key, prop.key)
				if err != nil {
					return nil, nil, err
				}
				if m {
					overlaps = true
				}
			}
			if overlaps {
				overlapping = append(overlapping, props[i])
			} else {
				others = append(others, props[i])
			}
		}
	}
	return overlapping, others, nil
}

func translateRequired(props []*property) (*ast.Pattern, error) {
	res := []*ast.Pattern{}
	for _, prop := range props {
		if prop.required {
			res = append(res, ast.NewContains(ast.NewTreeNode(prop.name, ast.NewZAny())))
		}
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return newAnd(res...), nil
}

func maxProperties(n int) *ast.Pattern {
	ps := make([]*ast.Pattern, n+1)
	// one more than the maxProperties
	for i := 0; i < n+1; i++ {
		ps[i] = ast.NewTreeNode(ast.NewAnyName(), ast.NewZAny())
	}
	res := ast.NewConcat(ps...)
	return ast.NewNot(ast.NewConcat(res, ast.NewZAny()))
}

func minProperties(n int) *ast.Pattern {
	ps := make([]*ast.Pattern, n)
	for i := 0; i < n; i++ {
		ps[i] = ast.NewTreeNode(ast.NewAnyName(), ast.NewZAny())
	}
	return ast.NewConcat(ast.NewConcat(ps...), ast.NewZAny())
}

func translateAdditionalProperties(parentId string, s *schema.Schema) (*ast.Pattern, error) {
	additional := ast.NewZAny()

	names := std.SortedKeys(s.GetProperties())
	patternNames := std.SortedKeys(s.PatternProperties)

	otherNames := ast.NewAnyName()
	if len(names) > 0 || len(patternNames) > 0 {
		nameExprs := make([]*ast.NameExpr, len(names)+len(patternNames))
		for i, name := range names {
			nameExprs[i] = ast.NewStringName(name)
		}
		for i, name := range patternNames {
			nameExprs[i+len(names)] = ast.NewRegexName(name)
		}
		otherNames = ast.NewAnyNameExcept(
			ast.NewNameChoice(nameExprs...),
		)
		additional = ast.NewZeroOrMore(
			ast.NewTreeNode(otherNames, ast.NewZAny()),
		)
	}
	if s.AdditionalProperties != nil {
		if s.AdditionalProperties.Bool != nil && !(*s.AdditionalProperties.Bool) {
			additional = ast.NewEmpty()
		} else if s.AdditionalProperties.Schema != nil {
			p, err := translate(getId(parentId, s), s.AdditionalProperties.Schema)
			if err != nil {
				return nil, err
			}
			additional = ast.NewZeroOrMore(
				ast.NewTreeNode(otherNames, p),
			)
		}
	}
	return additional, nil
}
