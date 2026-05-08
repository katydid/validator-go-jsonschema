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
	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/translate"
	"github.com/katydid/validator-go/validator"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/intern"
	"github.com/katydid/validator-go/validator/mem"
)

func Validate(schemaStr []byte, jsonStr []byte) (bool, error) {
	p := json.NewJSONSchemaParser()
	p.Init(jsonStr)
	return ValidateParser(schemaStr, p)
}

func ValidateParser(schemaStr []byte, p parse.Parser) (bool, error) {
	g, err := newGrammar(schemaStr)
	if err != nil {
		return false, err
	}
	if err := translate.CheckRefs(g); err != nil {
		return false, err
	}
	return intern.Interpret(g, true, p)
}

func newGrammar(schemaStr []byte) (*ast.Grammar, error) {
	schema, err := schema.ParseSchema(schemaStr)
	if err != nil {
		return nil, err
	}
	g, err := translate.TranslateDraft4(schema)
	if err != nil {
		return nil, err
	}
	return g, err
}

type Filter interface {
	Validate([]byte) (bool, error)
}

type filter struct {
	parser json.Parser
	mem    *mem.Mem
}

func NewFilter(schemaStr []byte) (Filter, error) {
	schema, err := schema.ParseSchema(schemaStr)
	if err != nil {
		return nil, err
	}
	g, err := translate.TranslateDraft4(schema)
	if err != nil {
		return nil, err
	}
	m, err := mem.NewRecord(g)
	if err != nil {
		return nil, err
	}
	p := json.NewJSONSchemaParser()
	return &filter{
		parser: p,
		mem:    m,
	}, nil
}

func (f *filter) Validate(jsonStr []byte) (bool, error) {
	f.parser.Init(jsonStr)
	return validator.Validate(f.mem, f.parser)
}
