package deviceid

import "testing"

func TestStabilizeSMBIOSUUID_EndianPair(t *testing.T) {
	// Same firmware UUID as typically reported from Linux sysfs vs Windows WMI (byte order of first fields differs).
	linuxForm := "03000200-0400-0500-0607-08090a0b0c0d"
	windowsForm := "00020003-0004-0005-0607-08090a0b0c0d"

	a, ok := stabilizeSMBIOSUUID(linuxForm)
	if !ok {
		t.Fatal("linuxForm")
	}
	b, ok := stabilizeSMBIOSUUID(windowsForm)
	if !ok {
		t.Fatal("windowsForm")
	}
	if a != b {
		t.Fatalf("want same stable id, got %q and %q", a, b)
	}
}

func TestStabilizeSMBIOSUUID_Idempotent(t *testing.T) {
	s := "aabbccdd-eeff-0011-2233-445566778899" // 8-4-4-4-12 hex
	x, ok := stabilizeSMBIOSUUID(s)
	if !ok {
		t.Fatal()
	}
	y, ok := stabilizeSMBIOSUUID(x)
	if !ok || x != y {
		t.Fatalf("second pass changed value: %q -> %q", x, y)
	}
}
