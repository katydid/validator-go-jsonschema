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
	// make sure the funcs are registered
	_ "github.com/katydid/validator-go-jsonschema/jsonschema/funcs"

	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func multipleOfExpr(d float64) *ast.Expr {
	return ast.NewFunction("multipleOf", combinator.DoubleConst(d))
}

func minimumExpr(d float64) *ast.Expr {
	return ast.NewFunction("minimum", combinator.DoubleConst(d))
}

func exclusiveMinimumExpr(d float64) *ast.Expr {
	return ast.NewFunction("exclusiveMinimum", combinator.DoubleConst(d))
}

func maximumExpr(d float64) *ast.Expr {
	return ast.NewFunction("maximum", combinator.DoubleConst(d))
}

func exclusiveMaximumExpr(d float64) *ast.Expr {
	return ast.NewFunction("exclusiveMaximum", combinator.DoubleConst(d))
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

func uriExpr() *ast.Expr {
	return ast.NewFunction("uri")
}

func uriReferenceExpr() *ast.Expr {
	return ast.NewFunction("uriReference")
}

func uriTemplateExpr() *ast.Expr {
	return ast.NewFunction("uriTemplate")
}

func periodExpr() *ast.Expr {
	return ast.NewFunction("period")
}

func semverExpr() *ast.Expr {
	return ast.NewFunction("semver")
}

func stringTypeExpr() *ast.Expr {
	return ast.NewFunction("stringType")
}

func anyExpr() *ast.Expr {
	return ast.NewFunction("any")
}

func regexExpr(s string) *ast.Expr {
	return ast.NewFunction("regex", combinator.StringConst(s), combinator.StringVar())
}

func nullTypeExpr() *ast.Expr {
	return ast.NewFunction("null")
}

func numberTypeExpr() *ast.Expr {
	return ast.NewFunction("number")
}

func integerTypeExpr() *ast.Expr {
	return ast.NewFunction("integer")
}

func boolTypeExpr() *ast.Expr {
	return ast.NewFunction("boolType")
}

func maxLengthExpr(d uint64) *ast.Expr {
	return ast.NewFunction("maxLength", combinator.IntConst(int64(d)))
}

func minLengthExpr(d uint64) *ast.Expr {
	return ast.NewFunction("minLength", combinator.IntConst(int64(d)))
}
