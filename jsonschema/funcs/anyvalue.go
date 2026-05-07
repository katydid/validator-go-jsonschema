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

// AnyValue always returns true for non tag values.
func AnyValue() (funcs.Bool, error) {
	return funcs.TrimBool(&anyValue{
		hash: funcs.Hash("anyValue"),
	}), nil
}

var _ funcs.Setter = &anyValue{}

func (this *anyValue) SetValue(v parse.Token) {
	this.Token = v
}

type anyValue struct {
	Token parse.Token
	hash  uint64
}

func (this *anyValue) HasVariable() bool {
	return true
}

func (this *anyValue) ToExpr() *ast.Expr {
	return ast.NewFunction("anyValue")
}

func (this *anyValue) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, _, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	return kind != parse.TagKind, nil
}

func (this *anyValue) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *anyValue) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("anyValue", AnyValue)
}
