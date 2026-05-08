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

func translateRef(name string) (*ast.Pattern, error) {
	if name == "#" {
		return ast.NewReference("main"), nil
	}
	if strings.HasPrefix(name, "#/") {
		refName, err := newRefName(name)
		if err != nil {
			return nil, err
		}
		return ast.NewReference(refName), nil
	}
	if strings.HasPrefix(name, "http") {
		return nil, fmt.Errorf("remoteRef is not supported")
	}
	if strings.HasPrefix(name, "file:/") {
		return nil, fmt.Errorf("remoteRef file is not supported")
	}
	return nil, fmt.Errorf("unsupported reference type %s", name)
}

func newRefName(s string) (string, error) {
	path, err := parsePointer(s)
	if err != nil {
		return "", err
	}
	return strings.Join(path, "/"), nil
}
