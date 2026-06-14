// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// We copied these tests from src/net/netip/netip_test.go in the Go standard libary.

package netip

import (
	"testing"
)

func TestValidIPv6(t *testing.T) {
	var validIPs = []struct {
		in    string
		valid bool
	}{
		// IPv4 address in windows-style "print all the digits" form.
		{
			in:    "010.000.015.001",
			valid: false,
		},
		// IPv4 address with a silly amount of leading zeros.
		{
			in:    "000001.00000002.00000003.000000004",
			valid: false,
		},
		// 4-in-6 with octet with leading zero
		{
			in:    "::ffff:1.2.03.4",
			valid: false,
		},
		// 4-in-6 with octet with unexpected character
		{
			in:    "::ffff:1.2.3.z",
			valid: false,
		},
		// Localhost IPv6.
		{
			in:    "::1",
			valid: true,
		},
		// Fully expanded IPv6 address.
		{
			in:    "fd7a:115c:a1e0:ab12:4843:cd96:626b:430b",
			valid: true,
		},
		// IPv6 with elided fields in the middle.
		{
			in:    "fd7a:115c::626b:430b",
			valid: true,
		},
		// IPv6 with elided fields at the end.
		{
			in:    "fd7a:115c:a1e0:ab12:4843:cd96::",
			valid: true,
		},
		// IPv6 with single elided field at the end.
		{
			in:    "fd7a:115c:a1e0:ab12:4843:cd96:626b::",
			valid: true,
		},
		// IPv6 with single elided field in the middle.
		{
			in:    "fd7a:115c:a1e0::4843:cd96:626b:430b",
			valid: true,
		},
		// IPv6 with the trailing 32 bits written as IPv4 dotted decimal. (4in6)
		{
			in:    "::ffff:192.168.140.255",
			valid: true,
		},
		// IPv6 with capital letters.
		{
			in:    "FD9E:1A04:F01D::1",
			valid: true,
		},
	}

	for _, test := range validIPs {
		t.Run(test.in, func(t *testing.T) {
			got := validIPv6([]byte(test.in))
			if got != test.valid {
				t.Fatalf("want %v got %v for %s", test.valid, got, test.in)
			}
		})
	}

	var invalidIPs = []string{
		// Empty string
		"",
		// Garbage non-IP
		"bad",
		// Single number. Some parsers accept this as an IPv4 address in
		// big-endian uint32 form, but we don't.
		"1234",
		// IPv4 with a zone specifier
		"1.2.3.4%eth0",
		// IPv4 field must have at least one digit
		".1.2.3",
		"1.2.3.",
		"1..2.3",
		// IPv4 address too long
		"1.2.3.4.5",
		// IPv4 in dotted octal form
		"0300.0250.0214.0377",
		// IPv4 in dotted hex form
		"0xc0.0xa8.0x8c.0xff",
		// IPv4 in class B form
		"192.168.12345",
		// IPv4 in class B form, with a small enough number to be
		// parseable as a regular dotted decimal field.
		"127.0.1",
		// IPv4 in class A form
		"192.1234567",
		// IPv4 in class A form, with a small enough number to be
		// parseable as a regular dotted decimal field.
		"127.1",
		// IPv4 field has value >255
		"192.168.300.1",
		// IPv4 with too many fields
		"192.168.0.1.5.6",
		// IPv6 with not enough fields
		"1:2:3:4:5:6:7",
		// IPv6 with too many fields
		"1:2:3:4:5:6:7:8:9",
		// IPv6 with 8 fields and a :: expander
		"1:2:3:4::5:6:7:8",
		// IPv6 with a field bigger than 2b
		"fe801::1",
		// IPv6 with non-hex values in field
		"fe80:tail:scal:e::",
		// IPv6 with a zone delimiter but no zone.
		"fe80::1%",
		// IPv6 (without ellipsis) with too many fields for trailing embedded IPv4.
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:192.168.140.255",
		// IPv6 (with ellipsis) with too many fields for trailing embedded IPv4.
		"ffff::ffff:ffff:ffff:ffff:ffff:ffff:192.168.140.255",
		// IPv6 with invalid embedded IPv4.
		"::ffff:192.168.140.bad",
		// IPv6 with multiple ellipsis ::.
		"fe80::1::1",
		// IPv6 with invalid non hex/colon character.
		"fe80:1?:1",
		// IPv6 with truncated bytes after single colon.
		"fe80:",
		// IPv6 with 5 zeros in last group
		"0:0:0:0:0:ffff:0:00000",
		// IPv6 with 5 zeros in one group and embedded IPv4
		"0:0:0:0:00000:ffff:127.1.2.3",
	}

	for _, s := range invalidIPs {
		t.Run(s, func(t *testing.T) {
			if validIPv6([]byte(s)) {
				t.Fatal()
			}
		})
	}
}
