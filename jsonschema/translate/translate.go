//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package translate

import (
	"fmt"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func TranslateDraft4(schema *schema.Schema) (*ast.Grammar, error) {
	p, err := translate(schema)
	if err != nil {
		return nil, err
	}
	return ast.NewGrammar(ast.RefLookup(map[string]*ast.Pattern{"main": p})), nil
}

func translates(schemas []*schema.Schema) ([]*ast.Pattern, error) {
	return std.MapErr(schemas, translate)
}

func translate(schema *schema.Schema) (*ast.Pattern, error) {
	var ps []*ast.Pattern
	if len(schema.Id) > 0 {
		return nil, fmt.Errorf("TODO: id not supported")
	}
	if schema.Default != nil {
		return nil, fmt.Errorf("TODO: default not supported")
	}
	if schema.Type != nil {
		p, err := translateTypes(*schema.Type)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if schema.HasNumericConstraints() {
		p, err := translateNumeric(schema.Numeric)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if schema.HasStringConstraints() {
		p, err := translateString(schema.String, schema.Format)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if schema.HasArrayConstraints() {
		return nil, fmt.Errorf("TODO: array not supported")
	}
	if schema.HasObjectConstraints() {
		p, err := translateObject(schema)
		if err != nil {
			return nil, err
		}
		if !hasObjectType(schema.Type) {
			p = ast.NewOr(p, notObjectType())
		}
		ps = append(ps, p)
	}
	if schema.HasOperatorConstraints() {
		p, err := translateOperators(schema)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if len(schema.Format) > 0 {
		expr, err := translateFormat(schema.Format)
		if err != nil {
			return nil, err
		}
		p := combinator.Value(expr)
		ps = append(ps, p)
	}
	if len(schema.Ref) > 0 {
		return nil, fmt.Errorf("TODO: ref not supported")
	}
	if len(ps) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewAnd(ps...), nil
}
