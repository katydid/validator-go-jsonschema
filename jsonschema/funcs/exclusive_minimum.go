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
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type exclusiveMinimum struct {
	Token parse.Token
	d     float64
	hash  uint64
}

var _ funcs.Setter = &exclusiveMinimum{}

func (this *exclusiveMinimum) SetValue(v parse.Token) {
	this.Token = v
}

func ExclusiveMinimum(d funcs.ConstDouble) (funcs.Bool, error) {
	evaluatedD, err := d.Eval()
	if err != nil {
		return nil, err
	}
	return &exclusiveMinimum{
		d:    evaluatedD,
		hash: funcs.Hash("exclusiveMinimum", d),
	}, nil
}

func (this *exclusiveMinimum) Eval() (bool, error) {
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
		// TODO: Consider not supporting Int64Kind here
		n = float64(cast.ToInt64(v))
	case parse.Float64Kind:
		n = cast.ToFloat64(v)
	default:
		// not a number is ignored
		return true, nil
	}
	return n > this.d, nil
}

func (this *exclusiveMinimum) ToExpr() *ast.Expr {
	return ast.NewFunction("exclusiveMinimum", ast.NewDoubleConst(this.d))
}

func (this *exclusiveMinimum) HasVariable() bool {
	return true
}

func (this *exclusiveMinimum) Hash() uint64 {
	return this.hash
}

func (this *exclusiveMinimum) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*exclusiveMinimum); ok {
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
	funcs.Register("exclusiveMinimum", ExclusiveMinimum)
}
