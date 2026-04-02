# fingerprint_poc

Go library and small CLI that compute a **device fingerprint** and a **SHA-256 hash** intended to identify a **machine** consistently across common variation: **dual-boot OS**, **WSL vs native Windows**, and **cloud instances**.

The hash is derived only from stable fields (`v`, `source`, `id`). The printed fingerprint also includes `os` and `arch` for visibility; changing OS or architecture does **not** change the hash if `source` and `id` stay the same.

## Requirements

- [Go](https://go.dev/dl/) (version compatible with `go.mod`, currently 1.25+)

**Windows:** PowerShell with `Get-CimInstance` for SMBIOS (typical on Windows 10/11 and Server).

**WSL:** [Windows interoperability](https://learn.microsoft.com/en-us/windows/wsl/filesystems#run-windows-tools-from-wsl) should be enabled so `powershell.exe` can run from the Linux side; otherwise resolution may fall back to the distro’s Linux identifiers and **not** match the Windows host.

**Linux (native):** SMBIOS via `product_uuid`, **`/sys/firmware/dmi/tables/DMI`**, **`entries/*/raw`**, and optionally `dmidecode`. Many paths work without `sudo`; some systems lock down DMI (see below).

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
- **`hash`**: `SHA256( JSON({ "v", "source", "id" }) )` as lowercase hex. Field order follows JSON marshaling of that struct. Schema `v` is currently **4** (SMBIOS focus + Linux DMI fallbacks for parity with Windows when `sudo` / readable DMI works).

Hardware UUIDs are normalized: trimmed, braces removed, lowercased, and the all-zero UUID is rejected. For SMBIOS product/system UUIDs, the first three fields are also **endianness-canonicalized** (the same firmware UUID is often shown differently by Windows WMI vs Linux sysfs); the code picks a single deterministic string so dual-boot pairs match.

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
| Windows (fallback) | Registry `MachineGuid` | `windows_machine_guid` |
| **WSL** | Host SMBIOS via `powershell.exe`, then guest SMBIOS chain | `smbios_system_uuid` |
| **Linux (native)** | SMBIOS (`product_uuid`, **`tables/DMI`**, **`entries/*/raw`**, `dmidecode`) | `smbios_system_uuid` |
| Linux (fallback) | `/etc/machine-id` or `/var/lib/dbus/machine-id` | `linux_machine_id` |
| **macOS** | `ioreg` `IOPlatformUUID` | `darwin_platform_uuid` |
| Other | — | `unsupported_os` (error) |

If DMI is not readable on Linux without elevated privileges, the app may fall back to **`linux_machine_id`**, which will **not** match Windows SMBIOS until SMBIOS data is available (e.g. readable sysfs, or `dmidecode` with appropriate permissions).

## Design goals and limits

**Goals:**

- **SMBIOS**-based `id` should match between **Windows** and **Ubuntu** when Linux can read the same firmware UUID (WSL uses host PowerShell; native Linux uses DMI sysfs / table walk / `dmidecode`).
- Switching OS or architecture should not change the hash if the underlying machine `id` and `source` are unchanged.

**Not guaranteed:**

- **Different machines** or **different VMs** should produce different fingerprints; there is no global “one id for every possible runtime.”
- **Containers** often hide or fake DMI; IDs may be unstable or generic.
- **WSL without interop** or blocked PowerShell may fall back to Linux-only ids and diverge from Windows.
- **Spoofing**, cloned disks, or firmware that reports missing or all-zero UUIDs can break or weaken stability.
- **Privacy**: the fingerprint may constitute personal/device data; handle it like any other identifier under your policies.

## Project layout

- `deviceid/` — `Compute`, cloud probes, per-OS physical resolution, SMBIOS/CIM helper, WSL detection, DMI table walk, UUID normalization.
- `main.go` — sample CLI.

## License

Use and licensing are determined by the repository owner; add a `LICENSE` file if you distribute this project.
