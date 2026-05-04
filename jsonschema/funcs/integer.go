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

type integer struct {
	U           funcs.Uint
	I           funcs.Int
	hasVariable bool
	hash        uint64
}

func Integer(U funcs.Uint, I funcs.Int) (funcs.Double, error) {
	return &integer{
		U:           U,
		I:           I,
		hasVariable: U.HasVariable() || I.HasVariable(),
		hash:        funcs.Hash("integer", U, I),
	}, nil
}

func (this *integer) Eval() (float64, error) {
	u, err := this.U.Eval()
	if err == nil {
		return float64(u), nil
	}
	i, err := this.I.Eval()
	if err == nil {
		return float64(i), nil
	}
	return 0, err
}

func (this *integer) ToExpr() *ast.Expr {
	return ast.NewFunction("integer", this.U.ToExpr(), this.I.ToExpr())
}

func (this *integer) HasVariable() bool {
	return this.hasVariable
}

func (this *integer) Hash() uint64 {
	return this.hash
}

func (this *integer) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*integer); ok {
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("integer", Integer)
}
