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
	"strings"

	"github.com/dlclark/regexp2/v2"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

type regex struct {
	Token parse.Token
	s     string
	r     *regexp2.Regexp
	hash  uint64
}

var _ funcs.Setter = &regex{}

func (this *regex) SetValue(v parse.Token) {
	this.Token = v
}

func compileRegex(s string) (*regexp2.Regexp, error) {
	return regexp2.Compile(s, regexp2.ECMAScript|regexp2.Unicode)
}

func Regex(S funcs.ConstString) (funcs.Bool, error) {
	s, err := S.Eval()
	if err != nil {
		return nil, err
	}
	r, err := compileRegex(s)
	if err != nil {
		return nil, err
	}
	return &regex{
		s:    s,
		r:    r,
		hash: funcs.Hash("regex", S),
	}, nil
}

func (this *regex) Eval() (bool, error) {
	if this.Token == nil {
		return false, errTokenNotSet
	}
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	if kind != parse.StringKind {
		// ignore non string values.
		return true, nil
	}
	s := cast.ToString(v)
	return this.r.MatchString(s)
}

func (this *regex) ToExpr() *ast.Expr {
	return ast.NewFunction("regex", ast.NewStringConst(this.s))
}

func (this *regex) HasVariable() bool {
	return true
}

func (this *regex) Hash() uint64 {
	return this.hash
}

func (this *regex) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*regex); ok {
		return strings.Compare(this.s, other.s)
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func init() {
	funcs.Register("regexp2", Regex)
}
