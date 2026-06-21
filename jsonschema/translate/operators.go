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

func translateOperators(schema *schema.Schema) (*ast.Pattern, error) {
	var res []*ast.Pattern
	if schema.Enum != nil {
		p, err := translateEnum(schema.Enum)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if len(schema.AllOf) > 0 {
		ps, err := std.MapErr(schema.AllOf, translate)
		if err != nil {
			return nil, err
		}
		res = append(res, newAnd(ps...))
	}
	if len(schema.AnyOf) > 0 {
		ps, err := std.MapErr(schema.AnyOf, translate)
		if err != nil {
			return nil, err
		}
		res = append(res, newOr(ps...))
	}
	if len(schema.OneOf) > 0 {
		p, err := translateOneOf(schema.OneOf)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if schema.Not != nil {
		p, err := translate(schema.Not)
		if err != nil {
			return nil, err
		}
		res = append(res, ast.NewNot(p))
	}
	if schema.If != nil {
		p, err := translateIf(schema.If, schema.Then, schema.Else)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if schema.Dependencies != nil {
		deps, err := translateDependencies(schema.Dependencies)
		if err != nil {
			return nil, err
		}
		res = append(res, NewObjectNode(deps))
	}
	if schema.DependentRequired != nil {
		deps, err := translateDependentRequired(schema.DependentRequired)
		if err != nil {
			return nil, err
		}
		res = append(res, NewObjectNode(deps))
	}
	if schema.DependentSchemas != nil {
		deps, err := translateDependentSchemas(schema.DependentSchemas)
		if err != nil {
			return nil, err
		}
		res = append(res, NewObjectNode(deps))
	}
	return newAnd(res...), nil
}
