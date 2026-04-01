//go:build !linux && !darwin && !windows

package deviceid

import (
	"fmt"
	"runtime"
)

func physicalIdentity() (Identity, error) {
	return Identity{}, fmt.Errorf("deviceid: unsupported GOOS %q for physical identity", runtime.GOOS)
}
