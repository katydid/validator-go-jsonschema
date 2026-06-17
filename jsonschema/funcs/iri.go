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

	jsonschema "github.com/katydid/validator-go-jsonschema/jsonschema/funcs/ianlancetaylor"
)

// IRI returns whether a string is a valid iri
func IRI() (funcs.Bool, error) {
	return funcs.TrimBool(&iri{
		hash: funcs.Hash("iri"),
	}), nil
}

var _ funcs.Setter = &iri{}

func (this *iri) SetValue(v parse.Token) {
	this.Token = v
}

type iri struct {
	Token parse.Token
	hash  uint64
}

func (this *iri) HasVariable() bool {
	return true
}

func (this *iri) ToExpr() *ast.Expr {
	return ast.NewFunction("iri")
}

func isIRI(str string) bool {
	err := jsonschema.ValidateIRI(str)
	return err == nil
}

func (this *iri) Eval() (bool, error) {
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
	valid := isIRI(str)
	return valid, nil
}

func (this *iri) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	return this.ToExpr().Compare(that.ToExpr())
}

func (this *iri) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("iri", IRI)
}
