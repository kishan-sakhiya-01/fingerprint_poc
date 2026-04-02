//go:build windows || linux

package deviceid

import "os/exec"

// Serial of the Windows system volume backing disk (physical disk), via CIM.
const psBootDiskSerialScript = `
$ErrorActionPreference = 'Stop'
$ld = ($env:SystemDrive.TrimEnd(':') + ':')
$p = Get-CimInstance -Query "ASSOCIATORS OF {Win32_LogicalDisk.DeviceID='$ld'} WHERE AssocClass=Win32_LogicalDiskToPartition" | Select-Object -First 1
if (-not $p) { exit 1 }
$d = Get-CimInstance -Query "ASSOCIATORS OF {Win32_DiskPartition.DeviceID='$($p.DeviceID)'} WHERE AssocClass=Win32_DiskPartitionToDiskDrive" | Select-Object -First 1
if (-not $d) { exit 1 }
Write-Output $d.SerialNumber.Trim()
`

func runPowerShellBootDiskSerial() (string, bool) {
	ps := resolvePowerShellForCIM()
	if ps == "" {
		return "", false
	}
	cmd := exec.Command(ps, "-NoProfile", "-Command", psBootDiskSerialScript)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	s, ok := normalizeBootDiskSerial(string(out))
	if !ok {
		return "", false
	}
	return s, true
}

func identityBootDiskFromPS() (Identity, bool) {
	s, ok := runPowerShellBootDiskSerial()
	if !ok {
		return Identity{}, false
	}
	return Identity{Source: SourceBootDisk, RawID: s}, true
}
