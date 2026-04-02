//go:build windows || linux

package deviceid

import (
	"os"
	"os/exec"
	"runtime"
)

// identityFromCIMProductUUID reads Win32_ComputerSystemProduct.UUID via Windows PowerShell.
// Used on native Windows and on WSL (Linux) so the ID matches the Windows host.
func identityFromCIMProductUUID() (identity, bool) {
	ps := resolvePowerShellForCIM()
	if ps == "" {
		return identity{}, false
	}
	cmd := exec.Command(ps, "-NoProfile", "-Command", "(Get-CimInstance -ClassName Win32_ComputerSystemProduct).UUID")
	out, err := cmd.Output()
	if err != nil {
		return identity{}, false
	}
	u, ok := stabilizeSMBIOSUUID(string(out))
	if !ok {
		return identity{}, false
	}
	return identity{source: SourceSMBIOS, rawID: u}, true
}

func resolvePowerShellForCIM() string {
	if runtime.GOOS == "windows" {
		if p, err := exec.LookPath("powershell.exe"); err == nil {
			return p
		}
		return "powershell.exe"
	}
	if p, err := exec.LookPath("powershell.exe"); err == nil {
		return p
	}
	const winPS = "/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe"
	if st, err := os.Stat(winPS); err == nil && !st.IsDir() {
		return winPS
	}
	return ""
}
