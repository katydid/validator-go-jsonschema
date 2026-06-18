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
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"
)

// EnumString returns a function that checks whether the element is contained in the list.
func EnumString(list funcs.ConstStrings) (funcs.Bool, error) {
	if list.HasVariable() {
		return nil, funcs.ErrContainsListNotConst{}
	}
	l, err := list.Eval()
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{})
	for i := range l {
		set[l[i]] = struct{}{}
	}
	return funcs.TrimBool(&inSetString{
		List:        list,
		set:         set,
		hash:        funcs.Hash("enum", list),
		hasVariable: true,
	}), nil
}

var _ funcs.Setter = &inSetString{}

func (this *inSetString) SetValue(v parse.Token) {
	this.Token = v
}

type inSetString struct {
	Token       parse.Token
	List        funcs.ConstStrings
	set         map[string]struct{}
	hash        uint64
	hasVariable bool
}

func (this *inSetString) Eval() (bool, error) {
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	if kind != parse.StringKind {
		return false, nil
	}
	s := cast.ToString(v)
	_, ok := this.set[s]
	return ok, nil
}

func (this *inSetString) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*inSetString); ok {
		if c := this.List.Compare(other.List); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *inSetString) ToExpr() *ast.Expr {
	return ast.NewFunction("enum", this.List.ToExpr())
}

func (this *inSetString) HasVariable() bool {
	return this.hasVariable
}

func (this *inSetString) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("enum", EnumString)
}
