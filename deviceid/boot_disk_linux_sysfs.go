//go:build linux

package deviceid

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reNVMePart = regexp.MustCompile(`^(nvme\d+n\d+)p\d+$`)
	reMMCPart  = regexp.MustCompile(`^(mmcblk\d+)p\d+$`)
	reSDPart   = regexp.MustCompile(`^(sd[a-z]+)\d+$`)
)

func tryLinuxBootDiskSerialSysfs() (Identity, bool) {
	src, err := rootMountDeviceFromProcMounts()
	if err != nil || src == "" {
		return Identity{}, false
	}
	if strings.HasPrefix(src, "/dev/mapper/") {
		if lp, err := os.Readlink(src); err == nil {
			src = filepath.Join("/dev", filepath.Base(lp))
		}
	}
	parent := sysfsDiskBase(strings.TrimPrefix(src, "/dev/"))
	if parent == "" {
		return Identity{}, false
	}
	raw, ok := readLinuxBlockSerialRaw(parent)
	if !ok {
		return Identity{}, false
	}
	s, ok := normalizeBootDiskSerial(raw)
	if !ok {
		return Identity{}, false
	}
	return Identity{Source: SourceBootDisk, RawID: s}, true
}

func rootMountDeviceFromProcMounts() (string, error) {
	b, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return "", err
	}
	var last string
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) < 3 || f[1] != "/" {
			continue
		}
		dev := f[0]
		if dev == "rootfs" || dev == "overlay" || strings.HasPrefix(dev, "/dev/loop") {
			continue
		}
		if f[2] == "tmpfs" || f[2] == "devtmpfs" {
			continue
		}
		if !strings.HasPrefix(dev, "/dev/") {
			continue
		}
		last = dev
	}
	if last == "" {
		return "", os.ErrNotExist
	}
	return last, nil
}

func sysfsDiskBase(name string) string {
	base := filepath.Base(name)
	for strings.HasPrefix(base, "dm-") {
		slavesDir := filepath.Join("/sys/block", base, "slaves")
		ents, err := os.ReadDir(slavesDir)
		if err != nil || len(ents) == 0 {
			break
		}
		base = ents[0].Name()
	}
	if m := reNVMePart.FindStringSubmatch(base); m != nil {
		return m[1]
	}
	if m := reMMCPart.FindStringSubmatch(base); m != nil {
		return m[1]
	}
	if m := reSDPart.FindStringSubmatch(base); m != nil {
		return m[1]
	}
	i := len(base) - 1
	for i >= 0 && base[i] >= '0' && base[i] <= '9' {
		i--
	}
	if i >= 0 && i < len(base)-1 {
		return base[:i+1]
	}
	return base
}

func readLinuxBlockSerialRaw(parent string) (string, bool) {
	candidates := []string{
		filepath.Join("/sys/block", parent, "queue", "serial_number"),
		filepath.Join("/sys/block", parent, "device", "serial"),
		filepath.Join("/sys/block", parent, "device", "wwid"),
	}
	for i := range candidates {
		b, err := os.ReadFile(candidates[i])
		if err != nil {
			continue
		}
		s := strings.TrimSpace(string(b))
		if s != "" && !strings.EqualFold(s, "none") {
			return s, true
		}
	}
	if ctrl := nvmeCtrlFromParent(parent); ctrl != "" {
		b, err := os.ReadFile(filepath.Join("/sys/class/nvme", ctrl, "serial"))
		if err == nil {
			s := strings.TrimSpace(string(b))
			if s != "" {
				return s, true
			}
		}
	}
	return "", false
}

func nvmeCtrlFromParent(parent string) string {
	// parent like nvme0n1 -> nvme0
	i := strings.LastIndex(parent, "n")
	if i <= 0 || !strings.HasPrefix(parent, "nvme") {
		return ""
	}
	return parent[:i]
}
