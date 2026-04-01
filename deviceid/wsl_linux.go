//go:build linux

package deviceid

import (
	"os"
	"strings"
)

// isWSL reports Linux running as WSL1/WSL2 (not a generic Linux VM).
func isWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}
	// WSLg / interop
	if os.Getenv("WSL_INTEROP") != "" {
		return true
	}
	// Strong signal for default WSL2 installs
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
		return true
	}
	if b, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		s := strings.ToLower(strings.TrimSpace(string(b)))
		if strings.Contains(s, "microsoft-standard-wsl") || strings.Contains(s, "wsl2") {
			return true
		}
	}
	// WSL1: kernel version string includes "Microsoft"
	if b, err := os.ReadFile("/proc/version"); err == nil {
		lo := strings.ToLower(string(b))
		if strings.Contains(lo, "microsoft") {
			return true
		}
	}
	return false
}
