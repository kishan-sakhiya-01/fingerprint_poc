//go:build darwin

package deviceid

import (
	"fmt"
	"os/exec"
	"strings"
)

func physicalIdentity() (identity, error) {
	out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return identity{}, fmt.Errorf("darwin ioreg: %w", err)
	}
	const key = `"IOPlatformUUID" = "`
	s := string(out)
	i := strings.Index(s, key)
	if i < 0 {
		return identity{}, fmt.Errorf("darwin: IOPlatformUUID not found")
	}
	s = s[i+len(key):]
	j := strings.IndexByte(s, '"')
	if j < 0 {
		return identity{}, fmt.Errorf("darwin: malformed IOPlatformUUID")
	}
	uuid := strings.TrimSpace(s[:j])
	if uuid == "" {
		return identity{}, fmt.Errorf("darwin: empty IOPlatformUUID")
	}
	return identity{source: SourceDarwin, rawID: uuid}, nil
}
