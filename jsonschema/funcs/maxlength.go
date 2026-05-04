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

type maxLength struct {
	S           funcs.String
	n           int64
	hasVariable bool
	hash        uint64
}

func MaxLength(S funcs.String, N funcs.ConstInt) (funcs.Bool, error) {
	n, err := N.Eval()
	if err != nil {
		return nil, err
	}
	return &maxLength{
		S:           S,
		n:           n,
		hasVariable: S.HasVariable(),
		hash:        funcs.Hash("maxLength", S, N),
	}, nil
}

func (this *maxLength) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	l := int64(0)
	for range s {
		l++
	}
	return l <= this.n, nil
}

func (this *maxLength) ToExpr() *ast.Expr {
	return ast.NewFunction("maxLength", this.S.ToExpr(), ast.NewIntConst(this.n))
}

func (this *maxLength) HasVariable() bool {
	return this.hasVariable
}

func (this *maxLength) Hash() uint64 {
	return this.hash
}

func (this *maxLength) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*maxLength); ok {
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
	funcs.Register("maxLength", MaxLength)
}
