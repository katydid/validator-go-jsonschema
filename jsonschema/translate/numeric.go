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
	list := []*ast.Expr{}
	if schema.MultipleOf != nil {
		list = append(list, multipleOfExpr(*schema.MultipleOf))
	}
	if schema.Maximum != nil {
		if schema.ExclusiveMaximum {
			list = append(list, exclusiveMaximumExpr(*schema.Maximum))
		} else {
			list = append(list, maximumExpr(*schema.Maximum))
		}
	}
	if schema.Minimum != nil {
		if schema.ExclusiveMinimum {
			list = append(list, exclusiveMinimumExpr(*schema.Minimum))
		} else {
			list = append(list, minimumExpr(*schema.Minimum))
		}
	}
	if len(list) == 0 {
		panic("unreachable")
	}
	return combinator.Value(and(list)), nil
}
