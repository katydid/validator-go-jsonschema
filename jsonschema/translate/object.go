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
	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func translateObject(s *schema.Schema) (*ast.Pattern, error) {
	var constraints []*ast.Pattern
	if s.MaxProperties != nil {
		constraints = append(constraints, maxProperties(int(*s.MaxProperties)))
	}
	if s.MinProperties > 0 {
		constraints = append(constraints, minProperties(int(s.MinProperties)))
	}
	if len(s.Required) > 0 {
		required, err := translateRequired(s.Required)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, required)
	}
	if s.Dependencies != nil {
		deps, err := translateDependencies(s.Dependencies)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, deps)
	}
	if s.DependentRequired != nil {
		deps, err := translateDependentRequired(s.DependentRequired)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, deps)
	}
	if s.DependentSchemas != nil {
		deps, err := translateDependentSchemas(s.DependentSchemas)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, deps)
	}

	additional, err := translateAdditionalProperties(s)
	if err != nil {
		return nil, err
	}

	props, err := newProperties(s)
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

	return ast.NewAnd(constraints...), nil
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
		res = append(res, ast.NewContains(ast.NewTreeNode(ast.NewStringName(req), ast.NewZAny())))
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewAnd(res...), nil
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

func translateAdditionalProperties(s *schema.Schema) (*ast.Pattern, error) {
	additional := ast.NewZAny()

	names := std.SortedKeys(s.Properties)
	patternNames := std.SortedKeys(s.PatternProperties)

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
		} else if s.AdditionalProperties.Schema != nil {
			p, err := translate(s.AdditionalProperties.Schema)
			if err != nil {
				return nil, err
			}
			additional = ast.NewZeroOrMore(
				ast.NewTreeNode(ast.NewAnyName(), p),
			)
		}
	}
	return additional, nil
}

func translateDependencies(deps *schema.Dependencies) (*ast.Pattern, error) {
	d := *deps
	dependentRequired := make(map[string][]string)
	dependentSchemas := make(map[string]*schema.Schema)
	for name := range d {
		if len(d[name].RequiredProperty) > 0 {
			dependentRequired[name] = d[name].RequiredProperty
		} else if d[name].Schema != nil {
			dependentSchemas[name] = d[name].Schema
		}
	}
	p1, err := translateDependentRequired(dependentRequired)
	if err != nil {
		return nil, err
	}
	p2, err := translateDependentSchemas(dependentSchemas)
	if err != nil {
		return nil, err
	}
	return ast.NewAnd(p1, p2), nil
}

func translateDependentRequired(deps map[string][]string) (*ast.Pattern, error) {
	res := []*ast.Pattern{}
	names := std.SortedKeys(deps)
	for _, name := range names {
		for _, reqName := range deps[name] {
			ifName := ast.NewContains(ast.NewTreeNode(ast.NewStringName(name), ast.NewZAny()))
			thenName := ast.NewContains(ast.NewTreeNode(ast.NewStringName(reqName), ast.NewZAny()))
			elseName := ast.NewZeroOrMore(ast.NewTreeNode(ast.NewAnyNameExcept(ast.NewStringName(name)), ast.NewZAny()))
			res = append(res, ast.NewOr(ast.NewAnd(ifName, thenName), elseName))
		}
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewAnd(res...), nil
}

// TODO: There is a problem there that the object and array tags are added, but then they nest object and arrays too deeply.
func translateDependentSchemas(deps map[string]*schema.Schema) (*ast.Pattern, error) {
	res := []*ast.Pattern{}
	names := std.SortedKeys(deps)
	for _, name := range names {
		thenPattern, err := translate(deps[name])
		if err != nil {
			return nil, err
		}
		ifName := ast.NewContains(ast.NewTreeNode(ast.NewStringName(name), ast.NewZAny()))
		elseName := ast.NewZeroOrMore(ast.NewTreeNode(ast.NewAnyNameExcept(ast.NewStringName(name)), ast.NewZAny()))
		res = append(res, ast.NewOr(ast.NewAnd(ifName, thenPattern), elseName))
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewAnd(res...), nil
}
