# Roadmap: API Types (Config Data Model)

## Package Path
`pkg/apis/k3sctl.k3sproject.io/v1beta1/`

## Source Files

### cluster.go (top-level)
- `APIVersion` constant = `"k3sctl.k3sproject.io/v1beta1"`
- `ClusterMetadata` struct: Name, User, Kubeconfig
- `Cluster` struct: APIVersion, Kind, Metadata, Spec, Origin
- Methods: UnmarshalYAML, String, SetDefaults, Validate, StorageType, Resolve
- `StorageType()`: check for etcd vs sqlite in K3s config
- Keep Manifests for extra manifest support

### cluster_test.go
- Validation tests, StorageType tests with K3s-specific storage type detection

### cluster/spec.go
- `Spec` struct: Hosts, K3s(*K3s), Options
- Methods:
  - `K3sLeader()` — picks first running server or falls back
  - `KubeAPIURL()`, `InternalKubeAPIURL()` — construct API URLs
  - `APIPort()` — defaults to 6443
  - `Resolve()` — resolve upload file paths
  - `NodeInternalKubeAPIURL` for per-node API access

### cluster/spec_test.go
- Tests for API URL construction, leader selection

### cluster/host.go (LARGE FILE — key type)
- `Host` struct embeds `rig.CompositeConfig` + `*rig.Client`
- Key fields:
  - Role: "server", "agent" (K3s terminology)
  - Reset, PrivateInterface, PrivateAddress
  - DataDir, KubeletRootDir
  - Environment map
  - UploadBinary, UseExistingK3s, K3sBinaryPath
  - K3sInstallPath, K3sDownloadURLOverride
  - InstallFlags (Flags type)
  - Files ([]*UploadFile)
  - OSIDOverride, HostnameOverride, NoTaints
  - Hooks (map of action/stage to commands)
  - Metadata (HostMetadata)
  - Configurer
  - OSRelease
- `HostMetadata`:
  - K3sBinaryVersion, K3sRunningVersion
  - K3sInstalled, K3sExistingConfig, K3sNewConfig
  - K3sTokenData, K3sStatusArgs
  - Arch, IsK3sLeader, Hostname, Ready
  - NeedsUpgrade, MachineID, DryRunFakeLeader
- Key methods:
  - `Connect(ctx)` — establish rig connection
  - `K3sInstallCommand()` — build systemd service file or equivalent install command
  - `K3sResetCommand()` — stop service, remove dirs
  - `K3sBackupCommand()` — `k3s etcd-snapshot save`
  - `K3sRestoreCommand()` — restore from etcd snapshot
  - `IsServer()` — checks server role
  - `K3sServiceName()` — returns "k3s" (server) or "k3s-agent" (agent)
  - `InstallK3sBinary()` — install binary at path
  - `K3sBinaryVersion()` — get version from binary
  - `UpdateK3sBinary()` — replace binary
  - `K3sDataDir()` — resolve data directory (`/var/lib/rancher/k3s` or custom)
  - `DrainNode()`, `CordonNode()`, `UncordonNode()`, `DeleteNode()`
  - `Taints()`, `AddTaint()`, `RemoveTaint()`
  - `CheckHTTPStatus()` — HTTP health check
  - `ExpandTokens()` — expand version/arch tokens in strings
  - `FlagsChanged()` — detect install flag changes
  - `RunHooks()` — execute hooks
  - Token handling: read `/var/lib/rancher/k3s/server/node-token` file

### cluster/host_test.go, cluster/host_posix_test.go
- Tests for service names, install commands, flag changes

### cluster/hosts.go
- `Hosts` type: `[]*Host`
- Methods: Validate, First, Last, Find, Filter, WithRole, Servers, Agents
- Each, ParallelEach, BatchedParallelEach
- `Servers()` — filter `h.IsServer()`
- `Agents()` — filter `h.Role == "agent"`
- Validation: at least one server required

### cluster/hosts_test.go
- Tests for host collection methods

### cluster/k3s.go
- `K3s` struct: Version, VersionChannel, Config(dig.Mapping), Metadata
- `K3sMetadata`: ClusterID, VersionDefaulted
- Key methods:
  - `ReadNodeToken()` — read from `/var/lib/rancher/k3s/server/node-token`
  - `GetClusterID()` — via kubectl
  - `ServerConfig()` / `AgentConfig()` — generate role-specific config YAML
- No DynamicConfig (K3s is file-based)
- Tokens are plain strings (no encode/decode)

### cluster/k3s_test.go
- K3s version parsing tests

### cluster/flags.go
- `Flags` type: string slice with flag manipulation
- Methods: Add, AddUnlessExist, AddOrReplace, Include, Index, Get, GetValue, GetBoolean, Delete, Merge, MergeOverwrite, Join, Each, Map, Equals, NewFlags
- Uses `shellescape.Unquote/Split` for safe flag parsing

### cluster/flags_test.go, cluster/flags_quoting_test.go
- Flag parsing and quoting tests

### cluster/options.go
- `Options` struct: Wait, Drain, Concurrency, EvictTaint
- `WaitOption`: Enabled bool pointer with default true
- `DrainOption`: Enabled, GracePeriod, Timeout, Force, IgnoreDaemonSets, DeleteEmptyDirData, PodSelector, SkipWaitForDeleteTimeout
- `ConcurrencyOption`: Limit(30), WorkerDisruptionPercent(10), Uploads(5)
- `EvictTaintOption`: Enabled, Taint, Effect, ControllerWorkers
- `DrainOption.ToKubectlArgs()` — build kubectl drain flags

### cluster/hook.go
- `Hooks` type: `map[string]map[string][]string`
- `ForActionAndStage(action, stage) []string`

### cluster/uploadfile.go
- `UploadFile` struct: Name, Source, Data, Sha256, DestinationDir, DestinationFile, PermMode, DirPermMode, User, Group
- Supports glob patterns, URLs, inline data
- `ResolveRelativeTo(baseDir)` — resolve local paths

### cluster/uploadfile_test.go
- Upload file resolution tests

### cluster/binary_providers.go
- `Host.defaultBinaryProvider()` — selects provider based on host config:
  - UseExistingK3s → binprovider.NewExisting
  - K3sBinaryPath → binprovider.NewLocalFile
  - UploadBinary → binprovider.NewLocalUpload
  - K3sDownloadURLOverride → binprovider.NewCustomURL
  - Default → binprovider.NewGitHub (k3s-io/k3s releases)

### cluster/quote.go
- `quoter` interface and `quote()` helper for shell-safe quoting

## Implementation Order
1. `cluster/quote.go` (utility)
2. `cluster/flags.go` (utility)
3. `cluster/hook.go` (utility)
4. `cluster/uploadfile.go` (utility)
5. `cluster/options.go` (config options)
6. `cluster/k3s.go` (K3s config type)
7. `cluster/host.go` (core host type)
8. `cluster/hosts.go` (host collection)
9. `cluster/spec.go` (cluster spec)
10. `cluster.go` (top-level cluster type)
11. `cluster/binary_providers.go` (binary provider selection)
12. All test files
