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

package schema

func (s *Schema) Walk(visit func(s *Schema)) {
	visit(s)
	for _, child := range s.Definitions {
		child.Walk(visit)
	}
	if child := s.Array.GetAdditionalItems().GetSchema(); child != nil {
		child.Walk(visit)
	}
	if child := s.Array.GetItems().GetObject(); child != nil {
		child.Walk(visit)
	}
	for _, child := range s.Array.GetItems().GetArray() {
		child.Walk(visit)
	}
	if child := s.Object.GetAdditionalProperties().GetSchema(); child != nil {
		child.Walk(visit)
	}
	for _, child := range s.Object.GetProperties() {
		child.Walk(visit)
	}
	for _, child := range s.Object.GetPatternProperties() {
		child.Walk(visit)
	}
	for _, child := range s.Operators.AllOf {
		child.Walk(visit)
	}
	for _, child := range s.Operators.AnyOf {
		child.Walk(visit)
	}
	for _, child := range s.Operators.OneOf {
		child.Walk(visit)
	}
	if child := s.Operators.Not; child != nil {
		child.Walk(visit)
	}
	if child := s.Operators.If; child != nil {
		child.Walk(visit)
	}
	if child := s.Operators.Then; child != nil {
		child.Walk(visit)
	}
	if child := s.Operators.Else; child != nil {
		child.Walk(visit)
	}
	if deps := s.Operators.Dependencies; deps != nil {
		for _, dep := range *deps {
			if child := dep.Schema; child != nil {
				child.Walk(visit)
			}
		}
	}
	for _, child := range s.Operators.DependentSchemas {
		child.Walk(visit)
	}
}
