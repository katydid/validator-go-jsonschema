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

func hostNameExpr() *ast.Expr {
	return ast.NewFunction("hostname")
}

func jsonPointerExpr() *ast.Expr {
	return ast.NewFunction("jsonPointer")
}

func relativeJSONPointerExpr() *ast.Expr {
	return ast.NewFunction("relativeJSONPointer")
}

func uuidExpr() *ast.Expr {
	return ast.NewFunction("uuid")
}

func durationExpr() *ast.Expr {
	return ast.NewFunction("duration")
}

func ipv4Expr() *ast.Expr {
	return ast.NewFunction("ipv4")
}

func ipv6Expr() *ast.Expr {
	return ast.NewFunction("ipv6")
}

func timeExpr() *ast.Expr {
	return ast.NewFunction("time")
}

func formatExpr(format string) (*ast.Expr, error) {
	switch format {
	case "date":
		return dateExpr(), nil
	case "date-time":
		return datetimeExpr(), nil
	case "email":
		return emailExpr(), nil
	case "hostname":
		return hostNameExpr(), nil
	case "json-pointer":
		return jsonPointerExpr(), nil
	case "relative-json-pointer":
		return relativeJSONPointerExpr(), nil
	case "uuid":
		return uuidExpr(), nil
	case "duration":
		return durationExpr(), nil
	case "ipv4":
		return ipv4Expr(), nil
	case "ipv6":
		return ipv6Expr(), nil
	case "time":
		return timeExpr(), nil
	}
	return nil, fmt.Errorf("format %s not supported", format)
}

func newNumberExpr() *ast.Expr {
	return ast.NewFunction("number")
}

func newIntegerExpr() *ast.Expr {
	return ast.NewFunction("integer")
}

func newTypeExpr(e *ast.Expr) *ast.Expr {
	return ast.NewFunction("type", e)
}

func maxLengthExpr(d uint64) *ast.Expr {
	return ast.NewFunction("maxLength", combinator.IntConst(int64(d)))
}

func minLengthExpr(d uint64) *ast.Expr {
	return ast.NewFunction("minLength", combinator.IntConst(int64(d)))
}
