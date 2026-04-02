package deviceid

import "fmt"

// smbiosUUIDBytesUsable reports whether the 16-byte DMI UUID field is set (not all-zero / all-0xFF).
func smbiosUUIDBytesUsable(b [16]byte) bool {
	all0, allFF := true, true
	for _, x := range b {
		if x != 0 {
			all0 = false
		}
		if x != 0xff {
			allFF = false
		}
	}
	return !all0 && !allFF
}

// parseSMBIOSType1UUID extracts the 16-byte system UUID from a Type-1 (System Information) SMBIOS structure raw blob.
func parseSMBIOSType1UUID(raw []byte) ([16]byte, bool) {
	var z [16]byte
	if len(raw) < 24 {
		return z, false
	}
	if raw[0] != 1 {
		return z, false
	}
	L := int(raw[1])
	if L < 24 {
		return z, false
	}
	copy(z[:], raw[8:24])
	if !smbiosUUIDBytesUsable(z) {
		return z, false
	}
	return z, true
}

// smbiosTableMemoryToUUIDString formats UUID bytes as Linux sysfs typically does (little-endian first three fields in the string).
func smbiosTableMemoryToUUIDString(b [16]byte) string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[3], b[2], b[1], b[0],
		b[5], b[4],
		b[7], b[6],
		b[8], b[9], b[10], b[11],
		b[12], b[13], b[14], b[15])
}
