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

	"github.com/katydid/validator-go/validator/ast"
)

func translateRef(parentId string, name string) (*ast.Pattern, error) {
	defName, err := refToDefName(parentId, name)
	if err != nil {
		return nil, err
	}
	return ast.NewReference(defName), nil
}

func refToDefName(parentId string, ref string) (string, error) {
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
	if strings.HasPrefix(ref, "#/") {
		path := strings.Join(paths, "/")
		return path, nil
	}
	return prependParentId(parentId, paths)
}

func definitionToPrefix(prefix string, name string, id string) string {
	return "/definitions/" + name
}

func prependParentId(parentId string, paths []string) (string, error) {
	if parentId == "" {
		return strings.Join(paths, "/"), nil
	}
	parentPaths, err := parsePointer(parentId)
	if err != nil {
		return "", err
	}
	i := 0
	for i < len(parentPaths) && i < len(paths) && parentPaths[i] == paths[i] {
		i++
	}
	if i != 0 {
		paths = append(parentPaths[:i], paths[i:]...)
		return strings.Join(paths, "/"), nil
	}
	// remove last slash or last item
	parentPaths = parentPaths[:len(parentPaths)-1]
	return strings.Join(append(parentPaths, paths...), "/"), nil
}

func definitionToDefName(prefix string, parentId string, name string, id string, anchor string) (string, error) {
	if len(anchor) > 0 {
		return "#" + anchor, nil
	}
	if len(id) > 0 {
		if strings.HasPrefix(id, "#") && !strings.HasPrefix(id, "#/") {
			// anchor
			return id, nil
		}
		paths, err := parsePointer(id)
		if err != nil {
			return "", err
		}
		return prependParentId(parentId, paths)
	}
	name = "/definitions/" + name
	s := prefix + name
	paths, err := parsePointer(s)
	if err != nil {
		return "", err
	}
	return prependParentId(parentId, paths)
}
