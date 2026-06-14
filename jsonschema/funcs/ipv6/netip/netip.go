// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// We copy src/net/netip/netip.go from Go standard library and modify it to work on bytes, not allocate and only return a bool.

package netip

import (
	"bytes"
	"math"
)

// parseIPv6 parses s as an IPv6 address (in form "2001:db8::68").
func validIPv6(in []byte) bool {
	s := in

	// Split off the zone right from the start. Yes it's a second scan
	// of the string, but trying to handle it inline makes a bunch of
	// other inner loop conditionals more expensive, and it ends up
	// being slower.
	var zone []byte
	i := bytes.IndexByte(s, '%')
	if i != -1 {
		s, zone = s[:i], s[i+1:]
		if len(zone) == 0 {
			// Not allowed to have an empty zone if explicitly specified.
			return false
		}
		if len(zone) > 0 {
			// Not allowed for JSON Schema IPv6
			return false
		}
	}

	var ip [16]byte
	ellipsis := -1 // position of ellipsis in ip

	// Might have leading ellipsis
	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]
		// Might be only ellipsis
		if len(s) == 0 {
			return true
		}
	}

	// Loop, parsing hex numbers followed by colon.
	i = 0
	for i < 16 {
		// Hex number. Similar to parseIPv4, inlining the hex number
		// parsing yields a significant performance increase.
		off := 0
		acc := uint32(0)
		for ; off < len(s); off++ {
			c := s[off]
			if c >= '0' && c <= '9' {
				acc = (acc << 4) + uint32(c-'0')
			} else if c >= 'a' && c <= 'f' {
				acc = (acc << 4) + uint32(c-'a'+10)
			} else if c >= 'A' && c <= 'F' {
				acc = (acc << 4) + uint32(c-'A'+10)
			} else {
				break
			}
			if off > 3 {
				//more than 4 digits in group, fail.
				return false
			}
			if acc > math.MaxUint16 {
				// Overflow, fail.
				return false
			}
		}
		if off == 0 {
			// No digits found, fail.
			return false
		}

		// If followed by dot, might be in trailing IPv4.
		if off < len(s) && s[off] == '.' {
			if ellipsis < 0 && i != 12 {
				// Not the right place.
				return false
			}
			if i+4 > 16 {
				// Not enough room.
				return false
			}

			end := len(in)
			if len(zone) > 0 {
				end -= len(zone) + 1
			}

			if !parseIPv4Fields(in, end-len(s), end, ip[i:i+4]) {
				return false
			}
			s = []byte{}
			i += 4
			break
		}

		// Save this 16-bit chunk.
		ip[i] = byte(acc >> 8)
		ip[i+1] = byte(acc)
		i += 2

		// Stop at end of string.
		s = s[off:]
		if len(s) == 0 {
			break
		}

		// Otherwise must be followed by colon and more.
		if s[0] != ':' {
			return false
		} else if len(s) == 1 {
			return false
		}
		s = s[1:]

		// Look for ellipsis.
		if s[0] == ':' {
			if ellipsis >= 0 { // already have one
				return false
			}
			ellipsis = i
			s = s[1:]
			if len(s) == 0 { // can be at end
				break
			}
		}
	}

	// Must have used entire string.
	if len(s) != 0 {
		return false
	}

	// If didn't parse enough, expand ellipsis.
	if i < 16 {
		if ellipsis < 0 {
			return false
		}
		n := 16 - i
		for j := i - 1; j >= ellipsis; j-- {
			ip[j+n] = ip[j]
		}
		clear(ip[ellipsis : ellipsis+n])
	} else if ellipsis >= 0 {
		// Ellipsis must represent at least one 0 group.
		return false
	}
	return true
}

func parseIPv4Fields(in []byte, off, end int, fields []uint8) bool {
	var val, pos int
	var digLen int // number of digits in current octet
	s := in[off:end]
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			if digLen == 1 && val == 0 {
				return false
			}
			val = val*10 + int(s[i]) - '0'
			digLen++
			if val > 255 {
				return false
			}
		} else if s[i] == '.' {
			// .1.2.3
			// 1.2.3.
			// 1..2.3
			if i == 0 || i == len(s)-1 || s[i-1] == '.' {
				return false
			}
			// 1.2.3.4.5
			if pos == 3 {
				return false
			}
			fields[pos] = uint8(val)
			pos++
			val = 0
			digLen = 0
		} else {
			return false
		}
	}
	if pos < 3 {
		return false
	}
	fields[3] = uint8(val)
	return true
}
