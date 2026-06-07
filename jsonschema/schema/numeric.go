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

import (
	"errors"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

// http://json-schema.org/latest/json-schema-validation.html#anchor13
type Numeric struct {
	MultipleOf       *float64  `json:"multipleOf,omitempty"`
	Maximum          *float64  `json:"maximum,omitempty"`
	ExclusiveMaximum *Exlusive `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64  `json:"minimum,omitempty"`
	ExclusiveMinimum *Exlusive `json:"exclusiveMinimum,omitempty"`
}

func (this Numeric) HasNumericConstraints() bool {
	return this.MultipleOf != nil || this.Maximum != nil || this.Minimum != nil || this.ExclusiveMaximum != nil || this.ExclusiveMinimum != nil
}

type Exlusive struct {
	isExclusive bool
	val         *float64
}

func (this *Exlusive) IsExclusive() bool {
	return this != nil && this.isExclusive
}

func (this *Exlusive) GetFloat() *float64 {
	if this == nil {
		return nil
	}
	return this.val
}

func (this *Exlusive) UnmarshalJSON(buf []byte) error {
	var b bool
	errBool := std.UnmarshalJSON(buf, &b)
	if errBool == nil {
		this.isExclusive = b
		return nil
	}
	var f float64
	errFloat := std.UnmarshalJSON(buf, &f)
	if errFloat == nil {
		this.isExclusive = true
		this.val = &f
		return nil
	}
	return errors.New(errBool.Error() + "\n" + errFloat.Error())
}
