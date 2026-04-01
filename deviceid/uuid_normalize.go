package deviceid

import (
	"strings"
)

var zeroUUID = strings.ToLower("00000000-0000-0000-0000-000000000000")

// normalizeHardwareUUID returns a lowercase canonical form or false if unusable.
func normalizeHardwareUUID(s string) (string, bool) {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "{}")
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	s = strings.ToLower(s)
	if s == zeroUUID {
		return "", false
	}
	return s, true
}
