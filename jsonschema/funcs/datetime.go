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
	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/datetimeformat"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

// DateTime returns whether a string is a valid datetime
func DateTime() (funcs.Bool, error) {
	return funcs.TrimBool(&datetime{
		hash: funcs.Hash("datetime"),
	}), nil
}

var _ funcs.Setter = &datetime{}

func (this *datetime) SetValue(v parse.Token) {
	this.Token = v
}

type datetime struct {
	Token parse.Token
	hash  uint64
}

func (this *datetime) HasVariable() bool {
	return true
}

func (this *datetime) ToExpr() *ast.Expr {
	return ast.NewFunction("datetime")
}

func (this *datetime) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	if kind != parse.StringKind && kind != parse.DateTimeKind {
		// ignore non appropriate kinds
		return true, nil
	}
	return datetimeformat.IsValid(v), nil
}

func (this *datetime) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *datetime) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("datetime", DateTime)
}
