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

	jsonschema "github.com/katydid/validator-go-jsonschema/jsonschema/funcs/santhosh-tekuri"
)

// UUID returns whether a string is a valid uuid
func UUID() (funcs.Bool, error) {
	return funcs.TrimBool(&uuid{
		hash: funcs.Hash("uuid"),
	}), nil
}

var _ funcs.Setter = &uuid{}

func (this *uuid) SetValue(v parse.Token) {
	this.Token = v
}

type uuid struct {
	Token parse.Token
	hash  uint64
}

func (this *uuid) HasVariable() bool {
	return true
}

func (this *uuid) ToExpr() *ast.Expr {
	return ast.NewFunction("uuid")
}

func (this *uuid) Eval() (bool, error) {
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
	str := cast.ToString(v)
	err = jsonschema.ValidateUUID(str)
	return err == nil, nil
}

func (this *uuid) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *uuid) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("uuid", UUID)
}
