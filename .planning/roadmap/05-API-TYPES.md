# 05: API Types — Cluster Config Model

## Overview

The API types define the YAML configuration format. This layer defines K3s-specific concepts, paths, and behaviors.

## API Version

- `k3sctl.k3sproject.io/v1beta1`

## Files

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster.go`
- `APIVersion = "k3sctl.k3sproject.io/v1beta1"`
- `ClusterMetadata` — fields: Name, User, Kubeconfig, EtcdMembers, Manifests
- `Cluster` struct — APIVersion, Kind, Metadata, Spec, Origin
- `StorageType()` — K3s defaults to "sqlite"; checks config for `datastore-endpoint`
- `Resolve(baseDir)` — delegates to Spec

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/spec.go`
- `Spec` struct: `Hosts`, `K3s`, `Options`
- `K3sLeader()` → returns first server node
- `KubeAPIURL()` — API URL from leader
- `InternalKubeAPIURL()` — internal API URL
- `NodeInternalKubeAPIURL(h)` — per-node internal API URL
- `APIPort()` — checks config for `https-listen-port`, defaults to 6443
- `clusterExternalAddress()` — checks config for `tls-san`

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/host.go` (~957 lines, largest file)

| Field/Method | Description |
|-------------|--------|
| `Role` | server, agent |
| `K3sBinaryPath` | Path to K3s binary |
| `K3sInstallPath` | Custom install location |
| `K3sDownloadURLOverride` | Custom download URL |
| `K3sServiceName()` | Returns `k3s` or `k3s-agent` |
| `K3sRole()` | Returns server/agent |
| `K3sInstallCommand()` | Systemd unit installation |
| `K3sResetCommand()` | `k3s-uninstall.sh` or `k3s-agent-uninstall.sh` |
| `K3sBinaryVersion()` | `k3s --version` |
| `K3sDataDir()` | Data dir with default |
| `K3sInstallFlags()` | Build K3s server/agent flags |
| `IsServer()` | Checks server role |
| `DrainNode()` | `k3s kubectl drain` |
| `Connect()` | rig v2 SSH connection |
| `ResolveConfigurer()` | OS module resolution |

- **K3s Install approach**:
  - Place binary, create systemd unit, start service
  - Flags: `--data-dir`, `--token`, `--server` (for agents), `--config`

- **K3s Join Token**:
  - Stored at `/var/lib/rancher/k3s/server/node-token`
  - Agents use `--token` flag
  - Read token from server node (no `token create` command)

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/hosts.go`
- `Servers()` — filter by server role
- `Agents()` — filter by agent role
- `ParallelEach`, `BatchedParallelEach`, `Each`, `Filter`, `First`

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/k3s.go`

```go
type K3s struct {
    Version        *version.Version `yaml:"version,omitempty"`
    VersionChannel string           `yaml:"versionChannel,omitempty"`
    Config         dig.Mapping      `yaml:"config,omitempty"`  // K3s config.yaml content
    Metadata       K3sMetadata      `yaml:"-"`
}

type K3sMetadata struct {
    ClusterID        string
    VersionDefaulted bool
}
```

- `GenerateToken()` — read `/var/lib/rancher/k3s/server/node-token` from leader
- `ParseToken()` — K3s tokens are plain strings
- `GetClusterID()` — kubectl approach
- `NodeConfig()` — K3s config structure
- `K3sMinVersion` — earliest supported K3s version
- No `DynamicConfig` — K3s uses config file or CLI flags
- Version format: `v1.30.0+k3s1`

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/options.go`
- Options, WaitOption, DrainOption, ConcurrencyOption, EvictTaintOption
- `DrainOption.ToKubectlArgs()` — standard kubectl drain
- UnmarshalYAML/Validate methods

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/flags.go`
- Generic flag manipulation

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/hook.go`
- Generic hook definitions

### `pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster/uploadfile.go`
- Generic file upload

## Test Files

- `cluster_test.go`
- `host_test.go`
- `hosts_test.go`
- `k3s_test.go`
- `flags_test.go`, `flags_quoting_test.go`
- `spec_test.go`
- `uploadfile_test.go`
- `options_test.go`
