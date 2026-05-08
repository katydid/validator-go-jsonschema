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
	"net/url"
	"strings"

	"github.com/qri-io/jsonpointer"
)

const reservedWordForEmpty = "reserved word for empty definition path"

func parsePointer(s string) ([]string, error) {
	// sometimes we forget to strip the hash from the front.
	if strings.HasPrefix(s, "#") {
		s = s[1:]
	}
	path, err := jsonpointer.Parse(s)
	if err != nil {
		return nil, err
	}
	// This decodes the percent encoding, changing %25 to %
	for i, p := range path {
		u, err := url.PathUnescape(p)
		if err == nil {
			// We ignore errors for paths that are already escaped.
			path[i] = u
		}
	}
	for i, p := range path {
		if p == "" {
			path[i] = reservedWordForEmpty
		}
	}
	return path, nil
}
