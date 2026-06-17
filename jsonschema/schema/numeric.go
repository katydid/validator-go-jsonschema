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
	"math/big"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

// http://json-schema.org/latest/json-schema-validation.html#anchor13
type Numeric struct {
	MultipleOf       *float64  `json:"multipleOf,omitempty"`
	Maximum          *Number   `json:"maximum,omitempty"`
	ExclusiveMaximum *Exlusive `json:"exclusiveMaximum,omitempty"`
	Minimum          *Number   `json:"minimum,omitempty"`
	ExclusiveMinimum *Exlusive `json:"exclusiveMinimum,omitempty"`
}

func (this Numeric) HasNumericConstraints() bool {
	return this.MultipleOf != nil || this.Maximum != nil || this.Minimum != nil || this.ExclusiveMaximum != nil || this.ExclusiveMinimum != nil
}

type Number struct {
	f        *float64
	bigFloat *string
}

func (this *Number) GetFloat() *float64 {
	if this == nil {
		return nil
	}
	return this.f
}

func (this *Number) GetBigFloat() *string {
	if this == nil {
		return nil
	}
	return this.bigFloat
}

func (this *Number) UnmarshalJSON(buf []byte) error {
	var f float64
	errFloat := std.UnmarshalJSON(buf, &f)
	if errFloat == nil {
		this.f = &f
		return nil
	}
	s := string(buf)
	_, _, errBig := new(big.Float).Parse(s, 10)
	if errBig == nil {
		this.bigFloat = &s
		return nil
	}
	return errors.New(errFloat.Error() + "\n" + errBig.Error())
}

type Exlusive struct {
	isExclusive bool
	val         *Number
}

func (this *Exlusive) IsExclusive() bool {
	return this != nil && this.isExclusive
}

func (this *Exlusive) GetNumber() *Number {
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
	var n *Number
	errFloat := std.UnmarshalJSON(buf, &n)
	if errFloat == nil {
		this.isExclusive = true
		this.val = n
		return nil
	}
	return errors.New(errBool.Error() + "\n" + errFloat.Error())
}
