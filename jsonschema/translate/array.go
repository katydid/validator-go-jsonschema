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
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func translateArray(s *schema.Schema) (*ast.Pattern, error) {
	constraints := []*ast.Pattern{}
	if s.UniqueItems {
		return nil, fmt.Errorf("uniqueItems are not supported")
	}
	if s.MaxItems != nil {
		constraints = append(constraints, maxItems(int(*s.MaxItems)))
	}
	if s.MinItems > 0 {
		constraints = append(constraints, minItems(int(s.MinItems)))
	}
	additionalItems := ast.NewZAny()
	if s.AdditionalItems != nil {
		if s.AdditionalItems.Bool != nil {
			if !*s.AdditionalItems.Bool {
				additionalItems = ast.NewEmpty()
			}
		}
		if s.AdditionalItems.Schema != nil {
			p, err := translate(s.AdditionalItems.Schema)
			if err != nil {
				return nil, err
			}
			additionalItems = ast.NewZeroOrMore(anyIndex(p))
		}
	}
	if s.Items != nil {
		// TODO: There is a problem here when items are arrays or objects.
		if s.Items.Object != nil {
			sch := s.Items.Object
			pattern, err := translate(sch)
			if err != nil {
				return nil, err
			}
			constraints = append(constraints, ast.NewZeroOrMore(anyIndex(pattern)))
		} else if s.Items.Array != nil {
			schs := s.Items.Array
			patterns, err := std.MapErr(schs, translate)
			if err != nil {
				return nil, err
			}
			patterns = std.Map(patterns, anyIndex)
			patterns = concatCombos(patterns, additionalItems)
			constraints = append(constraints, ast.NewOr(patterns...))
		}
	}
	if len(constraints) == 0 {
		return ast.NewZAny(), nil
	}
	return ast.NewAnd(constraints...), nil
}

func concatCombos(ps []*ast.Pattern, additionalItems *ast.Pattern) []*ast.Pattern {
	if len(ps) == 0 {
		return []*ast.Pattern{ast.NewEmpty()}
	}
	combos := make([]*ast.Pattern, 0, len(ps)+2)
	combos = append(combos, ast.NewEmpty())
	for i := range ps {
		if i == 0 {
			continue
		}
		combos = append(combos, ast.NewConcat(ps[:i]...))
	}
	psz := append(ps, additionalItems)
	combos = append(combos, ast.NewConcat(psz...))
	return combos
}

func anyIndex(p *ast.Pattern) *ast.Pattern {
	return ast.NewTreeNode(ast.NewAnyName(), p)
}

func maxItems(n int) *ast.Pattern {
	ps := make([]*ast.Pattern, n+1)
	// one more than the maxItems
	for i := 0; i < n+1; i++ {
		ps[i] = ast.NewTreeNode(ast.NewAnyName(), ast.NewZAny())
	}
	res := ast.NewConcat(ps...)
	return ast.NewNot(ast.NewConcat(res, ast.NewZAny()))
}

func minItems(n int) *ast.Pattern {
	ps := make([]*ast.Pattern, n)
	for i := 0; i < n; i++ {
		ps[i] = ast.NewTreeNode(ast.NewAnyName(), ast.NewZAny())
	}
	return ast.NewConcat(ast.NewConcat(ps...), ast.NewZAny())
}
