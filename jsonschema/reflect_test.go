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

package jsonschema

import (
	goreflect "reflect"

	"github.com/katydid/parser-go-reflect/reflect"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
)

func unmarshal(data []byte) (any, error) {
	var m any
	err := std.UnmarshalJSON(data, &m)
	return m, err
}

func newReflectParser() parse.ParserWithInit {
	p := reflect.NewJSONSchemaParser()
	return &reflectWithInit{p}
}

type reflectWithInit struct {
	reflect.Parser
}

func (r *reflectWithInit) Init(data []byte) {
	m, err := unmarshal(data)
	if err != nil {
		panic(err)
	}
	r.Parser.Init(goreflect.ValueOf(m))
}
