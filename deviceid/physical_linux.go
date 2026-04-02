//go:build linux

package deviceid

import (
	"fmt"
	"os"
	"strings"
)

func physicalIdentity() (Identity, error) {
	// WSL: use the Windows host SMBIOS UUID (same path as native Windows) so hash matches host + dual-boot Linux when firmware UUID matches.
	if isWSL() {
		if id, ok := identityFromCIMProductUUID(); ok {
			return id, nil
		}
	}
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
	u, ok := stabilizeSMBIOSUUID(string(b))
	if !ok {
		return Identity{}, false
	}
	return Identity{Source: SourceSMBIOS, RawID: u}, true
}
