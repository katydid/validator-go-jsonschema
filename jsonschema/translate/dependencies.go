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
	return newAnd(p1, p2), nil
}

func translateDependentRequired(deps map[string][]string) (*ast.Pattern, error) {
	res := []*ast.Pattern{}
	names := std.SortedKeys(deps)
	for _, name := range names {
		for _, reqName := range deps[name] {
			ifName := ast.NewContains(ast.NewTreeNode(ast.NewStringName(name), ast.NewZAny()))
			thenName := ast.NewContains(ast.NewTreeNode(ast.NewStringName(reqName), ast.NewZAny()))
			elseName := ast.NewNot(ifName)
			res = append(res, newOr(newAnd(ifName, thenName), elseName))
		}
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return newAnd(res...), nil
}

func translateDependentSchemas(deps map[string]*schema.Schema) (*ast.Pattern, error) {
	res := []*ast.Pattern{}
	names := std.SortedKeys(deps)
	for _, name := range names {
		thenPat, err := translate(deps[name])
		if err != nil {
			return nil, err
		}
		ifPat := ast.NewContains(ast.NewTreeNode(ast.NewStringName(name), ast.NewZAny()))
		elsePat := ast.NewNot(ifPat)
		res = append(res, newOr(newAnd(ifPat, thenPat), elsePat))
	}
	if len(res) == 0 {
		return ast.NewZAny(), nil
	}
	return newAnd(res...), nil
}
