package jsonschema

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// see https://www.rfc-editor.org/rfc/rfc6901#section-3
// json schema validator for format: json-pointer
func ValidateJSONPointer(s string) error {
	if s == "" {
		return nil
	}
	if !strings.HasPrefix(s, "/") {
		return fmt.Errorf("not starting with /")
	}
	for _, tok := range strings.Split(s, "/")[1:] {
		escape := false
		for _, ch := range tok {
			if escape {
				escape = false
				if ch != '0' && ch != '1' {
					return fmt.Errorf("~ must be followed by 0 or 1")
				}
				continue
			}
			if ch == '~' {
				escape = true
				continue
			}
			switch {
			case ch >= '\x00' && ch <= '\x2E':
			case ch >= '\x30' && ch <= '\x7D':
			case ch >= '\x7F' && ch <= '\U0010FFFF':
			default:
				return fmt.Errorf("invalid character %q", ch)
			}
		}
		if escape {
			return fmt.Errorf("~ must be followed by 0 or 1")
		}
	}
	return nil
}

// see https://tools.ietf.org/html/draft-handrews-relative-json-pointer-01#section-3
// json schema validator for format: relative-json-pointer
func ValidateRelativeJSONPointer(s string) error {
	// start with non-negative-integer
	numDigits := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			numDigits++
		} else {
			break
		}
	}
	if numDigits == 0 {
		return fmt.Errorf("must start with non-negative integer")
	}
	if numDigits > 1 && strings.HasPrefix(s, "0") {
		return fmt.Errorf("starts with zero")
	}
	s = s[numDigits:]

	// followed by either json-pointer or '#'
	if s == "#" {
		return nil
	}
	return ValidateJSONPointer(s)
}

// see https://datatracker.ietf.org/doc/html/rfc4122#page-4
// json schema validator for format: uuid
func ValidateUUID(s string) error {
	hexGroups := []int{8, 4, 4, 4, 12}
	groups := strings.Split(s, "-")
	if len(groups) != len(hexGroups) {
		return fmt.Errorf("must have %d elements", len(hexGroups))
	}
	for i, group := range groups {
		if len(group) != hexGroups[i] {
			return fmt.Errorf("element %d must be %d characters long", i+1, hexGroups[i])
		}
		for _, ch := range group {
			switch {
			case ch >= '0' && ch <= '9':
			case ch >= 'a' && ch <= 'f':
			case ch >= 'A' && ch <= 'F':
			default:
				return fmt.Errorf("non-hex character %q", ch)
			}
		}
	}
	return nil
}

// see https://datatracker.ietf.org/doc/html/rfc3339#appendix-A
// json schema validator for format: duration
func ValidateDuration(s string) error {
	// must start with 'P'
	var ok bool
	s, ok = strings.CutPrefix(s, "P")
	if !ok {
		return fmt.Errorf("must start with P")
	}
	if s == "" {
		return fmt.Errorf("nothing after P")
	}

	// dur-week
	if s, ok := strings.CutSuffix(s, "W"); ok {
		if s == "" {
			return fmt.Errorf("no number in week")
		}
		for _, ch := range s {
			if ch < '0' || ch > '9' {
				return fmt.Errorf("invalid week")
			}
		}
		return nil
	}

	allUnits := []string{"YMD", "HMS"}
	for i, s := range strings.Split(s, "T") {
		if i != 0 && s == "" {
			return fmt.Errorf("no time elements")
		}
		if i >= len(allUnits) {
			return fmt.Errorf("more than one T")
		}
		units := allUnits[i]
		for s != "" {
			digitCount := 0
			for _, ch := range s {
				if ch >= '0' && ch <= '9' {
					digitCount++
				} else {
					break
				}
			}
			if digitCount == 0 {
				return fmt.Errorf("missing number")
			}
			s = s[digitCount:]
			if s == "" {
				return fmt.Errorf("missing unit")
			}
			unit := s[0]
			j := strings.IndexByte(units, unit)
			if j == -1 {
				if strings.IndexByte(allUnits[i], unit) != -1 {
					return fmt.Errorf("unit %q out of order", unit)
				}
				return fmt.Errorf("invalid unit %q", unit)
			}
			units = units[j+1:]
			s = s[1:]
		}
	}

	return nil
}

// see see https://datatracker.ietf.org/doc/html/rfc3339#section-5.6
// json schema validator for format: date
func validateDate(s string) error {
	_, err := time.Parse("2006-01-02", s)
	return err
}

