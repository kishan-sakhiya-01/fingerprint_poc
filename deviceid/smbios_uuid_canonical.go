package deviceid

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// stabilizeSMBIOSUUID normalizes SMBIOS product / system UUID strings so the same
// underlying firmware UUID matches across Windows WMI and Linux sysfs even when
// the OSes use different endianness for the first three RFC fields (common on PCs).
func stabilizeSMBIOSUUID(s string) (string, bool) {
	u, ok := normalizeHardwareUUID(s)
	if !ok {
		return "", false
	}
	b, err := parseRFC4122UUID(u)
	if err != nil {
		return u, true
	}
	sw := smbiosEndianSwapRFCBytes(b)
	a := formatRFC4122UUID(b)
	c := formatRFC4122UUID(sw)
	if strings.Compare(a, c) <= 0 {
		return a, true
	}
	return c, true
}

// smbiosEndianSwapRFCBytes converts between the two common interpretations of the
// SMBIOS 16-byte UUID (first DWORD and two WORDs are stored little-endian in the table).
func smbiosEndianSwapRFCBytes(b [16]byte) [16]byte {
	var o [16]byte
	o[0], o[1], o[2], o[3] = b[3], b[2], b[1], b[0]
	o[4], o[5] = b[5], b[4]
	o[6], o[7] = b[7], b[6]
	copy(o[8:], b[8:])
	return o
}

func parseRFC4122UUID(s string) ([16]byte, error) {
	var z [16]byte
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return z, fmt.Errorf("uuid: invalid layout")
	}
	h := strings.ReplaceAll(s, "-", "")
	if len(h) != 32 {
		return z, fmt.Errorf("uuid: invalid hex length")
	}
	raw, err := hex.DecodeString(h)
	if err != nil {
		return z, err
	}
	copy(z[:], raw)
	return z, nil
}

func formatRFC4122UUID(b [16]byte) string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5], b[6], b[7],
		b[8], b[9], b[10], b[11],
		b[12], b[13], b[14], b[15])
}
