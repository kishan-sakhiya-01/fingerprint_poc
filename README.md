# fingerprint_poc

Go library and small CLI that compute a **device fingerprint** and a **SHA-256 hash** intended to identify a **machine** consistently across common variation: **dual-boot OS**, **WSL vs native Windows**, and **cloud instances**.

The hash is derived only from stable fields (`v`, `source`, `id`). The printed fingerprint also includes `os` and `arch` for visibility; changing OS or architecture does **not** change the hash if `source` and `id` stay the same.

## Requirements

- [Go](https://go.dev/dl/) (version compatible with `go.mod`, currently 1.25+)

**Windows:** PowerShell is used for the **boot disk serial** (CIM) first, then SMBIOS, then `MachineGuid`.

**WSL:** [interop](https://learn.microsoft.com/en-us/windows/wsl/filesystems#run-windows-tools-from-wsl) must allow **`powershell.exe`** so the **same** boot-disk script runs as on Windows; otherwise host parity breaks.

**Linux (native):** **`/sys/...` boot disk serial** for the root filesystem’s physical disk usually works **without sudo**, before any SMBIOS paths. Extra SMBIOS sources (`product_uuid`, **`/sys/firmware/dmi/tables/DMI`**, `entries/*/raw`, `dmidecode`) run when disk serial is unavailable.

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
- **`hash`**: `SHA256( JSON({ "v", "source", "id" }) )` as lowercase hex. Field order follows JSON marshaling of that struct. Schema `v` is currently **5** (boot-disk serial preferred for physical machines so Ubuntu without DMI/sudo can match Windows on the **same physical disk**).

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
| **Windows** | Serial of the disk that holds the **system volume** (CIM), normalized | `boot_disk_serial` |
| Windows (fallback) | SMBIOS UUID via PowerShell | `smbios_system_uuid` |
| Windows (last resort) | `MachineGuid` | `windows_machine_guid` |
| **WSL** | PowerShell host boot-disk serial, then host SMBIOS CIM, then guest SMBIOS | `boot_disk_serial` or `smbios_system_uuid` |
| **Linux (native)** | **`/sys/block/...`** serial for **root’s disk** (LVM → `dm-*` slaves), normalized | `boot_disk_serial` |
| Linux | SMBIOS: `product_uuid`, **`tables/DMI`**, **`entries/*/raw`**, `dmidecode` | `smbios_system_uuid` |
| Linux (fallback) | `machine-id` | `linux_machine_id` |
| **macOS** | `ioreg` `IOPlatformUUID` | `darwin_platform_uuid` |
| Other | — | `unsupported_os` (error) |

WSL is detected the same way as before. **Native** Linux does **not** use a guest virtual disk for `boot_disk_serial`; WSL uses the host via PowerShell first.

Note: **Windows and Ubuntu must share the same physical boot disk** for `boot_disk_serial` to match. If each OS is on a different drive, ids will differ by design.

## Design goals and limits

**Goals:**

- Same **bare-metal** machine under **Windows**, **WSL**, and **dual-boot Linux** on the **same physical system disk** should share the same `id` and hash via **`boot_disk_serial`** first (no `sudo` on typical Ubuntu).
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
