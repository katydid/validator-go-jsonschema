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
	"errors"
	"strings"

	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/regexformat"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

var errRegexVar = errors.New("regex requires a constant expression as its first parameter, but it has a variable parameter")

var errRegexNotVar = errors.New("regex requires a variable expression as its second parameter, but it has a variable parameter")

// Regex returns a new regex function given the first parameter as the expression string that needs to compiled and the second as the regex that should be matched.
func Regex(pattern funcs.ConstString, input funcs.String) (funcs.Bool, error) {
	if pattern.HasVariable() {
		return nil, errRegexVar
	}
	if !input.HasVariable() {
		return nil, errRegexNotVar
	}
	p, err := pattern.Eval()
	if err != nil {
		return nil, err
	}
	matchString, err := regexformat.Compile(p)
	if err != nil {
		return nil, err
	}
	return funcs.TrimBool(&regex{
		pattern:     p,
		matchString: matchString,
		S:           input,
		hash:        funcs.Hash("regex", pattern, input),
		hasVariable: input.HasVariable(),
	}), nil
}

type regex struct {
	pattern     string
	matchString func(string) bool
	S           funcs.String
	hash        uint64
	hasVariable bool
}

func (this *regex) HasVariable() bool {
	return this.hasVariable
}

func (this *regex) ToExpr() *ast.Expr {
	return ast.NewFunction("regex", ast.NewStringConst(this.pattern), this.S.ToExpr())
}

func (this *regex) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		// this.S is always a varString in the case of jsonschema
		// This means an error will only be returned if it is not a string type
		// But json schema says regex has to ignore non string types
		return true, nil
	}
	return this.matchString(s), nil
}

func (this *regex) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*regex); ok {
		if c := strings.Compare(this.pattern, other.pattern); c != 0 {
			return c
		}
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *regex) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("regex", Regex)
}
