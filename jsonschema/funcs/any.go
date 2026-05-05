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

// Any always returns true
func Any() (funcs.Bool, error) {
	return funcs.TrimBool(&any{
		hash: funcs.Hash("any"),
	}), nil
}

var _ funcs.Setter = &any{}

func (this *any) SetValue(v parse.Token) {
	this.Token = v
}

type any struct {
	Token parse.Token
	hash  uint64
}

func (this *any) HasVariable() bool {
	return true
}

func (this *any) ToExpr() *ast.Expr {
	return ast.NewFunction("any")
}

func (this *any) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	return true, nil
}

func (this *any) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *any) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("any", Any)
}
