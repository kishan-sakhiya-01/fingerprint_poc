//go:build linux

package deviceid

import (
	"fmt"
	"os"
	"strings"
)

func physicalIdentity() (Identity, error) {
	// WSL: host SMBIOS first (same order as Windows), then host boot disk; avoid guest sysfs disk for host parity.
	if isWSL() {
		if id, ok := identityFromCIMProductUUID(); ok {
			return id, nil
		}
		if id, ok := identityBootDiskFromPS(); ok {
			return id, nil
		}
		if id, ok := linuxSMBIOSIdentity(); ok {
			return id, nil
		}
	} else {
		if id, ok := linuxSMBIOSIdentity(); ok {
			return id, nil
		}
		if id, ok := tryLinuxBootDiskSerialSysfs(); ok {
			return id, nil
		}
	}
	for _, p := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		id := strings.TrimSpace(string(b))
		if id != "" {
			return Identity{Source: SourceLinux, RawID: id}, nil
		}
	}
	return Identity{}, fmt.Errorf("linux: no machine-id found")
}
