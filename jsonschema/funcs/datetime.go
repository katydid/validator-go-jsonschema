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
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/funcs"

	jsonschema "github.com/katydid/validator-go-jsonschema/jsonschema/funcs/santhosh-tekuri"
)

// DateTime returns whether a string is a valid datetime
func DateTime(input funcs.String) (funcs.Bool, error) {
	return funcs.TrimBool(&datetime{
		S:           input,
		hash:        funcs.Hash("datetime", input),
		hasVariable: input.HasVariable(),
	}), nil
}

type datetime struct {
	S           funcs.String
	hash        uint64
	hasVariable bool
}

func (this *datetime) HasVariable() bool {
	return this.hasVariable
}

func (this *datetime) ToExpr() *ast.Expr {
	return ast.NewFunction("datetime", this.S.ToExpr())
}

func (this *datetime) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	err = jsonschema.ValidateDateTime(s)
	return err == nil, nil
}

func (this *datetime) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*datetime); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *datetime) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("datetime", DateTime)
}
