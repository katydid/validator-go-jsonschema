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

package jsonschema

import (
	"fmt"

	_ "github.com/katydid/validator-go-jsonschema/jsonschema/funcs"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func multipleOfExpr(d float64) *ast.Expr {
	return ast.NewFunction("multipleOf", combinator.DoubleConst(d))
}

func emailExpr() *ast.Expr {
	return ast.NewFunction("email")
}

func datetimeExpr() *ast.Expr {
	return ast.NewFunction("datetime")
}

func dateExpr() *ast.Expr {
	return ast.NewFunction("date")
}

func formatExpr(format string) (*ast.Expr, error) {
	switch format {
	case "date":
		return dateExpr(), nil
	case "date-time":
		return datetimeExpr(), nil
	case "email":
		return emailExpr(), nil
	}
	return nil, fmt.Errorf("format %s not supported", format)
}
