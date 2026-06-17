// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"fmt"

	"github.com/jtacoma/uritemplates"
)

// uriTemplateFormat requires a valid URI template.
func uriTemplateFormat(s string) error {
	if _, err := uritemplates.Parse(s); err != nil {
		return fmt.Errorf("%q is not a valid URI template: %v", s, err)
	}
	return nil
}
