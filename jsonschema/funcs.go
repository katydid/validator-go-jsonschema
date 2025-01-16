// Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package jsonschema

import (
	"fmt"
	"strings"

	"github.com/katydid/validator-go-jsonschema/validator/funcs"
)

type multipleOf struct {
	N           funcs.Double
	d           float64
	hash        uint64
	hasVariable bool
}

func MultipleOf(n funcs.Double, d funcs.ConstDouble) (funcs.Bool, error) {
	evaluatedD, err := d.Eval()
	if err != nil {
		return nil, err
	}
	return &multipleOf{
		N:           n,
		d:           evaluatedD,
		hash:        funcs.Hash("multipleOf", n, d),
		hasVariable: n.HasVariable(),
	}, nil
}

func (this *multipleOf) Eval() (bool, error) {
	n, err := this.N.Eval()
	if err != nil {
		return false, err
	}
	v := n / this.d
	return v == float64(int64(v)) || v == float64(uint64(v)), nil
}

func (this *multipleOf) String() string {
	return "multipleOf(" + this.N.String() + "," + fmt.Sprintf("%v", this.d) + ")"
}

func (this *multipleOf) HasVariable() bool {
	return this.hasVariable
}

func (this *multipleOf) Hash() uint64 {
	return this.hash
}

func (this *multipleOf) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*multipleOf); ok {
		if this.d != other.d {
			if this.d < other.d {
				return -1
			}
			return 1
		}
		if c := this.N.Compare(other.N); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func init() {
	funcs.Register("multipleOf", MultipleOf)
}

type integer struct {
	U    funcs.Uint
	I    funcs.Int
	hash uint64
}

func Integer() (funcs.Double, error) {
	return &integer{
		U:    funcs.UintVar(),
		I:    funcs.IntVar(),
		hash: funcs.Hash("integer"),
	}, nil
}

func (this *integer) Eval() (float64, error) {
	u, err := this.U.Eval()
	if err == nil {
		return float64(u), nil
	}
	i, err := this.I.Eval()
	if err == nil {
		return float64(i), nil
	}
	return 0, err
}

func (this *integer) String() string {
	return "integer"
}

func (this *integer) HasVariable() bool {
	return true
}

func (this *integer) Hash() uint64 {
	return this.hash
}

func (this *integer) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*integer); ok {
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func init() {
	funcs.Register("integer", Integer)
}

type number struct {
	I    funcs.Double
	D    funcs.Double
	hash uint64
}

func Number() (funcs.Double, error) {
	i, err := Integer()
	return &number{
		I:    i,
		D:    funcs.DoubleVar(),
		hash: funcs.Hash("number"),
	}, err
}

func (this *number) Eval() (float64, error) {
	i, err := this.I.Eval()
	if err == nil {
		return i, nil
	}
	return this.D.Eval()
}

func (this *number) String() string {
	return "number"
}

func (this *number) HasVariable() bool {
	return true
}

func (this *number) Hash() uint64 {
	return this.hash
}

func (this *number) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if _, ok := that.(*number); ok {
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func init() {
	funcs.Register("number", Number)
}

type maxLength struct {
	S           funcs.String
	n           int64
	hasVariable bool
	hash        uint64
}

func MaxLength(S funcs.String, N funcs.ConstInt) (funcs.Bool, error) {
	n, err := N.Eval()
	if err != nil {
		return nil, err
	}
	return &maxLength{
		S:           S,
		n:           n,
		hasVariable: S.HasVariable(),
		hash:        funcs.Hash("maxLength", S, N),
	}, nil
}

func (this *maxLength) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	l := int64(0)
	for range s {
		l++
	}
	return l <= this.n, nil
}

func (this *maxLength) String() string {
	return "maxLength(" + this.S.String() + "," + fmt.Sprintf("%d", this.n) + ")"
}

func (this *maxLength) HasVariable() bool {
	return this.hasVariable
}

func (this *maxLength) Hash() uint64 {
	return this.hash
}

func (this *maxLength) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*maxLength); ok {
		if this.n != other.n {
			if this.n < other.n {
				return -1
			}
			return 1
		}
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func init() {
	funcs.Register("maxLength", MaxLength)
}

type minLength struct {
	S           funcs.String
	n           int64
	hasVariable bool
	hash        uint64
}

func MinLength(S funcs.String, N funcs.ConstInt) (funcs.Bool, error) {
	n, err := N.Eval()
	if err != nil {
		return nil, err
	}
	return &minLength{
		S:           S,
		n:           n,
		hasVariable: S.HasVariable(),
		hash:        funcs.Hash("minLength", S, N),
	}, nil
}

func (this *minLength) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	l := int64(0)
	for range s {
		l++
	}
	return l >= this.n, nil
}

func (this *minLength) String() string {
	return "minLength(" + this.S.String() + "," + fmt.Sprintf("%d", this.n) + ")"
}

func (this *minLength) HasVariable() bool {
	return this.hasVariable
}

func (this *minLength) Hash() uint64 {
	return this.hash
}

func (this *minLength) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*minLength); ok {
		if this.n != other.n {
			if this.n < other.n {
				return -1
			}
			return 1
		}
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func init() {
	funcs.Register("minLength", MinLength)
}
