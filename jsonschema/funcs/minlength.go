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

type minLength struct {
	S           funcs.String
	n           int64
	hasVariable bool
	hash        uint64
}

func MinLength(S funcs.String, N funcs.ConstInt) (funcs.Bool, error) {
	n, err := N.Eval()
	if err != nil {
		return nil, err
	}
	return &minLength{
		S:           S,
		n:           n,
		hasVariable: S.HasVariable(),
		hash:        funcs.Hash("minLength", S, N),
	}, nil
}

func (this *minLength) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	l := int64(0)
	for range s {
		l++
	}
	return l >= this.n, nil
}

func (this *minLength) ToExpr() *ast.Expr {
	return ast.NewFunction("minLength", this.S.ToExpr(), ast.NewIntConst(this.n))
}

func (this *minLength) HasVariable() bool {
	return this.hasVariable
}

func (this *minLength) Hash() uint64 {
	return this.hash
}

func (this *minLength) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*minLength); ok {
		if this.n != other.n {
			if this.n < other.n {
				return -1
			}
			return 1
		}
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("minLength", MinLength)
}
