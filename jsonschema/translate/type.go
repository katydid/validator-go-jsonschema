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
	"github.com/katydid/validator-go/validator/combinator"
)

func translateType(typ schema.SimpleType) (*ast.Pattern, error) {
	switch typ {
	case schema.TypeArray, schema.TypeObject:
		//TODO: This does not distinguish between arrays and objects
		return combinator.Many(combinator.InAny(combinator.Any())), nil
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
