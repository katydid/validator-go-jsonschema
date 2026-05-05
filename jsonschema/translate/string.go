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

func translateString(schema schema.String, format string) (*ast.Pattern, error) {
	v := combinator.StringVar()
	list := []*ast.Expr{}
	if schema.MaxLength != nil {
		list = append(list, maxLengthExpr(*schema.MaxLength))
	}
	if schema.MinLength > 0 {
		list = append(list, minLengthExpr(schema.MinLength))
	}
	if schema.Pattern != nil {
		list = append(list, regexExpr(*schema.Pattern))
	}
	if len(format) > 0 {
		formatExpr, err := formatExpr(format)
		if err != nil {
			return nil, err
		}
		list = append(list, formatExpr)
	}
	if len(list) == 0 {
		return combinator.Value(newTypeExpr(v)), nil
	}
	return combinator.Value(and(list)), nil
}
