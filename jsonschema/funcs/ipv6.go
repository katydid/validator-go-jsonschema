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

	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/ipv6/lexer"
)

// IPv6 returns whether a string is a valid ipv6
func IPv6() (funcs.Bool, error) {
	return funcs.TrimBool(&ipv6{
		hash:  funcs.Hash("ipv6"),
		lexer: lexer.NewLexer(nil),
	}), nil
}

var _ funcs.Setter = &ipv6{}

func (this *ipv6) SetValue(v parse.Token) {
	this.Token = v
}

type ipv6 struct {
	Token parse.Token
	lexer *lexer.Lexer
	hash  uint64
}

func (this *ipv6) HasVariable() bool {
	return true
}

func (this *ipv6) ToExpr() *ast.Expr {
	return ast.NewFunction("ipv6")
}

func (this *ipv6) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	if kind != parse.StringKind {
		// ignore non appropriate kinds
		return true, nil
	}
	return this.lexer.IsValid(v), nil
}

func (this *ipv6) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *ipv6) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("ipv6", IPv6)
}
