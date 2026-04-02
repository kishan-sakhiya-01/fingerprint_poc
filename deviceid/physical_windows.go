//go:build windows

package deviceid

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func physicalIdentity() (identity, error) {
	if id, ok := identityFromCIMProductUUID(); ok {
		return id, nil
	}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		return identity{}, fmt.Errorf("windows registry: %w", err)
	}
	defer k.Close()
	guid, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return identity{}, fmt.Errorf("windows MachineGuid: %w", err)
	}
	guid = strings.TrimSpace(guid)
	if guid == "" {
		return identity{}, fmt.Errorf("windows: empty MachineGuid")
	}
	return identity{source: SourceWindows, rawID: guid}, nil
}