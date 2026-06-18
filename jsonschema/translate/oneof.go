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

func translateOneOf(schemas []*schema.Schema) (*ast.Pattern, error) {
	ps, err := std.MapErr(schemas, translate)
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, fmt.Errorf("oneof of zero schemas not supported")
	}
	if len(ps) == 1 {
		return ps[0], nil
	}
	return newXor(ps...), nil
}
