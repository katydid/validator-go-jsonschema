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

type minmaxLength struct {
	Token parse.Token
	min   int
	max   int
	hash  uint64
}

var _ funcs.Setter = &minmaxLength{}

func (this *minmaxLength) SetValue(v parse.Token) {
	this.Token = v
}

func MinMaxLength(Min funcs.ConstInt, Max funcs.ConstInt) (funcs.Bool, error) {
	min, err := Min.Eval()
	if err != nil {
		return nil, err
	}
	max, err := Max.Eval()
	if err != nil {
		return nil, err
	}
	return &minmaxLength{
		min:  int(min),
		max:  int(max),
		hash: funcs.Hash("minmaxLength", Min, Max),
	}, nil
}

func (this *minmaxLength) Eval() (bool, error) {
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
	return runeCountRange(v, this.min, this.max), nil
}

// returns if number of runes is in range
func runeCountRange(bs []byte, min int, max int) bool {
	np := len(bs)
	var n int
	// can only create less characters than bytes, so if already less then we are done
	if len(bs) < min {
		return false
	}
	for ; n < np; n++ {
		if c := bs[n]; c >= utf8.RuneSelf {
			// non-ASCII slow path
			s := cast.ToString(bs[n:])
			return runeCountRangeString(s, min-n, max-n)
		}
		if n > max {
			return false
		}
	}
	return n >= min && n <= max
}

func runeCountRangeString(s string, min int, max int) bool {
	n := 0
	for range s {
		if n > max {
			return false
		}
		n++
	}
	return n >= min && n <= max
}

func (this *minmaxLength) ToExpr() *ast.Expr {
	return ast.NewFunction("minmaxLength", ast.NewIntConst(int64(this.min)), ast.NewIntConst(int64(this.max)))
}

func (this *minmaxLength) HasVariable() bool {
	return true
}

func (this *minmaxLength) Hash() uint64 {
	return this.hash
}

func (this *minmaxLength) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*minmaxLength); ok {
		if this.min != other.min {
			if this.min < other.min {
				return -1
			}
			return 1
		}
		if this.max != other.max {
			if this.max < other.max {
				return -1
			}
			return 1
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("minmaxLength", MinMaxLength)
}
