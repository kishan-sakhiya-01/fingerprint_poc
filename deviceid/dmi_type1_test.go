package deviceid

import "testing"

func TestParseSMBIOSType1UUID(t *testing.T) {
	// Type 1: formatted length 24 ends right after 16-byte UUID (offset 8..23).
	header := []byte{0x01, 24, 0, 0, 0x01, 0x02, 0x03, 0x04}
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

func TestWalkSMBIOSStructureTable(t *testing.T) {
	u := [16]byte{0x03, 0x00, 0x02, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d}
	header := []byte{0x01, 24, 0, 0, 0x01, 0x02, 0x03, 0x04}
	s1 := append(append([]byte{}, header...), u[:]...)
	s1 = append(s1, 0, 0)
	end := []byte{127, 4, 0, 0, 0, 0}
	table := append(s1, end...)
	got, ok := walkSMBIOSStructureTable(table)
	if !ok || got != u {
		t.Fatalf("got %v ok=%v want %v", got, ok, u)
	}
}

func TestParseSMBIOSType1UUID_AllZeroRejected(t *testing.T) {
	raw := make([]byte, 32)
	raw[0] = 1
	raw[1] = 24
	_, ok := parseSMBIOSType1UUID(raw)
	if ok {
		t.Fatal("expected reject all-zero uuid")
	}
}
