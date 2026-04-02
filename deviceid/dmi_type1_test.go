package deviceid

import "testing"

func TestParseSMBIOSType1UUID(t *testing.T) {
	// Type 1 header through serial-number index (8 bytes), then 16-byte UUID at offset 8.
	header := []byte{0x01, 0x19, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04}
	u := [16]byte{0x03, 0x00, 0x02, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d}
	raw := append(append([]byte{}, header...), u[:]...)

	got, ok := parseSMBIOSType1UUID(raw)
	if !ok {
		t.Fatal("expected ok")
	}
	if got != u {
		t.Fatalf("uuid mismatch: %v vs %v", got, u)
	}
}

func TestParseSMBIOSType1UUID_AllZeroRejected(t *testing.T) {
	raw := make([]byte, 32)
	raw[0] = 1
	raw[1] = 25 // formatted length includes UUID at 8..23
	_, ok := parseSMBIOSType1UUID(raw)
	if ok {
		t.Fatal("expected reject all-zero uuid")
	}
}
