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
	"katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/validator-go/validator/ast"
	"katydid.org.za/go/validator-go/validator/funcs"
)

// EnumDouble returns a function that checks whether the element is contained in the list.
func EnumDouble(list funcs.ConstDoubles) (funcs.Bool, error) {
	if list.HasVariable() {
		return nil, funcs.ErrContainsListNotConst{}
	}
	l, err := list.Eval()
	if err != nil {
		return nil, err
	}
	set := make(map[uint64]struct{})
	for i := range l {
		set[uint64(l[i])] = struct{}{}
	}
	return funcs.TrimBool(&inSetDouble{
		List:        list,
		set:         set,
		hash:        funcs.Hash("enum", list),
		hasVariable: true,
	}), nil
}

var _ funcs.Setter = &inSetDouble{}

func (this *inSetDouble) SetValue(v parse.Token) {
	this.Token = v
}

type inSetDouble struct {
	Token       parse.Token
	List        funcs.ConstDoubles
	set         map[uint64]struct{}
	hash        uint64
	hasVariable bool
}

func (this *inSetDouble) Eval() (bool, error) {
	kind, v, err := this.Token.Token()
	if err != nil {
		return false, err
	}
	switch kind {
	case parse.Float64Kind:
		var u uint64
		cast.ToFloat64BitsPtr(v, &u)
		_, ok := this.set[u]
		return ok, nil
	case parse.Int64Kind:
		var i int64
		cast.ToInt64Ptr(v, &i)
		_, ok := this.set[uint64(i)]
		return ok, nil
	}
	return false, nil
}

func (this *inSetDouble) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*inSetDouble); ok {
		if c := this.List.Compare(other.List); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *inSetDouble) ToExpr() *ast.Expr {
	return ast.NewFunction("enum", this.List.ToExpr())
}

func (this *inSetDouble) HasVariable() bool {
	return this.hasVariable
}

func (this *inSetDouble) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("enum", EnumDouble)
}
