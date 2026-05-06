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

type stringType struct {
	Token parse.Token
	hash  uint64
}

var _ funcs.Setter = &stringType{}

func (this *stringType) SetValue(v parse.Token) {
	this.Token = v
}

func StringType() (funcs.Bool, error) {
	return &stringType{
		hash: funcs.Hash("stringType"),
	}, nil
}

func (this *stringType) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, _, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	return kind == parse.StringKind, nil
}

func (this *stringType) ToExpr() *ast.Expr {
	return ast.NewFunction("stringType")
}

func (this *stringType) HasVariable() bool {
	return true
}

func (this *stringType) Hash() uint64 {
	return this.hash
}

func (this *stringType) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*stringType); ok {
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("stringType", StringType)
}
