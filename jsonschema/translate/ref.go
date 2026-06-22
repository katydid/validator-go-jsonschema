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
	if ref == "#" {
		return "main", nil
	}
	if strings.HasPrefix(ref, "file:/") {
		return "", fmt.Errorf("remoteRef file is not supported")
	}
	if strings.HasPrefix(ref, "#/") {
		path, err := parsePointer(ref)
		if err != nil {
			return "", err
		}
		refName := id + strings.Join(path, "/")
		return refName, nil
	}
	if strings.HasPrefix(ref, "#") {
		return ref, nil
	}
	s := ref
	path, err := parsePointer(s)
	if err != nil {
		return "", err
	}
	return strings.Join(path, "/"), nil
}

func definitionToPrefix(prefix string, name string, sch *schema.Schema) string {
	name = "/definitions/" + name
	if len(sch.Id) > 0 {
		return sch.Id
	}
	if len(sch.Anchor) > 0 {
		return "#" + sch.Anchor
	}
	name = prefix + name
	return name
}

func definitionToDefName(prefix string, name string, sch *schema.Schema) (string, error) {
	name = "/definitions/" + name
	if len(sch.Id) > 0 {
		return sch.Id, nil
	}
	if len(sch.Anchor) > 0 {
		return "#" + sch.Anchor, nil
	}
	s := prefix + name
	if strings.HasPrefix(s, "#") && !strings.HasPrefix(s, "#/") && s != "#" {
		// anchors like #bla are also allowed
		return s, nil
	}
	path, err := parsePointer(s)
	if err != nil {
		return "", err
	}
	return strings.Join(path, "/"), nil
}
