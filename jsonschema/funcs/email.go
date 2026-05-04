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
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/pool"
	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/email/lexer"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

// Email returns whether a string is a valid email
func Email(input funcs.String) (funcs.Bool, error) {
	return funcs.TrimBool(&email{
		lexer:       lexer.NewLexer([]byte{}),
		pool:        pool.New(),
		S:           input,
		hash:        funcs.Hash("email", input),
		hasVariable: input.HasVariable(),
	}), nil
}

type email struct {
	lexer       *lexer.Lexer
	pool        pool.Pool
	S           funcs.String
	hash        uint64
	hasVariable bool
}

func (this *email) HasVariable() bool {
	return this.hasVariable
}

func (this *email) ToExpr() *ast.Expr {
	return ast.NewFunction("email", this.S.ToExpr())
}

func isEmail(pool pool.Pool, lexer *lexer.Lexer, s string) bool {
	bytes := cast.FromString(s, pool.Alloc)
	lexer.Init(bytes)
	_, err := lexer.Next()
	valid := err == nil
	pool.FreeAll()
	return valid
}

func (this *email) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	return isEmail(this.pool, this.lexer, s), nil
}

func (this *email) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*email); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *email) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("email", Email)
}
