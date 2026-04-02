package deviceid

import "testing"

func TestNormalizeBootDisk_serialSpacing(t *testing.T) {
	a, ok := normalizeBootDiskSerial("  ABC 12 3 ")
	if !ok || a != "abc123" {
		t.Fatalf("got %q %v", a, ok)
	}
}
