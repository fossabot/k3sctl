# Roadmap: Configurer Layer (OS-Specific Behavior)

## Architecture

The configurer layer provides OS-specific behavior through a registry pattern:
1. `configurer/interface.go` defines the `Configurer` interface
2. `configurer/registry.go` provides `RegisterOSModule` / `ResolveOSModule`
3. OS modules register via `init()` functions
4. Each module matches by OS release info and provides path/command implementations

## Source Files

### configurer/interface.go
- `Configurer` interface methods:
  - `Kind() string` — OS identifier
  - `OSKind() string` — "linux"
  - `K3sCmdf(format, args...) string` — format K3s command string
  - `K3sBinaryPath() string` — path to K3s binary (`/usr/local/bin/k3s`)
  - `K3sConfigPath() string` — path to K3s config file (`/etc/rancher/k3s/config.yaml`)
  - `DataDirDefaultPath() string` — default data directory (`/var/lib/rancher/k3s`)
  - `K3sJoinTokenPath() string` — path to join token file (`/var/lib/rancher/k3s/server/node-token`)
  - `UpdateEnvironment(map[string]string) string` — update env vars
  - `KubectlCmdf(host, dataDir, format, args...) string` — format kubectl command
  - `KubeconfigPath(host, dataDir) string` — kubeconfig file path (`/etc/rancher/k3s/k3s.yaml`)
  - `FixContainer() string` — container fix command
  - `PrivateInterface() string` — command to detect private interface
  - `PrivateAddress() string` — command to detect private address
  - `K3sctlLockFilePath() string` — lock file path (`/var/lock/k3sctl.lock`)
  - `SetPath(key, value)` — set configurer path
- `HostValidator` interface:
  - `ValidateHost(host, dataDir) error` — validate host for compatibility

### configurer/host.go
- `Host` interface: requires `cmd.SimpleRunner`, `Sudo() *rig.Client`, `FS() remotefs.FS`

### configurer/registry.go
- `RegisterOSModule(matcher, builder)` — register OS module with matcher function
- `ResolveOSModule(osRelease) (Configurer, error)` — resolve OS module from os-release info
- Matcher: `func(os *rig.OSRelease) bool`
- Builder: `func() Configurer`

### configurer/linux.go
- `Linux` struct with paths map
- Default paths:
  - `K3sBinaryPath` = `/usr/local/bin/k3s`
  - `K3sConfigPath` = `/etc/rancher/k3s/config.yaml`
  - `K3sJoinTokenPath` = `/var/lib/rancher/k3s/server/node-token`
  - `DataDirDefaultPath` = `/var/lib/rancher/k3s`
- Methods:
  - `K3sCmdf`: formats command with binary path
  - `KubectlCmdf`: uses `k3s kubectl` subcommand
  - `KubeconfigPath`: `/etc/rancher/k3s/k3s.yaml`
  - `PrivateInterface`: uses `ip route` to detect
  - `PrivateAddress`: uses `ip -o addr show` to detect
  - `UpdateEnvironment`: patches `/etc/environment`
  - `FixContainer`: runs `mount --make-rshared /`

## Linux Distribution Modules

### configurer/linux/alpine.go
- Alpine Linux configurer
- Matcher: ID="alpine"
- Handles: apk package manager, OpenRC service manager

### configurer/linux/debian.go
- Debian configurer
- Matcher: ID="debian"

### configurer/linux/ubuntu.go
- Ubuntu configurer (extends Debian)
- Matcher: ID="ubuntu"

### configurer/linux/enterpriselinux.go
- Base for RHEL-family distros

### configurer/linux/enterpriselinux/rhel.go
- RHEL configurer

### configurer/linux/enterpriselinux/centos.go
- CentOS configurer

### configurer/linux/enterpriselinux/fedora.go
- Fedora configurer

### configurer/linux/enterpriselinux/rocky.go
- Rocky Linux configurer

### configurer/linux/enterpriselinux/almalinux.go
- AlmaLinux configurer

### configurer/linux/enterpriselinux/oracle.go
- Oracle Linux configurer

### configurer/linux/enterpriselinux/amazon.go
- Amazon Linux configurer

### configurer/linux/archlinux.go
- Arch Linux configurer

### configurer/linux/opensuse.go
- openSUSE configurer

### configurer/linux/sles.go
- SUSE Linux Enterprise Server configurer

### configurer/linux/flatcar.go
- Flatcar Container Linux configurer

### configurer/linux/coreos.go
- CoreOS configurer

## Test Files
- configurer/linux/alpine_test.go
- configurer/linux/linux_test.go

## Implementation Order
1. `configurer/interface.go` (K3s interface)
2. `configurer/host.go`
3. `configurer/registry.go`
4. `configurer/linux.go` (K3s paths)
5. Linux distro modules (one by one):
   - debian.go, ubuntu.go (most common)
   - enterpriselinux.go + rhel.go, centos.go, rocky.go, almalinux.go
   - fedora.go, amazon.go, oracle.go
   - alpine.go, archlinux.go
   - opensuse.go, sles.go
   - flatcar.go, coreos.go
6. Tests
