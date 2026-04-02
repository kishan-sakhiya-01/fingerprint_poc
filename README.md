# fingerprint_poc

Go library and small CLI that compute a **device fingerprint** and a **SHA-256 hash** intended to identify a **machine** consistently across common variation: **dual-boot OS**, **WSL vs native Windows**, and **cloud instances**.

The hash is derived only from stable fields (`v`, `source`, `id`). The printed fingerprint also includes `os` and `arch` for visibility; changing OS or architecture does **not** change the hash if `source` and `id` stay the same.

## Requirements

- [Go](https://go.dev/dl/) (version compatible with `go.mod`, currently 1.25+)

**Windows:** PowerShell **SMBIOS** (`Win32_ComputerSystemProduct`) first, then **boot disk serial** (CIM) if SMBIOS fails, then `MachineGuid`.

**WSL:** **`powershell.exe`** for host SMBIOS, then host boot-disk serial, then guest SMBIOS—[interop](https://learn.microsoft.com/en-us/windows/wsl/filesystems#run-windows-tools-from-wsl) required.

**Linux (native):** SMBIOS first (`product_uuid`, **`tables/DMI`**, `entries/*/raw`, `dmidecode`). If that whole chain fails, **boot disk serial** from **`/sys/block/...`** (no `sudo` when sysfs is readable), then `machine-id`.

## Quick start

```bash
go run .
```

The program prints two lines: the hex-encoded SHA-256 digest, then the `Fingerprint` struct.

Use the library from another module:

```go
import "github.com/kishan-sakhiya-01/fingerprint_poc/deviceid"

fp, hash, err := deviceid.Compute()
```

## Fingerprint and hash

- **`Fingerprint`**: `v` (schema/version), `source` (how `id` was obtained), `id` (normalized identifier string), plus `os` and `arch` from the Go runtime.
- **`hash`**: `SHA256( JSON({ "v", "source", "id" }) )` as lowercase hex. Field order follows JSON marshaling of that struct. Schema `v` is currently **6** (SMBIOS preferred; boot disk is fallback so sudo/readable DMI matches Windows again and avoids WMI-vs-sysfs disk serial skew).

Hardware UUIDs are normalized: trimmed, braces removed, lowercased, and the all-zero UUID is rejected. For SMBIOS product/system UUIDs, the first three fields are also **endianness-canonicalized** (the same firmware UUID is often shown differently by Windows WMI vs Linux `/sys/class/dmi/id/product_uuid`); the code picks a single deterministic string so dual-boot pairs match.

## How the identifier is chosen

Resolution order:

1. **AWS EC2** — instance id from IMDS (`tryAWS`).
2. **Google Compute Engine** — instance id from metadata (`tryGCP`).
3. **Azure VM** — `vmId` from Instance Metadata Service (`tryAzure`).
4. **Physical / local** — platform-specific (see below).

Cloud checks use short HTTP timeouts; if no cloud metadata responds, the physical path is used.

### Physical / local (`physicalIdentity`)

| Environment | Primary identifier | `source` value |
|-------------|-------------------|----------------|
| **Windows** | SMBIOS UUID via PowerShell | `smbios_system_uuid` |
| Windows (fallback) | Boot disk serial (CIM), normalized | `boot_disk_serial` |
| Windows (last resort) | `MachineGuid` | `windows_machine_guid` |
| **WSL** | Host SMBIOS CIM, then host boot-disk PS, then guest SMBIOS | `smbios_system_uuid` or `boot_disk_serial` |
| **Linux (native)** | SMBIOS chain (see Requirements) | `smbios_system_uuid` |
| Linux (fallback) | Boot disk sysfs serial | `boot_disk_serial` |
| Linux (last resort) | `machine-id` | `linux_machine_id` |
| **macOS** | `ioreg` `IOPlatformUUID` | `darwin_platform_uuid` |
| Other | — | `unsupported_os` (error) |

WSL avoids using the **guest** virtual disk for `boot_disk_serial` until host PowerShell paths are exhausted.

**No `sudo` and locked-down DMI:** Linux may only get `boot_disk_serial` while Windows still gets SMBIOS—they will **not** match until DMI is readable on Linux or you standardize on disk-only (not implemented as default).

## Design goals and limits

**Goals:**

- **SMBIOS**-based `id` should match between **Windows** and **Ubuntu** when Linux can read DMI (often with `sudo` for `dmidecode`, or readable **`tables/DMI`** / `product_uuid`).
- Switching OS or architecture should not change the hash if the underlying machine id and source are unchanged.

**Not guaranteed:**

- **Different machines** or **different VMs** should produce different fingerprints; there is no global “one id for every possible runtime.”
- **Containers** often hide or fake DMI; IDs may be unstable or generic.
- **WSL without interop** or blocked PowerShell may fall back to Linux-only ids and diverge from Windows.
- **Spoofing**, cloned disks, or firmware that reports missing or all-zero UUIDs can break or weaken stability.
- **Privacy**: the fingerprint may constitute personal/device data; handle it like any other identifier under your policies.

## Project layout

- `deviceid/` — `Compute`, cloud probes, per-OS physical resolution, SMBIOS/CIM helper, WSL detection, UUID normalization.
- `main.go` — sample CLI.

## License

Use and licensing are determined by the repository owner; add a `LICENSE` file if you distribute this project.
