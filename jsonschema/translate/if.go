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

package translate

import (
	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go/validator/ast"
)

func translateIf(cnd, thn, els *schema.Schema) (*ast.Pattern, error) {
	cndp, err := translate(cnd)
	if err != nil {
		return nil, err
	}
	thnp := ast.NewZAny()
	if thn != nil {
		thnp, err = translate(thn)
		if err != nil {
			return nil, err
		}
	}
	elsp := ast.NewZAny()
	if els != nil {
		elsp, err = translate(els)
		if err != nil {
			return nil, err
		}
	}
	return newOr(newAnd(cndp, thnp), newAnd(ast.NewNot(cndp), elsp)), nil
}
