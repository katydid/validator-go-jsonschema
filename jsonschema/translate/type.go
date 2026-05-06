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
	"github.com/katydid/validator-go/validator/combinator"
)

func translateTypes(typs []schema.SimpleType) (*ast.Pattern, error) {
	ps, err := std.MapErr(typs, translateType)
	if err != nil {
		return nil, err
	}
	return ast.NewOr(ps...), nil
}

func translateType(typ schema.SimpleType) (*ast.Pattern, error) {
	switch typ {
	case schema.TypeObject:
		return objectType(), nil
	case schema.TypeArray:
		return arrayType(), nil
	case schema.TypeBoolean:
		return combinator.Value(boolTypeExpr()), nil
	case schema.TypeInteger:
		return combinator.Value(integerTypeExpr()), nil
	case schema.TypeNull:
		return combinator.Value(nullTypeExpr()), nil
	case schema.TypeNumber:
		return combinator.Value(numberTypeExpr()), nil
	case schema.TypeString:
		return combinator.Value(stringTypeExpr()), nil
	}
	panic(fmt.Sprintf("unknown simpletype: %s", typ))
}

func hasType(typs *schema.Type, typ schema.SimpleType) bool {
	if typs == nil {
		return false
	}
	for _, typ := range *typs {
		if typ == schema.TypeObject {
			return true
		}
	}
	return false
}

func arrayType() *ast.Pattern {
	return NewArrayNode(ast.NewZAny())
}

func objectType() *ast.Pattern {
	return NewObjectNode(ast.NewZAny())
}

func anyFieldType() *ast.Pattern {
	return combinator.Value(anyExpr())
}

func notObjectType() *ast.Pattern {
	return ast.NewOr(arrayType(), anyFieldType())
}

func notArrayType() *ast.Pattern {
	return ast.NewOr(objectType(), anyFieldType())
}

func stringType() *ast.Pattern {
	return combinator.Value(stringTypeExpr())
}

func notStringType() *ast.Pattern {
	return ast.NewOr(objectType(), arrayType(), ast.NewNot(stringType()))
}
