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
	"bytes"

	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"

	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/hostname/lexer"
)

// Hostname returns whether a string is a valid hostname
func Hostname() (funcs.Bool, error) {
	return funcs.TrimBool(&hostname{
		hash:  funcs.Hash("hostname"),
		lexer: lexer.NewLexer([]byte{}),
	}), nil
}

var _ funcs.Setter = &hostname{}

func (this *hostname) SetValue(v parse.Token) {
	this.Token = v
}

type hostname struct {
	Token parse.Token
	lexer *lexer.Lexer
	hash  uint64
}

func (this *hostname) HasVariable() bool {
	return true
}

func (this *hostname) ToExpr() *ast.Expr {
	return ast.NewFunction("hostname")
}

var dot = []byte(".")

// https://en.wikipedia.org/wiki/Hostname#Restrictions_on_valid_host_names
// > Each label must be 1 to 63 octets long.
// > The entire hostname, including the delimiting dots, has a maximum of 253 ASCII characters.
func isHostname(lexer *lexer.Lexer, data []byte) bool {
	if len(data) >= 253 {
		return false
	}
	valid := lexer.IsValid(data)
	if !valid {
		return false
	}
	index := 0
	for {
		nextIndex := bytes.Index(data[index:], dot)
		if nextIndex < 0 {
			if len(data[index:]) > 63 {
				return false
			}
			return true
		}
		if nextIndex > 63 {
			return false
		}
		index += nextIndex + 1
	}
}

func (this *hostname) Eval() (bool, error) {
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
	return isHostname(this.lexer, v), nil
}

func (this *hostname) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *hostname) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("hostname", Hostname)
}
