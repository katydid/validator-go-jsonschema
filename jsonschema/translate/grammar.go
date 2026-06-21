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
	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go/validator/ast"
)

func NewGrammar(schemaStr []byte, version schema.Version) (*ast.Grammar, error) {
	s, err := schema.ParseSchema(schemaStr)
	if err != nil {
		return nil, err
	}
	s.SetDefaultVersion(version)
	g, err := Translate(s)
	if err != nil {
		return nil, err
	}
	if err := CheckRefs(g); err != nil {
		return nil, err
	}
	return g, err
}
