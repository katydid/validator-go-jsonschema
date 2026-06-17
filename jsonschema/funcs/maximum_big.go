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
	"math/big"
	"strings"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type maximumbig struct {
	Token parse.Token
	s     string
	big   *big.Float
	hash  uint64
}

var _ funcs.Setter = &maximumbig{}

func (this *maximumbig) SetValue(v parse.Token) {
	this.Token = v
}

func MaximumBig(s funcs.ConstString) (funcs.Bool, error) {
	evaluatedS, err := s.Eval()
	if err != nil {
		return nil, err
	}
	evaluatedB, _, err := new(big.Float).Parse(evaluatedS, 10)
	if err != nil {
		return nil, err
	}
	return &maximumbig{
		s:    evaluatedS,
		big:  evaluatedB,
		hash: funcs.Hash("maximumbig", s),
	}, nil
}

func (this *maximumbig) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	var n *big.Float
	switch kind {
	case parse.Int64Kind:
		n = big.NewFloat(float64(cast.ToInt64(v)))
	case parse.Float64Kind:
		n = big.NewFloat(cast.ToFloat64(v))
	case parse.DecimalKind:
		s := cast.ToString(v)
		n, _, err = new(big.Float).Parse(s, 10)
		if err != nil {
			return false, nil
		}
	default:
		// not a number is ignored
		return true, nil
	}
	return n.Cmp(this.big) <= 0, nil
}

func (this *maximumbig) ToExpr() *ast.Expr {
	return ast.NewFunction("maximumbig", ast.NewStringConst(this.s))
}

func (this *maximumbig) HasVariable() bool {
	return true
}

func (this *maximumbig) Hash() uint64 {
	return this.hash
}

func (this *maximumbig) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*maximumbig); ok {
		return strings.Compare(this.s, other.s)
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("maximumbig", MaximumBig)
}
