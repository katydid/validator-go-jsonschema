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

package funcs

import (
	"unicode/utf8"

	"katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/validator-go/validator/ast"
	"katydid.org.za/go/validator-go/validator/funcs"
)

type minLength struct {
	Token parse.Token
	n     int64
	hash  uint64
}

var _ funcs.Setter = &minLength{}

func (this *minLength) SetValue(v parse.Token) {
	this.Token = v
}

func MinLength(N funcs.ConstInt) (funcs.Bool, error) {
	n, err := N.Eval()
	if err != nil {
		return nil, err
	}
	return &minLength{
		n:    n,
		hash: funcs.Hash("minLength", N),
	}, nil
}

func (this *minLength) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	if kind != parse.StringKind {
		// ignore non string values.
		return true, nil
	}
	return runeCountGe(v, int(this.n)), nil
}

// returns if number of runes is greater than of equal to min.
func runeCountGe(bs []byte, min int) bool {
	np := len(bs)
	var n int
	// can only create less characters than bytes, so if already less then we are done
	if len(bs) < min {
		return false
	}
	for ; n < np; n++ {
		if c := bs[n]; c >= utf8.RuneSelf {
			// non-ASCII slow path
			var s string
			cast.ToStringPtr(bs[n:], &s)
			return runeCountStringGe(s, min-n)
		}
		if n >= min {
			return true
		}
	}
	return n >= min
}

func runeCountStringGe(s string, min int) bool {
	n := 0
	for range s {
		if n >= min {
			return true
		}
		n++
	}
	return n >= min
}

func (this *minLength) ToExpr() *ast.Expr {
	return ast.NewFunction("minLength", ast.NewIntConst(this.n))
}

func (this *minLength) HasVariable() bool {
	return true
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
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("minLength", MinLength)
}
