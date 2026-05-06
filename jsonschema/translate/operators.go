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
	"fmt"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go/validator/ast"
)

func translateOperators(schema *schema.Schema) (*ast.Pattern, error) {
	if len(schema.Definitions) > 0 {
		return nil, fmt.Errorf("definitions not supported")
	}
	if len(schema.Enum) > 0 {
		return nil, fmt.Errorf("enum not supported")
	}
	if len(schema.AllOf) > 0 {
		ps, err := translates(schema.AllOf)
		if err != nil {
			return nil, err
		}
		return ast.NewAnd(ps...), nil
	}
	if len(schema.AnyOf) > 0 {
		ps, err := translates(schema.AnyOf)
		if err != nil {
			return nil, err
		}
		return ast.NewOr(ps...), nil
	}
	if len(schema.OneOf) > 0 {
		return translateOneOf(schema.OneOf)
	}
	if schema.Not != nil {
		p, err := translate(schema.Not)
		if err != nil {
			return nil, err
		}
		return ast.NewNot(p), nil
	}
	panic("unreachable object")
}
