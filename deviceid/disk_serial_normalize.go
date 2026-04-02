package deviceid

import "strings"

// normalizeBootDiskSerial returns a compact comparable form or false if empty after normalization.
func normalizeBootDiskSerial(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch r {
		case ' ', '\t', '\n', '\r', '.', '_', '-', ':', '\\', '/':
			continue
		default:
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "", false
	}
	return out, true
}
