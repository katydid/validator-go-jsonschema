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
)

func Translate(s *schema.Schema) (*ast.Grammar, error) {
	defs, err := translateDefinitions(s)
	if err != nil {
		return nil, err
	}
	return ast.NewGrammar(ast.RefLookup(defs)), nil
}

func translate(s *schema.Schema) (*ast.Pattern, error) {
	if s.Const != nil {
		// If there is a const no other constraints are necessary.
		return translateConst(*s.Const)
	}
	if s.Default != nil {
		// default works if we do nothing
	}
	ptype, err := translateTypeConstraints(s)
	if err != nil {
		return nil, err
	}
	ps := []*ast.Pattern{ptype}
	if s.HasOperatorConstraints() {
		p, err := translateOperators(s)
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
		ps = append(ps, p)
	}
	return newAnd(ps...), nil
}

func translateTypeConstraints(s *schema.Schema) (*ast.Pattern, error) {
	var ps []*ast.Pattern
	if hasType(s.Type, schema.TypeNull) {
		ps = append(ps, nullType())
	}
	if hasType(s.Type, schema.TypeBoolean) {
		ps = append(ps, boolType())
	}
	if hasType(s.Type, schema.TypeInteger) || hasType(s.Type, schema.TypeNumber) {
		typ := integerType()
		if hasType(s.Type, schema.TypeNumber) {
			typ = numberType()
		}
		if s.HasNumericConstraints() {
			p, err := translateNumeric(s.Numeric)
			if err != nil {
				return nil, err
			}
			ps = append(ps, newAnd(p, typ))
		} else {
			ps = append(ps, typ)
		}
	} else if s.HasNumericConstraints() {
		p, err := translateNumeric(s.Numeric)
		if err != nil {
			return nil, err
		}
		ps = append(ps, newOr(p, notNumberType()))
	}
	if hasType(s.Type, schema.TypeString) {
		typ := stringType()
		if s.HasStringConstraints() {
			p, err := translateString(s.String)
			if err != nil {
				return nil, err
			}
			ps = append(ps, newAnd(typ, p))
		} else {
			ps = append(ps, typ)
		}
	} else if s.HasStringConstraints() {
		p, err := translateString(s.String)
		if err != nil {
			return nil, err
		}
		ps = append(ps, newOr(p, notStringType()))
	}
	if hasType(s.Type, schema.TypeArray) {
		if s.HasArrayConstraints() {
			p, err := translateArray(s)
			if err != nil {
				return nil, err
			}
			p = NewArrayNode(p)
			ps = append(ps, p)
		} else {
			typ := arrayType()
			ps = append(ps, typ)
		}
	} else if s.HasArrayConstraints() {
		p, err := translateArray(s)
		if err != nil {
			return nil, err
		}
		p = NewArrayNode(p)
		ps = append(ps, newOr(p, notArrayType()))
	}
	if hasType(s.Type, schema.TypeObject) {
		if s.HasObjectConstraints() {
			p, err := translateObject(s)
			if err != nil {
				return nil, err
			}
			p = NewObjectNode(p)
			ps = append(ps, p)
		} else {
			typ := objectType()
			ps = append(ps, typ)
		}
	} else if s.HasObjectConstraints() {
		p, err := translateObject(s)
		if err != nil {
			return nil, err
		}
		p = NewObjectNode(p)
		ps = append(ps, newOr(p, notObjectType()))
	}
	if s.Type != nil && len(*s.Type) == 1 {
		// If there is only one type, then there is no options to consider and all constraints and conjunctions.
		return newAnd(ps...), nil
	}
	// If there zero or two types, we can type the union of the constriants.
	return newOr(ps...), nil
}