// see https://datatracker.ietf.org/doc/html/rfc3339#section-5.6
// NOTE: golang time package does not support leap seconds.
// json schema validator for format: time
func ValidateTime(str string) error {
	// min: hh:mm:ssZ
	if len(str) < 9 {
		return fmt.Errorf("less than 9 characters long")
	}
	if str[2] != ':' || str[5] != ':' {
		return fmt.Errorf("missing colon in correct place")
	}

	// parse hh:mm:ss
	var hms []int
	for _, tok := range strings.SplitN(str[:8], ":", 3) {
		i, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("invalid hour/min/sec")
		}
		if i < 0 {
			return fmt.Errorf("non-positive hour/min/sec")
		}
		hms = append(hms, i)
	}
	if len(hms) != 3 {
		return fmt.Errorf("missing hour/min/sec")
	}
	h, m, s := hms[0], hms[1], hms[2]
	if h > 23 || m > 59 || s > 60 {
		return fmt.Errorf("hour/min/sec out of range")
	}
	str = str[8:]

	// parse sec-frac if present
	if rem, ok := strings.CutPrefix(str, "."); ok {
		numDigits := 0
		for _, ch := range rem {
			if ch >= '0' && ch <= '9' {
				numDigits++
			} else {
				break
			}
		}
		if numDigits == 0 {
			return fmt.Errorf("no digits in second fraction")
		}
		str = rem[numDigits:]
	}

	if str != "z" && str != "Z" {
		// parse time-numoffset
		if len(str) != 6 {
			return fmt.Errorf("offset must be 6 characters long")
		}
		var sign int
		switch str[0] {
		case '+':
			sign = -1
		case '-':
			sign = +1
		default:
			return fmt.Errorf("offset must begin with plus/minus")
		}
		str = str[1:]
		if str[2] != ':' {
			return fmt.Errorf("missing colon in offset in correct place")
		}

		var zhm []int
		for _, tok := range strings.SplitN(str, ":", 2) {
			i, err := strconv.Atoi(tok)
			if err != nil {
				return fmt.Errorf("invalid hour/min in offset")
			}
			if i < 0 {
				return fmt.Errorf("non-positive hour/min in offset")
			}
			zhm = append(zhm, i)
		}
		zh, zm := zhm[0], zhm[1]
		if zh > 23 || zm > 59 {
			return fmt.Errorf("hour/min in offset out of range")
		}

		// apply timezone
		hm := (h*60 + m) + sign*(zh*60+zm)
		if hm < 0 {
			hm += 24 * 60
		}
		h, m = hm/60, hm%60
	}

	// check leap second
	if s >= 60 && (h != 23 || m != 59) {
		return fmt.Errorf("invalid leap second")
	}

	return nil
}

// see https://datatracker.ietf.org/doc/html/rfc3339#section-5.6
// json schema validator for format: date-time
func validateDateTime(s string) error {
	// min: yyyy-mm-ddThh:mm:ssZ
	if len(s) < 20 {
		return fmt.Errorf("less than 20 characters long")
	}

	if s[10] != 't' && s[10] != 'T' {
		return fmt.Errorf("11th character must be t or T")
	}
	if err := validateDate(s[:10]); err != nil {
		return fmt.Errorf("invalid date element: %v", err)
	}
	if err := ValidateTime(s[11:]); err != nil {
		return fmt.Errorf("invalid time element: %v", err)
	}
	return nil
}

// json schema validator for format: period
func ValidatePeriod(s string) error {
	slash := strings.IndexByte(s, '/')
	if slash == -1 {
		return fmt.Errorf("missing slash")
	}

	start, end := s[:slash], s[slash+1:]
	if strings.HasPrefix(start, "P") {
		if err := ValidateDuration(start); err != nil {
			return fmt.Errorf("invalid start duration: %v", err)
		}
		if err := validateDateTime(end); err != nil {
			return fmt.Errorf("invalid end date-time: %v", err)
		}
	} else {
		if err := validateDateTime(start); err != nil {
			return fmt.Errorf("invalid start date-time: %v", err)
		}
		if strings.HasPrefix(end, "P") {
			if err := ValidateDuration(end); err != nil {
				return fmt.Errorf("invalid end duration: %v", err)
			}
		} else if err := validateDateTime(end); err != nil {
			return fmt.Errorf("invalid end date-time: %v", err)
		}
	}

	return nil
}

// see https://semver.org/#backusnaur-form-grammar-for-valid-semver-versions
// json schema validator for format: semver
func ValidateSemver(s string) error {
	// build --
	if i := strings.IndexByte(s, '+'); i != -1 {
		build := s[i+1:]
		if build == "" {
			return fmt.Errorf("build is empty")
		}
		for _, buildID := range strings.Split(build, ".") {
			if buildID == "" {
				return fmt.Errorf("build identifier is empty")
			}
			for _, ch := range buildID {
				switch {
				case ch >= '0' && ch <= '9':
				case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '-':
				default:
					return fmt.Errorf("invalid character %q in build identifier", ch)
				}
			}
		}
		s = s[:i]
	}

	// pre-release --
	if i := strings.IndexByte(s, '-'); i != -1 {
		preRelease := s[i+1:]
		for _, preReleaseID := range strings.Split(preRelease, ".") {
			if preReleaseID == "" {
				return fmt.Errorf("pre-release identifier is empty")
			}
			allDigits := true
			for _, ch := range preReleaseID {
				switch {
				case ch >= '0' && ch <= '9':
				case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '-':
					allDigits = false
				default:
					return fmt.Errorf("invalid character %q in pre-release identifier", ch)
				}
			}
			if allDigits && len(preReleaseID) > 1 && preReleaseID[0] == '0' {
				return fmt.Errorf("pre-release numeric identifier starts with zero")
			}
		}
		s = s[:i]
	}

	// versionCore --
	versions := strings.Split(s, ".")
	if len(versions) != 3 {
		return fmt.Errorf("versionCore must have 3 numbers separated by dot")
	}
	names := []string{"major", "minor", "patch"}
	for i, version := range versions {
		if version == "" {
			return fmt.Errorf("%s is empty", names[i])
		}
		if len(version) > 1 && version[0] == '0' {
			return fmt.Errorf("%s starts with zero", names[i])
		}
		for _, ch := range version {
			if ch < '0' || ch > '9' {
				return fmt.Errorf("%s contains non-digit", names[i])
			}
		}
	}

	return nil
}
