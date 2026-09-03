// Copyright 2026 Walter Schulze
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
	"math"
	"math/big"

	"katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/validator-go/validator/ast"
	"katydid.org.za/go/validator-go/validator/funcs"
)

type exclusiveMaximum struct {
	Token parse.Token
	d     float64
	big   *big.Float
	hash  uint64
}

var _ funcs.Setter = &exclusiveMaximum{}

func (this *exclusiveMaximum) SetValue(v parse.Token) {
	this.Token = v
}

func ExclusiveMaximum(d funcs.ConstDouble) (funcs.Bool, error) {
	evaluatedD, err := d.Eval()
	if err != nil {
		return nil, err
	}
	return &exclusiveMaximum{
		d:    evaluatedD,
		big:  big.NewFloat(evaluatedD),
		hash: funcs.Hash("exclusiveMaximum", d),
	}, nil
}

func (this *exclusiveMaximum) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	var n float64
	switch kind {
	case parse.Int64Kind:
		var i int64
		cast.ToInt64Ptr(v, &i)
		n = float64(i)
	case parse.Float64Kind:
		var u uint64
		cast.ToFloat64BitsPtr(v, &u)
		n = math.Float64frombits(u)
	case parse.DecimalKind:
		var s string
		cast.ToStringPtr(v, &s)
		n, _, err := new(big.Float).Parse(s, 10)
		if err != nil {
			return false, nil
		}
		return n.Cmp(this.big) < 0, nil
	default:
		// not a number is ignored
		return true, nil
	}
	return n < this.d, nil
}

func (this *exclusiveMaximum) ToExpr() *ast.Expr {
	return ast.NewFunction("exclusiveMaximum", ast.NewDoubleConst(this.d))
}

func (this *exclusiveMaximum) HasVariable() bool {
	return true
}

func (this *exclusiveMaximum) Hash() uint64 {
	return this.hash
}

func (this *exclusiveMaximum) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*exclusiveMaximum); ok {
		if this.d != other.d {
			if this.d < other.d {
				return -1
			}
			return 1
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("exclusiveMaximum", ExclusiveMaximum)
}
