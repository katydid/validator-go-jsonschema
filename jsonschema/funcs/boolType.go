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
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type boolType struct {
	Token parse.Token
	hash  uint64
}

var _ funcs.Setter = &boolType{}

func (this *boolType) SetValue(v parse.Token) {
	this.Token = v
}

func BoolType() (funcs.Bool, error) {
	return &boolType{
		hash: funcs.Hash("boolType"),
	}, nil
}

func (this *boolType) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, _, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	return kind == parse.TrueKind || kind == parse.FalseKind, nil
}

func (this *boolType) ToExpr() *ast.Expr {
	return ast.NewFunction("boolType")
}

func (this *boolType) HasVariable() bool {
	return true
}

func (this *boolType) Hash() uint64 {
	return this.hash
}

func (this *boolType) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*boolType); ok {
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("boolType", BoolType)
}
