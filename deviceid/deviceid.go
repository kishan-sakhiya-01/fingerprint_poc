package deviceid

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"runtime"
)

// fingerprintVersion 4: v3 + Linux reads SMBIOS Type-1 from firmware DMI entries (and dmidecode fallback) when product_uuid sysfs is missing.
const fingerprintVersion = 4

type Source string

const (
	SourceAWS           Source = "aws_ec2"
	SourceAzure         Source = "azure_vm"
	SourceGCP           Source = "gcp_gce"
	SourceSMBIOS        Source = "smbios_system_uuid"
	SourceLinux         Source = "linux_machine_id"
	SourceWindows       Source = "windows_machine_guid"
	SourceDarwin        Source = "darwin_platform_uuid"
	SourceUnsupportedOS Source = "unsupported_os"
)

type Identity struct {
	Source Source `json:"source"`
	RawID  string `json:"raw_id"`
}

type Fingerprint struct {
	Version int    `json:"v"`
	Source  Source `json:"source"`
	ID      string `json:"id"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
}

// Compute returns a versioned fingerprint and SHA-256 hex hash.
func Compute() (Fingerprint, string, error) {
	id, err := resolveIdentity()
	if err != nil {
		return Fingerprint{}, "", err
	}
	fp := Fingerprint{
		Version: fingerprintVersion,
		Source:  id.Source,
		ID:      id.RawID,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
	hash, err := hashFingerprint(fp)
	if err != nil {
		return Fingerprint{}, "", err
	}
	return fp, hash, nil
}

func resolveIdentity() (Identity, error) {
	if id, ok := tryAWS(); ok {
		return id, nil
	}
	if id, ok := tryGCP(); ok {
		return id, nil
	}
	if id, ok := tryAzure(); ok {
		return id, nil
	}
	return physicalIdentity()
}

// stableHashPayload is what gets hashed so dual-boot OS (or arch) changes do not change the digest.
type stableHashPayload struct {
	Version int    `json:"v"`
	Source  Source `json:"source"`
	ID      string `json:"id"`
}

func hashFingerprint(fp Fingerprint) (string, error) {
	payload := stableHashPayload{
		Version: fp.Version,
		Source:  fp.Source,
		ID:      fp.ID,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
