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

	if m := schema.ExclusiveMaximum.GetNumber(); m != nil {
		if f := m.GetFloat(); f != nil {
			list = append(list, exclusiveMaximumExpr(*f))
		} else if f := m.GetBigFloat(); f != nil {
			list = append(list, exclusiveMaximumBigExpr(*f))
		}
	} else if schema.Maximum != nil {
		if schema.ExclusiveMaximum.IsExclusive() {
			m := schema.Maximum
			if f := m.GetFloat(); f != nil {
				list = append(list, exclusiveMaximumExpr(*f))
			} else if f := m.GetBigFloat(); f != nil {
				list = append(list, exclusiveMaximumBigExpr(*f))
			}
		} else {
			m := schema.Maximum
			if f := m.GetFloat(); f != nil {
				list = append(list, maximumExpr(*f))
			} else if f := m.GetBigFloat(); f != nil {
				list = append(list, maximumBigExpr(*f))
			}
		}
	}
	if m := schema.ExclusiveMinimum.GetNumber(); m != nil {
		if f := m.GetFloat(); f != nil {
			list = append(list, exclusiveMinimumExpr(*f))
		} else if f := m.GetBigFloat(); f != nil {
			list = append(list, exclusiveMinimumBigExpr(*f))
		}
	} else if schema.Minimum != nil {
		if schema.ExclusiveMinimum.IsExclusive() {
			m := schema.Minimum
			if f := m.GetFloat(); f != nil {
				list = append(list, exclusiveMinimumExpr(*f))
			} else if f := m.GetBigFloat(); f != nil {
				list = append(list, exclusiveMinimumBigExpr(*f))
			}
		} else {
			m := schema.Minimum
			if f := m.GetFloat(); f != nil {
				list = append(list, minimumExpr(*f))
			} else if f := m.GetBigFloat(); f != nil {
				list = append(list, minimumBigExpr(*f))
			}
		}
	}
	return combinator.Value(andExpr(list)), nil
}
