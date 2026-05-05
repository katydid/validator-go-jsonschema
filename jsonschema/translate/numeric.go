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
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func translateNumeric(schema schema.Numeric) (*ast.Pattern, error) {
	v := newNumberExpr()
	list := []*ast.Expr{}
	notNum := combinator.Not(newTypeExpr(newNumberExpr()))
	if schema.MultipleOf != nil {
		mult := multipleOfExpr(*schema.MultipleOf)
		list = append(list, mult)
	}
	if schema.Maximum != nil {
		lt := combinator.LE(v, combinator.DoubleConst(*schema.Maximum))
		if schema.ExclusiveMaximum {
			lt = combinator.LT(v, combinator.DoubleConst(*schema.Maximum))
		}
		list = append(list, combinator.Or(lt, notNum))
	}
	if schema.Minimum != nil {
		lt := combinator.GE(v, combinator.DoubleConst(*schema.Minimum))
		if schema.ExclusiveMinimum {
			lt = combinator.GT(v, combinator.DoubleConst(*schema.Minimum))
		}
		list = append(list, combinator.Or(lt, notNum))
	}
	if len(list) == 0 {
		return combinator.Value(newTypeExpr(v)), nil
	}
	return combinator.Value(and(list)), nil
}
