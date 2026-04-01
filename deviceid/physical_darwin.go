//go:build darwin

package deviceid

import (
	"fmt"
	"os/exec"
	"strings"
)

func physicalIdentity() (Identity, error) {
	out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return Identity{}, fmt.Errorf("darwin ioreg: %w", err)
	}
	const key = `"IOPlatformUUID" = "`
	s := string(out)
	i := strings.Index(s, key)
	if i < 0 {
		return Identity{}, fmt.Errorf("darwin: IOPlatformUUID not found")
	}
	s = s[i+len(key):]
	j := strings.IndexByte(s, '"')
	if j < 0 {
		return Identity{}, fmt.Errorf("darwin: malformed IOPlatformUUID")
	}
	uuid := strings.TrimSpace(s[:j])
	if uuid == "" {
		return Identity{}, fmt.Errorf("darwin: empty IOPlatformUUID")
	}
	return Identity{Source: SourceDarwin, RawID: uuid}, nil
}
