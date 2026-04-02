//go:build linux

package deviceid

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var linuxProductUUIDPaths = []string{
	"/sys/class/dmi/id/product_uuid",
	"/sys/devices/virtual/dmi/id/product_uuid",
}

func tryLinuxSysfsProductUUID() (Identity, bool) {
	for _, p := range linuxProductUUIDPaths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if id, ok := identityFromLinuxProductUUIDBytes(b); ok {
			return id, true
		}
	}
	return Identity{}, false
}

func identityFromLinuxProductUUIDBytes(b []byte) (Identity, bool) {
	s := strings.TrimSpace(string(b))
	if s == "" {
		return Identity{}, false
	}
	lo := strings.ToLower(s)
	if strings.Contains(lo, "not specified") || strings.Contains(lo, "not settable") || strings.Contains(lo, "unknown") {
		return Identity{}, false
	}
	u, ok := stabilizeSMBIOSUUID(s)
	if !ok {
		return Identity{}, false
	}
	return Identity{Source: SourceSMBIOS, RawID: u}, true
}

func tryLinuxFirmwareDMIType1() (Identity, bool) {
	const base = "/sys/firmware/dmi/entries"
	ents, err := os.ReadDir(base)
	if err != nil {
		return Identity{}, false
	}
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(base, e.Name())
		tb, err := os.ReadFile(filepath.Join(dir, "type"))
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(tb)) != "1" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, "raw"))
		if err != nil {
			continue
		}
		u16, ok := parseSMBIOSType1UUID(raw)
		if !ok {
			continue
		}
		s := smbiosTableMemoryToUUIDString(u16)
		u, ok := stabilizeSMBIOSUUID(s)
		if !ok {
			continue
		}
		return Identity{Source: SourceSMBIOS, RawID: u}, true
	}
	return Identity{}, false
}

func tryLinuxDmidecodeSystemUUID() (Identity, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "dmidecode", "-s", "system-uuid")
	out, err := cmd.Output()
	if err != nil {
		return Identity{}, false
	}
	return identityFromLinuxProductUUIDBytes(out)
}

func linuxSMBIOSIdentity() (Identity, bool) {
	if id, ok := tryLinuxSysfsProductUUID(); ok {
		return id, true
	}
	if id, ok := tryLinuxFirmwareDMIType1(); ok {
		return id, true
	}
	if id, ok := tryLinuxDmidecodeSystemUUID(); ok {
		return id, true
	}
	return Identity{}, false
}
