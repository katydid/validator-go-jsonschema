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
	"fmt"

	"github.com/katydid/validator-go/validator/ast"
)

type visitor struct {
	refs []string
}

func (v *visitor) Visit(node interface{}) interface{} {
	p, ok := node.(*ast.Pattern)
	if !ok {
		return v
	}
	if p.Reference != nil {
		v.refs = append(v.refs, p.Reference.GetName())
	}
	return v
}

func CheckRefs(g *ast.Grammar) error {
	v := &visitor{refs: []string{}}
	refLookup := ast.NewRefLookup(g)
	for _, p := range refLookup {
		p.Walk(v)
	}
	for _, refName := range v.refs {
		if _, ok := refLookup[refName]; !ok {
			return fmt.Errorf("reference to unknown definition %s", refName)
		}
	}
	return nil
}
