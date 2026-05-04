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
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type number struct {
	U           funcs.Uint
	I           funcs.Int
	D           funcs.Double
	hasVariable bool
	hash        uint64
}

func Number(U funcs.Uint, I funcs.Int, D funcs.Double) (funcs.Double, error) {
	return &number{
		U:           U,
		I:           I,
		D:           D,
		hash:        funcs.Hash("number", U, I, D),
		hasVariable: U.HasVariable() || I.HasVariable() || D.HasVariable(),
	}, nil
}

func (this *number) Eval() (float64, error) {
	u, err := this.U.Eval()
	if err == nil {
		return float64(u), nil
	}
	i, err := this.I.Eval()
	if err == nil {
		return float64(i), nil
	}
	return this.D.Eval()
}

func (this *number) ToExpr() *ast.Expr {
	return ast.NewFunction("number", this.U.ToExpr(), this.I.ToExpr(), this.D.ToExpr())
}

func (this *number) HasVariable() bool {
	return true
}

func (this *number) Hash() uint64 {
	return this.hash
}

func (this *number) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*number); ok {
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("number", Number)
}
