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

func translateArray(s *schema.Schema) (*ast.Pattern, error) {
	if s.Type != nil {
		if len(*s.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with array constraints %#v", s)
		}
		if s.GetType()[0] != schema.TypeArray {
			return nil, fmt.Errorf("%v not supported with array constraints", s.GetType()[0])
		}
	}
	if s.UniqueItems {
		return nil, fmt.Errorf("uniqueItems are not supported")
	}
	if s.MaxItems != nil {
		return nil, fmt.Errorf("maxItems are not supported")
	}
	if s.MinItems > 0 {
		return nil, fmt.Errorf("minItems are not supported")
	}
	additionalItems := true
	if s.AdditionalItems != nil {
		if s.Items == nil {
			//any
		}
		if s.AdditionalItems.Bool != nil {
			additionalItems = *s.AdditionalItems.Bool
		}
		if !additionalItems && (s.MaxLength != nil || s.MinLength > 0) {
			return nil, fmt.Errorf("additionalItems: false and (maxItems|minItems) are not supported together")
		}
		return nil, fmt.Errorf("additionalItems are not supported")
	}
	if s.Items != nil {
		if s.Items.Object != nil {
			if s.Items.Object.Type == nil {
				//any
			} else {
				typ := s.Items.Object.GetType()[0]
				_ = typ
			}
			//TODO this specifies the type of every item in the list
		} else if s.Items.Array != nil {
			if !additionalItems {
				//TODO this specifies the length of the list as well as each ordered element's type
				//  if no type is set then any type is accepted
				maxLength := len(s.Items.Array)
				_ = maxLength
			} else {
				//TODO this specifies the types of the first few ordered items in the list
				//  if no type is set then any type is accepted
			}

		}
		return nil, fmt.Errorf("items are not supported")
	}
	return nil, nil
}
