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
	"encoding/json"
	"fmt"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func exactMatch(a any) (*ast.Pattern, error) {
	if a == nil {
		return combinator.Value(nullTypeExpr()), nil
	}
	switch v := a.(type) {
	case bool:
		return combinator.Value(ast.NewEqual(ast.NewBoolConst(v))), nil
	case string:
		return combinator.Value(ast.NewEqual(ast.NewStringConst(v))), nil
	case json.Number:
		i, interr := v.Int64()
		if interr == nil {
			return combinator.Value(ast.NewEqual(ast.NewIntConst(i))), nil
		}
		f, floaterr := v.Float64()
		if floaterr == nil {
			return combinator.Value(ast.NewEqual(ast.NewDoubleConst(f))), nil
		}
		return nil, fmt.Errorf("unsupported type %T for value %v given errs: [%v, %v]", a, a, interr, floaterr)
	case []any:
		ps := make([]*ast.Pattern, 0, len(v))
		for _, vv := range v {
			p, err := exactMatch(vv)
			if err != nil {
				return nil, err
			}
			ps = append(ps, ast.NewTreeNode(ast.NewAnyName(), p))
		}
		if len(ps) == 0 {
			return NewArrayNode(ast.NewEmpty()), nil
		}
		return NewArrayNode(ast.NewConcat(ps...)), nil
	case map[string]any:
		ps := make([]*ast.Pattern, 0, len(v))
		for k, vv := range v {
			p, err := exactMatch(vv)
			if err != nil {
				return nil, err
			}
			ps = append(ps, ast.NewTreeNode(ast.NewStringName(k), p))
		}
		if len(ps) == 0 {
			return NewObjectNode(ast.NewEmpty()), nil
		}
		return NewObjectNode(ast.NewInterleave(ps...)), nil
	}
	return nil, fmt.Errorf("unsupported type %T for value %v", a, a)
}

func translateEnum(enum []any) (*ast.Pattern, error) {
	if len(enum) == 0 {
		return ast.NewNot(ast.NewZAny()), nil
	}
	if len(enum) > 1 {
		if p := tryAllStrings(enum); p != nil {
			return p, nil
		}
		if p := tryAllNumbers(enum); p != nil {
			return p, nil
		}
	}
	exacts, err := std.MapErr(enum, exactMatch)
	if err != nil {
		return nil, err
	}
	return newOr(exacts...), nil
}

func tryAllStrings(enum []any) *ast.Pattern {
	enums := make([]string, len(enum))
	for i, e := range enum {
		if _, ok := e.(string); !ok {
			return nil
		}
		enums[i] = enum[i].(string)
	}
	return combinator.Value(enumStringExpr(enums))
}

func tryAllNumbers(enum []any) *ast.Pattern {
	enums := make([]float64, len(enum))
	for i, e := range enum {
		if _, ok := e.(json.Number); !ok {
			return nil
		}
		var err error
		enums[i], err = enum[i].(json.Number).Float64()
		if err != nil {
			return nil
		}
	}
	return combinator.Value(enumDoubleExpr(enums))
}
