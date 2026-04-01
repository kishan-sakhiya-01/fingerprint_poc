//go:build windows

package deviceid

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func physicalIdentity() (Identity, error) {
	if id, ok := identityFromCIMProductUUID(); ok {
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