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

// http://json-schema.org/latest/json-schema-validation.html#anchor36
type Array struct {
	AdditionalItems *Additional `json:"additionalItems,omitempty"`
	Items           *Items      `json:"items,omitempty"`
	MaxItems        *uint64     `json:"maxItems,omitempty"`
	MinItems        uint64      `json:"minItems,omitempty"`
	UniqueItems     bool        `json:"uniqueItems,omitempty"`
}

func (this Array) GetAdditionalItems() *Additional {
	if this.AdditionalItems == nil {
		return nil
	}
	return this.AdditionalItems
}

func (this Array) GetItems() *Items {
	if this.Items == nil {
		return nil
	}
	return this.Items
}

func (this Array) HasArrayConstraints() bool {
	return this.AdditionalItems != nil || this.Items != nil ||
		this.MaxItems != nil || this.MinItems > 0 || this.UniqueItems
}
