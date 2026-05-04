// Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package funcs

import (
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type multipleOf struct {
	N           funcs.Double
	d           float64
	hash        uint64
	hasVariable bool
}

func MultipleOf(n funcs.Double, d funcs.ConstDouble) (funcs.Bool, error) {
	evaluatedD, err := d.Eval()
	if err != nil {
		return nil, err
	}
	return &multipleOf{
		N:           n,
		d:           evaluatedD,
		hash:        funcs.Hash("multipleOf", n, d),
		hasVariable: n.HasVariable(),
	}, nil
}

func isMultipleOf(n float64, d float64) bool {
	v := n / d
	return v == float64(int64(v)) || v == float64(uint64(v))
}

func (this *multipleOf) Eval() (bool, error) {
	n, err := this.N.Eval()
	if err != nil {
		return false, err
	}
	return isMultipleOf(n, this.d), nil
}

func (this *multipleOf) ToExpr() *ast.Expr {
	return ast.NewFunction("multipleOf", this.N.ToExpr(), ast.NewDoubleConst(this.d))
}

func (this *multipleOf) HasVariable() bool {
	return this.hasVariable
}

func (this *multipleOf) Hash() uint64 {
	return this.hash
}

func (this *multipleOf) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*multipleOf); ok {
		if this.d != other.d {
			if this.d < other.d {
				return -1
			}
			return 1
		}
		if c := this.N.Compare(other.N); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("multipleOf", MultipleOf)
}
