//go:build linux

package deviceid

import (
	"fmt"
	"os"
	"strings"
)

func physicalIdentity() (Identity, error) {
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
