//go:build windows

package deviceid

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func physicalIdentity() (Identity, error) {
	if id, ok := windowsSMBIOSUUID(); ok {
		return id, nil
	}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		return Identity{}, fmt.Errorf("windows registry: %w", err)
	}
	defer k.Close()
	guid, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return Identity{}, fmt.Errorf("windows MachineGuid: %w", err)
	}
	guid = strings.TrimSpace(guid)
	if guid == "" {
		return Identity{}, fmt.Errorf("windows: empty MachineGuid")
	}
	return Identity{Source: SourceWindows, RawID: guid}, nil
}

func windowsSMBIOSUUID() (Identity, bool) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "(Get-CimInstance -ClassName Win32_ComputerSystemProduct).UUID")
	out, err := cmd.Output()
	if err != nil {
		return Identity{}, false
	}
	u, ok := normalizeHardwareUUID(string(out))
	if !ok {
		return Identity{}, false
	}
	return Identity{Source: SourceSMBIOS, RawID: u}, true
}