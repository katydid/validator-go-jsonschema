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
	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func Translate(s *schema.Schema) (*ast.Grammar, error) {
	defs, err := translateDefinitions(s)
	if err != nil {
		return nil, err
	}
	return ast.NewGrammar(ast.RefLookup(defs)), nil
}

func translate(s *schema.Schema) (*ast.Pattern, error) {
	var ps []*ast.Pattern
	if s.Default != nil {
		// default works if we do nothing
	}
	if s.Type != nil {
		p, err := translateTypes(*s.Type)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if s.HasNumericConstraints() {
		p, err := translateNumeric(s.Numeric)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if s.HasStringConstraints() {
		p, err := translateString(s.String, s.Format)
		if err != nil {
			return nil, err
		}
		if !hasType(s.Type, schema.TypeString) {
			p = ast.NewOr(p, notStringType())
		}
		ps = append(ps, p)
	}
	if s.HasArrayConstraints() {
		p, err := translateArray(s)
		if err != nil {
			return nil, err
		}
		p = NewArrayNode(p)
		if !hasType(s.Type, schema.TypeArray) {
			p = ast.NewOr(p, notArrayType())
		}
		ps = append(ps, p)
	}
	if s.HasObjectConstraints() {
		p, err := translateObject(s)
		if err != nil {
			return nil, err
		}
		p = NewObjectNode(p)
		if !hasType(s.Type, schema.TypeObject) {
			p = ast.NewOr(p, notObjectType())
		}
		ps = append(ps, p)
	}
	if s.HasOperatorConstraints() {
		p, err := translateOperators(s)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if len(s.Format) > 0 {
		expr, err := translateFormat(s.Format)
		if err != nil {
			return nil, err
		}
		p := combinator.Value(expr)
		ps = append(ps, p)
	}
	if s.Const != nil {
		p, err := translateConst(*s.Const)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if len(s.Ref) > 0 {
		prefix := ""
		if len(s.Id) > 0 {
			prefix = s.Id
		}
		p, err := translateRef(prefix, s.Ref)
		if err != nil {
			return nil, err
		}
		// ref overrides others if it is successful
		return p, nil
	}
	if len(ps) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewAnd(ps...), nil
}
