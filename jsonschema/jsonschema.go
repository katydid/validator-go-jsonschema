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
	"github.com/katydid/validator-go/validator/auto"
	"github.com/katydid/validator-go/validator/intern"
	"github.com/katydid/validator-go/validator/mem"
)

func MatchBytes(schemaStr []byte, jsonStr []byte) (bool, error) {
	i, err := NewInterpreter(schemaStr)
	if err != nil {
		return false, err
	}
	return i.MatchBytes(jsonStr)
}

func MatchParser(schemaStr []byte, p parse.Parser) (bool, error) {
	i, err := NewInterpreter(schemaStr)
	if err != nil {
		return false, err
	}
	return i.MatchParser(p)
}

type Matcher interface {
	MatchBytes([]byte) (bool, error)
	MatchParser(p parse.Parser) (bool, error)
}

type interpret struct {
	parser json.Parser
	g      *ast.Grammar
}

func NewInterpreter(schemaStr []byte) (Matcher, error) {
	g, err := newGrammar(schemaStr)
	if err != nil {
		return nil, err
	}
	p := json.NewJSONSchemaParser()
	return &interpret{
		parser: p,
		g:      g,
	}, nil
}

func (i *interpret) MatchBytes(jsonStr []byte) (bool, error) {
	i.parser.Init(jsonStr)
	return i.MatchParser(i.parser)
}

func (i *interpret) MatchParser(p parse.Parser) (bool, error) {
	return intern.Interpret(i.g, true, p)
}

type memoize struct {
	parser json.Parser
	mem    *mem.Mem
}

func NewMemoizer(schemaStr []byte) (Matcher, error) {
	g, err := newGrammar(schemaStr)
	if err != nil {
		return nil, err
	}
	m, err := mem.NewRecord(g)
	if err != nil {
		return nil, err
	}
	p := json.NewJSONSchemaParser()
	return &memoize{
		parser: p,
		mem:    m,
	}, nil
}

func (m *memoize) MatchBytes(jsonStr []byte) (bool, error) {
	m.parser.Init(jsonStr)
	return m.MatchParser(m.parser)
}

func (m *memoize) MatchParser(p parse.Parser) (bool, error) {
	return validator.Validate(m.mem, p)
}

type compiled struct {
	parser json.Parser
	auto   *auto.Auto
}

func Compile(schemaStr []byte) (Matcher, error) {
	g, err := newGrammar(schemaStr)
	if err != nil {
		return nil, err
	}
	a, err := auto.CompileRecord(g)
	if err != nil {
		return nil, err
	}
	p := json.NewJSONSchemaParser()
	return &compiled{
		parser: p,
		auto:   a,
	}, nil
}

func (c *compiled) MatchBytes(jsonStr []byte) (bool, error) {
	c.parser.Init(jsonStr)
	return c.MatchParser(c.parser)
}

func (c *compiled) MatchParser(p parse.Parser) (bool, error) {
	return c.auto.Validate(p)
}

func newGrammar(schemaStr []byte) (*ast.Grammar, error) {
	schema, err := schema.ParseSchema(schemaStr)
	if err != nil {
		return nil, err
	}
	g, err := translate.Translate(schema)
	if err != nil {
		return nil, err
	}
	if err := translate.CheckRefs(g); err != nil {
		return nil, err
	}
	return g, err
}
