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

type minimum struct {
	Token parse.Token
	d     float64
	hash  uint64
}

var _ funcs.Setter = &minimum{}

func (this *minimum) SetValue(v parse.Token) {
	this.Token = v
}

func Minimum(d funcs.ConstDouble) (funcs.Bool, error) {
	evaluatedD, err := d.Eval()
	if err != nil {
		return nil, err
	}
	return &minimum{
		d:    evaluatedD,
		hash: funcs.Hash("minimum", d),
	}, nil
}

func (this *minimum) Eval() (bool, error) {
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
	return n >= this.d, nil
}

func (this *minimum) ToExpr() *ast.Expr {
	return ast.NewFunction("minimum", ast.NewDoubleConst(this.d))
}

func (this *minimum) HasVariable() bool {
	return true
}

func (this *minimum) Hash() uint64 {
	return this.hash
}

func (this *minimum) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*minimum); ok {
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
	funcs.Register("minimum", Minimum)
}
