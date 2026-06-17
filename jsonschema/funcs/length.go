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

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type length struct {
	Token parse.Token
	n     int64
	hash  uint64
}

var _ funcs.Setter = &length{}

func (this *length) SetValue(v parse.Token) {
	this.Token = v
}

func Length(N funcs.ConstInt) (funcs.Bool, error) {
	n, err := N.Eval()
	if err != nil {
		return nil, err
	}
	return &length{
		n:    n,
		hash: funcs.Hash("length", N),
	}, nil
}

func (this *length) Eval() (bool, error) {
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
	expected := int(this.n)
	return runeCountEq(v, expected), nil
}

// returns if number of runes is equal to want
func runeCountEq(bs []byte, want int) bool {
	// there is no way to create more characters from fewer bytes, so the length cannot be equal.
	if len(bs) < want {
		return false
	}
	np := len(bs)
	var n int
	for ; n < np; n++ {
		if c := bs[n]; c >= utf8.RuneSelf {
			// non-ASCII slow path
			s := cast.ToString(bs[n:])
			return runeCountStringEq(s, want-n)
		}
		if n > want {
			return false
		}
	}
	return n == want
}

func runeCountStringEq(s string, expected int) bool {
	n := 0
	for range s {
		if n > expected {
			return false
		}
		n++
	}
	return n == expected
}

func (this *length) ToExpr() *ast.Expr {
	return ast.NewFunction("length", ast.NewIntConst(this.n))
}

func (this *length) HasVariable() bool {
	return true
}

func (this *length) Hash() uint64 {
	return this.hash
}

func (this *length) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*length); ok {
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
	funcs.Register("length", Length)
}
