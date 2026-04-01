# fingerprint_poc

Go library and small CLI that compute a **device fingerprint** and a **SHA-256 hash** intended to identify a **machine** consistently across common variation: **dual-boot OS**, **WSL vs native Windows**, and **cloud instances**.

The hash is derived only from stable fields (`v`, `source`, `id`). The printed fingerprint also includes `os` and `arch` for visibility; changing OS or architecture does **not** change the hash if `source` and `id` stay the same.

## Requirements

- [Go](https://go.dev/dl/) (version compatible with `go.mod`, currently 1.25+)

**Windows (SMBIOS via CIM):** PowerShell with `Get-CimInstance` (typical on Windows 10/11 and Server).

**WSL:** [Windows interoperability](https://learn.microsoft.com/en-us/windows/wsl/filesystems#run-windows-tools-from-wsl) should be enabled so `powershell.exe` can run from the Linux side; otherwise resolution may fall back to the distro’s Linux identifiers and **not** match the Windows host.

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
- **`hash`**: `SHA256( JSON({ "v", "source", "id" }) )` as lowercase hex. Field order follows JSON marshaling of that struct.

Hardware UUIDs are normalized: trimmed, braces removed, lowercased, and the all-zero UUID is rejected.

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
| **Windows** | `Win32_ComputerSystemProduct.UUID` (SMBIOS) via PowerShell | `smbios_system_uuid` |
| Windows (fallback) | Registry `HKLM\SOFTWARE\Microsoft\Cryptography` `MachineGuid` | `windows_machine_guid` |
| **Linux** | Same SMBIOS as Windows when running under **WSL**: host query via `powershell.exe` | `smbios_system_uuid` |
| Linux | `/sys/class/dmi/id/product_uuid` | `smbios_system_uuid` |
| Linux (fallback) | `/etc/machine-id` or `/var/lib/dbus/machine-id` | `linux_machine_id` |
| **macOS** | `ioreg` `IOPlatformUUID` | `darwin_platform_uuid` |
| Other | — | `unsupported_os` (error) |

WSL is detected via `WSL_DISTRO_NAME`, `WSL_INTEROP`, `WSLInterop`, WSL-style `osrelease`, or WSL1-style `/proc/version` branding. On WSL, the **Windows host** SMBIOS path is preferred so the result aligns with native Windows on the same PC.

## Design goals and limits

**Goals:**

- Same **bare-metal** machine under **Windows**, **WSL**, and **dual-boot Linux** should share the same `id` and hash when SMBIOS/product UUID is available and WSL can reach PowerShell.
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
