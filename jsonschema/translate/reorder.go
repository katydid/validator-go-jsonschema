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
	"slices"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func newAnd(ps ...*ast.Pattern) *ast.Pattern {
	ps = slices.DeleteFunc(ps, func(p *ast.Pattern) bool {
		return p.ZAny != nil
	})
	if len(ps) == 0 {
		return ast.NewZAny()
	}
	return ast.NewAnd(ps...)
}

func newOr(ps ...*ast.Pattern) *ast.Pattern {
	if len(ps) == 0 {
		return ast.NewZAny()
	}
	return ast.NewOr(ps...)
}

func newXor(ps ...*ast.Pattern) *ast.Pattern {
	if len(ps) == 0 {
		return ast.NewZAny()
	}
	return ast.NewXor(ps...)
}

func andExpr(list []*ast.Expr) *ast.Expr {
	return std.MustFoldA(list, combinator.And)
}

func newInterleave(ps ...*ast.Pattern) *ast.Pattern {
	ps = slices.DeleteFunc(ps, func(p *ast.Pattern) bool {
		return p.Empty != nil
	})
	if len(ps) == 0 {
		return ast.NewEmpty()
	}
	return ast.NewInterleave(ps...)
}
