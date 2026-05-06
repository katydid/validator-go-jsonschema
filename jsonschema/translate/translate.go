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

func translate(schema *schema.Schema) (*ast.Pattern, error) {
	pattern, err := translateOne(schema)
	if err != nil {
		return nil, err
	}
	if schema.Type != nil {
		types := *schema.Type
		if len(types) == 1 {
			p, err := translateType(types[0])
			if err != nil {
				return nil, err
			}
			pattern = ast.NewAnd(p, pattern)
		} else {
			ps := make([]*ast.Pattern, len(types))
			for i := range types {
				var err error
				ps[i], err = translateType(types[i])
				if err != nil {
					return nil, err
				}
			}
			ors := ast.NewOr(ps...)
			pattern = ast.NewAnd(ors, pattern)
		}
	}
	return pattern, nil
}

func translateOne(schema *schema.Schema) (*ast.Pattern, error) {
	if len(schema.Id) > 0 {
		return nil, fmt.Errorf("id not supported")
	}
	if schema.Default != nil {
		return nil, fmt.Errorf("default not supported")
	}
	if schema.HasNumericConstraints() {
		p, err := translateNumeric(schema.Numeric)
		return p, err
	}
	if schema.HasStringConstraints() {
		if schema.Type != nil && len(*schema.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with string constraints %#v", schema)
		}
		p, err := translateString(schema.String, schema.Format)
		return p, err
	}
	if schema.HasArrayConstraints() {
		return nil, fmt.Errorf("array not supported")
	}
	if schema.HasObjectConstraints() {
		p, err := translateObject(schema)
		return p, err
	}
	if schema.HasOperatorConstraints() {
		p, err := translateInstance(schema)
		return p, err
	}
	if len(schema.Format) > 0 {
		expr, err := translateFormat(schema.Format)
		if err != nil {
			return nil, err
		}
		return combinator.Value(expr), nil
	}

	if len(schema.Ref) > 0 {
		return nil, fmt.Errorf("ref not supported")
	}
	return ast.NewZAny(), nil
}

func translates(schemas []*schema.Schema) ([]*ast.Pattern, error) {
	return std.MapErr(schemas, translate)
}

func translateInstance(schema *schema.Schema) (*ast.Pattern, error) {
	if len(schema.Definitions) > 0 {
		return nil, fmt.Errorf("definitions not supported")
	}
	if len(schema.Enum) > 0 {
		return nil, fmt.Errorf("enum not supported")
	}
	if len(schema.AllOf) > 0 {
		ps, err := translates(schema.AllOf)
		if err != nil {
			return nil, err
		}
		return ast.NewAnd(ps...), nil
	}
	if len(schema.AnyOf) > 0 {
		ps, err := translates(schema.AnyOf)
		if err != nil {
			return nil, err
		}
		return ast.NewOr(ps...), nil
	}
	if len(schema.OneOf) > 0 {
		return translateOneOf(schema.OneOf)
	}
	if schema.Not != nil {
		p, err := translate(schema.Not)
		if err != nil {
			return nil, err
		}
		return ast.NewNot(p), nil
	}
	panic("unreachable object")
}

func translateType(typ schema.SimpleType) (*ast.Pattern, error) {
	switch typ {
	case schema.TypeArray, schema.TypeObject:
		//This does not distinguish between arrays and objects
		return combinator.Many(combinator.InAny(combinator.Any())), nil
	case schema.TypeBoolean:
		return combinator.Value(boolTypeExpr()), nil
	case schema.TypeInteger:
		return combinator.Value(integerTypeExpr()), nil
	case schema.TypeNull:
		return combinator.Value(nullTypeExpr()), nil
	case schema.TypeNumber:
		return combinator.Value(numberTypeExpr()), nil
	case schema.TypeString:
		return combinator.Value(stringTypeExpr()), nil
	}
	panic(fmt.Sprintf("unknown simpletype: %s", typ))
}

func and(list []*ast.Expr) *ast.Expr {
	return std.MustFoldA(list, combinator.And)
}
