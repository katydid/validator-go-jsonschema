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
	"io"

	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"

	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/ipv4/lexer"
)

// IPv4 returns whether a string is a valid ipv4
func IPv4() (funcs.Bool, error) {
	return funcs.TrimBool(&ipv4{
		lexer: lexer.NewLexer([]byte{}),
		hash:  funcs.Hash("ipv4"),
	}), nil
}

var _ funcs.Setter = &ipv4{}

func (this *ipv4) SetValue(v parse.Token) {
	this.Token = v
}

type ipv4 struct {
	Token parse.Token
	lexer *lexer.Lexer
	hash  uint64
}

func (this *ipv4) HasVariable() bool {
	return true
}

func (this *ipv4) ToExpr() *ast.Expr {
	return ast.NewFunction("ipv4")
}

func isIPV4(lexer *lexer.Lexer, data []byte) bool {
	lexer.Init(data)
	_, err := lexer.Next()
	valid := err == nil
	if valid {
		// check that there is only one email address
		_, err = lexer.Next()
		if err != io.EOF {
			return false
		}
	}
	return valid
}

func (this *ipv4) Eval() (bool, error) {
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
	valid := isIPV4(this.lexer, v)
	return valid, nil
}

func (this *ipv4) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *ipv4) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("ipv4", IPv4)
}
