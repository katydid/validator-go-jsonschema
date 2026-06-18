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
	"github.com/katydid/validator-go/validator/combinator"
)

func boolType() *ast.Pattern {
	return combinator.Value(boolTypeExpr())
}

func integerType() *ast.Pattern {
	return combinator.Value(integerTypeExpr())
}

func nullType() *ast.Pattern {
	return combinator.Value(nullTypeExpr())
}

func numberType() *ast.Pattern {
	return combinator.Value(numberTypeExpr())
}

func stringType() *ast.Pattern {
	return combinator.Value(stringTypeExpr())
}

func arrayType() *ast.Pattern {
	return NewArrayNode(ast.NewZAny())
}

func objectType() *ast.Pattern {
	return NewObjectNode(ast.NewZAny())
}

func hasType(typs *schema.Type, theType schema.SimpleType) bool {
	if typs == nil {
		return false
	}
	for _, typ := range *typs {
		if typ == theType {
			return true
		}
	}
	return false
}

func anyFieldType() *ast.Pattern {
	return combinator.Value(anyValueExpr())
}

func notObjectType() *ast.Pattern {
	return newOr(arrayType(), anyFieldType())
}

func notArrayType() *ast.Pattern {
	return newOr(objectType(), anyFieldType())
}

func notStringType() *ast.Pattern {
	return newOr(objectType(), arrayType(), ast.NewNot(stringType()))
}

func notNumberType() *ast.Pattern {
	return newOr(objectType(), arrayType(), ast.NewNot(numberType()))
}
