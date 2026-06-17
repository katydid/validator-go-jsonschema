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

package regexformat

import (
	"github.com/dlclark/regexp2/v2"
)

const regexOption regexp2.RegexOptions = regexp2.ECMAScript | regexp2.Unicode

func Compile(expr string) (func(s string) bool, error) {
	faster := tryFastPath(expr)
	if faster != nil {
		return faster, nil
	}
	c, err := regexp2.Compile(expr, regexOption)
	if err != nil {
		return nil, err
	}
	return func(s string) bool {
		res, err := c.MatchString(s)
		return err == nil && res
	}, nil
}
