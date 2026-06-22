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
	"strings"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go/validator/ast"
)

func translateRef(id string, name string) (*ast.Pattern, error) {
	defName, err := refToDefName(id, name)
	if err != nil {
		return nil, err
	}
	return ast.NewReference(defName), nil
}

func refToDefName(id string, ref string) (string, error) {
	if strings.HasPrefix(ref, "file:/") {
		return "", fmt.Errorf("remoteRef file is not supported")
	}
	if ref == "#" {
		return "main", nil
	}
	if strings.HasPrefix(ref, "#") && !strings.HasPrefix(ref, "#/") {
		// anchor
		return ref, nil
	}
	paths, err := parsePointer(ref)
	if err != nil {
		return "", err
	}
	path := strings.Join(paths, "/")
	if strings.HasPrefix(ref, "#/") {
		return id + path, nil
	}
	return path, nil
}

func definitionToPrefix(prefix string, name string, sch *schema.Schema) string {
	return "/definitions/" + name
}

func definitionToDefName(prefix string, name string, sch *schema.Schema) (string, error) {
	if len(sch.Id) > 0 {
		return prefix + sch.Id, nil
	}
	if len(sch.Anchor) > 0 {
		return "#" + sch.Anchor, nil
	}
	name = "/definitions/" + name
	s := prefix + name
	path, err := parsePointer(s)
	if err != nil {
		return "", err
	}
	return strings.Join(path, "/"), nil
}
