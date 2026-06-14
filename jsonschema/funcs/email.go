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
	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/email/lexer"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

// Email returns whether a string is a valid email
func Email() (funcs.Bool, error) {
	return funcs.TrimBool(&email{
		lexer: lexer.NewLexer([]byte{}),
		hash:  funcs.Hash("email"),
	}), nil
}

var _ funcs.Setter = &email{}

func (this *email) SetValue(v parse.Token) {
	this.Token = v
}

type email struct {
	Token parse.Token
	lexer *lexer.Lexer
	hash  uint64
}

func (this *email) HasVariable() bool {
	return true
}

func (this *email) ToExpr() *ast.Expr {
	return ast.NewFunction("email")
}

func isEmail(lexer *lexer.Lexer, data []byte) bool {
	return lexer.IsValid(data)
}

func (this *email) Eval() (bool, error) {
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
	return isEmail(this.lexer, v), nil
}

func (this *email) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *email) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("email", Email)
}
