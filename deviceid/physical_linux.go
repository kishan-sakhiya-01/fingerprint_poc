//go:build linux

package deviceid

import (
	"fmt"
	"os"
	"strings"
)

func physicalIdentity() (identity, error) {
	if isWSL() {
		if id, ok := identityFromCIMProductUUID(); ok {
			return id, nil
		}
	}
	if id, ok := linuxSMBIOSIdentity(); ok {
		return id, nil
	}
	for _, p := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		mid := strings.TrimSpace(string(b))
		if mid != "" {
			return identity{source: SourceLinux, rawID: mid}, nil
		}
	}
	return identity{}, fmt.Errorf("linux: no machine-id found")
}
