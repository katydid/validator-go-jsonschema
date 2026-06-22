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
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
)

func translateOperators(s *schema.Schema) (*ast.Pattern, error) {
	var res []*ast.Pattern
	if s.Enum != nil {
		p, err := translateEnum(s.Enum)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if len(s.AllOf) > 0 {
		ps, err := std.MapErr(s.AllOf, translateWithParentId(s.Id))
		if err != nil {
			return nil, err
		}
		res = append(res, newAnd(ps...))
	}
	if len(s.AnyOf) > 0 {
		ps, err := std.MapErr(s.AnyOf, translateWithParentId(s.Id))
		if err != nil {
			return nil, err
		}
		res = append(res, newOr(ps...))
	}
	if len(s.OneOf) > 0 {
		p, err := translateOneOf(s.Id, s.OneOf)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if s.Not != nil {
		p, err := translate(s.Id, s.Not)
		if err != nil {
			return nil, err
		}
		res = append(res, ast.NewNot(p))
	}
	if s.If != nil {
		p, err := translateIf(s.Id, s.If, s.Then, s.Else)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if s.Dependencies != nil {
		deps, err := translateDependencies(s.Id, s.Dependencies)
		if err != nil {
			return nil, err
		}
		res = append(res, deps)
	}
	if s.DependentRequired != nil {
		deps, err := translateDependentRequired(s.DependentRequired)
		if err != nil {
			return nil, err
		}
		res = append(res, deps)
	}
	if s.DependentSchemas != nil {
		deps, err := translateDependentSchemas(s.Id, s.DependentSchemas)
		if err != nil {
			return nil, err
		}
		res = append(res, deps)
	}
	return newAnd(res...), nil
}
