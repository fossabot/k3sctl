# 06: Core Phases — Connection, Detection, Facts

## Phase: Connect (`phase/connect.go`)
- Connects to all hosts in parallel via rig v2
- Auth rejection limit (12 attempts ~60 seconds)
- Retries with 10-minute timeout
- Handles `rig.ErrNonRetryable`, `host key mismatch`, `rig.ErrAuthFailed`
- SSH connection is tool-agnostic

## Phase: DetectOS (`phase/detect_os.go`)
- Blank imports for OS modules:
  ```go
  _ "github.com/nirantarasi/k3sctl/configurer"
  _ "github.com/nirantarasi/k3sctl/configurer/linux"
  _ "github.com/nirantarasi/k3sctl/configurer/linux/enterpriselinux"
  ```
- No Windows import
- `SetPath("K3sBinaryPath", ...)` for custom install path
- `After()` hook: runs "connect: after" hooks

## Phase: GatherFacts (`phase/gather_facts.go`)
- `K3sLeader()` — identify leader node
- `K3s.Version` — version tracking
- `K3sBinaryPath` / `K3sInstallPath` — binary location
- Architecture detection, hostname resolution, private interface/address
- `useExistingK3s` — reuse existing binary

### `phase/gather_facts_test.go`

## Phase: GatherK3sFacts (`phase/gather_k3s_facts.go`)
- Detect installed K3s binary version: `k3s --version`
- Detect running K3s version: check systemd service status
- Parse K3s version from output (format: `k3s version v1.30.0+k3s1 (hash)`)
- Read existing K3s config from `/etc/rancher/k3s/config.yaml`
- Detect if K3s is leader (first server node)
- Get cluster ID via `k3s kubectl`
- Parse status args from systemd service file

### `phase/gather_k3s_facts_test.go`

## Phase: ValidateFacts (`phase/validate_facts.go`)
- Downgrade check: compare K3s versions
- Version defaulting: if no version specified, use running version or latest
- `K3s` references throughout
- Version format: `v1.30.0+k3s1`

### `phase/validate_facts_test.go`

## Phase: ValidateHosts (`phase/validate_hosts.go`)
- Calls host.Validate() and configurer.ValidateHost()
- No Windows validation

### `phase/validate_hosts_test.go`

## Phase: ValidateEtcdMembers (`phase/validate_etcd_members.go`)
- K3s uses embedded SQLite by default, not etcd
  - If external etcd configured, validate members
  - If SQLite (default), skip entirely
  - If embedded etcd (`--cluster-init`), validate

## Phase: PrepareHosts (`phase/prepare_hosts.go`)
- Installs required packages via configurer.Prepare()
- Handles: curl, iptables (for agents), inetutils (for hostname)

## Phase: Lock/Unlock (`phase/lock.go`, `phase/unlock.go`)
- Lock file path uses `K3sctlLockFilePath`
- Advisory locking mechanism

## Phase: Disconnect (`phase/disconnect.go`)
- Generic disconnect all hosts
