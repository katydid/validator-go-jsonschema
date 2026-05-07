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

func translateRef(refName string) (*ast.Pattern, error) {
	ref := "main"
	if refName != "#" {
		ref = refName
	}
	if strings.HasPrefix(ref, "http") {
		return nil, fmt.Errorf("remoteRef is not supported")
	}
	if strings.HasPrefix(ref, "file:/") {
		return nil, fmt.Errorf("remoteRef file is not supported")
	}
	return ast.NewReference(ref), nil
}
