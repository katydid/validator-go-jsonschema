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
	"encoding/json"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

type TagType byte

const UnknownTag = TagType(0)

const FieldTag = TagType(1)

const ArrayTag = TagType(2)

const ObjectTag = TagType(3)

func (t TagType) String() string {
	switch t {
	case FieldTag:
		return ""
	case ArrayTag:
		return "array"
	case ObjectTag:
		return "object"
	case UnknownTag:
		return "<unknowntag>"
	}
	panic("unreachable")
}

func GetTags(s *schema.Schema) []TagType {
	tags := map[TagType]struct{}{}
	for _, typ := range s.GetType() {
		tag := getTag(typ)
		if tag != UnknownTag {
			tags[tag] = struct{}{}
		}
	}
	if len(tags) > 0 {
		return std.Keys(tags)
	}
	if guessStringTag(s) {
		tags[FieldTag] = struct{}{}
	}
	if guessNumeric(s) {
		tags[FieldTag] = struct{}{}
	}
	if tag := guessDefault(s); tag != UnknownTag {
		tags[tag] = struct{}{}
	}
	if guessArray(s) {
		tags[ArrayTag] = struct{}{}
	}
	if guessObject(s) {
		tags[ObjectTag] = struct{}{}
	}
	if guessFormat(s) {
		tags[FieldTag] = struct{}{}
	}
	return std.Keys(tags)
}

func getTag(typ schema.SimpleType) TagType {
	switch typ {
	case schema.TypeUnknown:
		return UnknownTag
	case schema.TypeArray:
		return ArrayTag
	case schema.TypeBoolean:
		return FieldTag
	case schema.TypeInteger:
		return FieldTag
	case schema.TypeNull:
		return FieldTag
	case schema.TypeNumber:
		return FieldTag
	case schema.TypeObject:
		return ObjectTag
	case schema.TypeString:
		return FieldTag
	}
	panic("unreachable")
}

func guessNumeric(s *schema.Schema) bool {
	n := s.Numeric
	return n.Maximum != nil || n.Minimum != nil || n.MultipleOf != nil
}

func guessStringTag(s *schema.Schema) bool {
	str := s.String
	return str.MaxLength != nil || str.MinLength > 0 || str.Pattern != nil
}

func guessArray(s *schema.Schema) bool {
	return s.HasArrayConstraints()
}

func guessObject(s *schema.Schema) bool {
	return s.HasObjectConstraints()
}

func guessFormat(s *schema.Schema) bool {
	return s.Format != ""
}

func guessDefault(s *schema.Schema) TagType {
	switch s.Default.(type) {
	case map[string]any:
		return ObjectTag
	case json.Number:
		return FieldTag
	case string:
		return FieldTag
	case bool:
		return FieldTag
	case []any:
		return ArrayTag
	}
	if s.Default == nil {
		return FieldTag
	}
	return UnknownTag
}
