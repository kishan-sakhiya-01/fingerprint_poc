//go:build linux

package deviceid

import (
	"fmt"
	"os"
	"strings"
)

func physicalIdentity() (Identity, error) {
	if id, ok := linuxSMBIOSUUID(); ok {
		return id, nil
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

func linuxSMBIOSUUID() (Identity, bool) {
	b, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return Identity{}, false
	}
	u, ok := normalizeHardwareUUID(string(b))
	if !ok {
		return Identity{}, false
	}
	return Identity{Source: SourceSMBIOS, RawID: u}, true
}
