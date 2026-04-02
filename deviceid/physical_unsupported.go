//go:build !linux && !darwin && !windows

package deviceid

import (
	"fmt"
	"runtime"
)

func physicalIdentity() (identity, error) {
	return identity{}, fmt.Errorf("deviceid: unsupported GOOS %q for physical identity", runtime.GOOS)
}
