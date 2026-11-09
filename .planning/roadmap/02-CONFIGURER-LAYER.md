# 02: Configurer Layer — OS Abstraction

## Overview

The configurer layer provides OS-specific command builders and path resolvers. It uses a registry pattern where each OS module registers itself via `init()` with a matcher function.

## Files

### `configurer/interface.go`
- Defines the `Configurer` interface with K3s-specific methods:
  - `K3sCmdf` — builds `k3s <cmd>` commands
  - `K3sBinaryPath` — returns `/usr/local/bin/k3s`
  - `K3sConfigPath` — returns `/etc/rancher/k3s/config.yaml`
  - `K3sJoinTokenPath` — returns `/etc/rancher/k3s/token`
  - `DataDirDefaultPath` — returns `/var/lib/rancher/k3s`
  - `K3sctlLockFilePath` — returns `/run/lock/k3sctl`
  - `KubectlCmdf` — uses `k3s kubectl`
  - `KubeconfigPath` — returns `/etc/rancher/k3s/k3s.yaml`
  - Keep: `Kind()`, `OSKind()`, `UpdateEnvironment()`, `FixContainer()`, `PrivateInterface()`, `PrivateAddress()`, `SetPath()`
- Also includes `HostValidator` interface with `ValidateHost(Host) error`

### `configurer/host.go`
- Same interface wrapping rig's SimpleRunner, Sudo, FS

### `configurer/registry.go`
- Same `RegisterOSModule`/`ResolveOSModule` pattern with mutex-protected registry

### `configurer/linux.go` (~200 lines)
- The main implementation file. All K3s paths:

| Method | Value |
|--------|-------|
| `Kind()` | `"linux"` |
| `OSKind()` | `"linux"` |
| `K3sBinaryPath()` | `/usr/local/bin/k3s` |
| `K3sConfigPath()` | `/etc/rancher/k3s/config.yaml` |
| `DataDirDefaultPath()` | `/var/lib/rancher/k3s` |
| `K3sJoinTokenPath()` | `/etc/rancher/k3s/token` |
| `K3sctlLockFilePath()` | `/run/lock/k3sctl` |
| `K3sCmdf()` | `/usr/local/bin/k3s <cmd>` |
| `KubectlCmdf()` | `k3s kubectl --data-dir=...` |
| `KubeconfigPath()` | `/etc/rancher/k3s/k3s.yaml` |

- `PrivateInterface()` — uses `ip route list`
- `PrivateAddress()` — uses `ip -o addr show`
- `UpdateEnvironment()` — patches `/etc/environment`
- `FixContainer()` — `mount --make-rshared /`
- `SetPath()` — dynamic path override

### No Windows modules — K3s is Linux-only

### Linux Distro Modules (all follow same pattern)

Each file: registers via `init()`, embeds `configurer.Linux` (directly or via a base type), optionally overrides `Prepare()` for package installation.

| File | Notes |
|------|-------|
| `configurer/linux/alpine.go` | Alpine Linux |
| `configurer/linux/archlinux.go` | Arch Linux |
| `configurer/linux/coreos.go` | CoreOS/Flatcar |
| `configurer/linux/debian.go` | Debian |
| `configurer/linux/enterpriselinux.go` | Base type for RHEL family |
| `configurer/linux/enterpriselinux/almalinux.go` | AlmaLinux |
| `configurer/linux/enterpriselinux/amazon.go` | Amazon Linux |
| `configurer/linux/enterpriselinux/centos.go` | CentOS |
| `configurer/linux/enterpriselinux/fedora.go` | Fedora |
| `configurer/linux/enterpriselinux/oracle.go` | Oracle Linux |
| `configurer/linux/enterpriselinux/rhel.go` | RHEL |
| `configurer/linux/enterpriselinux/rocky.go` | Rocky Linux |
| `configurer/linux/flatcar.go` | Overrides `K3sBinaryPath` → `/opt/bin/k3s` |
| `configurer/linux/opensuse.go` | openSUSE |
| `configurer/linux/sles.go` | SLES |
| `configurer/linux/ubuntu.go` | Ubuntu |

### Test Files

| File | Notes |
|------|-------|
| `configurer/linux/linux_test.go` | Test K3s paths |
| `configurer/linux/alpine_test.go` | Alpine-specific tests |

## Key Implementation Notes

1. The `SetPath` mechanism allows runtime override of `K3sBinaryPath` — essential for custom install locations
2. `KubeconfigPath` is always `/etc/rancher/k3s/k3s.yaml` on the server node (no fallback logic needed)
3. `K3sCmdf` builds commands like `/usr/local/bin/k3s server`, `/usr/local/bin/k3s agent`, etc.
4. No Windows modules needed — K3s is Linux-only
