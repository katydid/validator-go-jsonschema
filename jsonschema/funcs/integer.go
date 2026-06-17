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
	"math/big"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type integer struct {
	Token parse.Token
	big   big.Int
	hash  uint64
}

var _ funcs.Setter = &integer{}

func (this *integer) SetValue(v parse.Token) {
	this.Token = v
}

func Integer() (funcs.Bool, error) {
	return &integer{
		hash: funcs.Hash("integer"),
	}, nil
}

func (this *integer) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	if kind == parse.Int64Kind {
		return true, nil
	}
	if kind != parse.DecimalKind {
		return false, nil
	}
	s := cast.ToString(v)
	_, ok := this.big.SetString(s, 10)
	return ok, nil
}

func (this *integer) ToExpr() *ast.Expr {
	return ast.NewFunction("integer")
}

func (this *integer) HasVariable() bool {
	return true
}

func (this *integer) Hash() uint64 {
	return this.hash
}

func (this *integer) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*integer); ok {
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("integer", Integer)
}
