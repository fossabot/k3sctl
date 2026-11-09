# 07: Install and Configuration Phases

## Phase: DefaultK3sVersion (`phase/default_k3s_version.go`)
- Fetch latest K3s version from channels: `https://update.k3s.io/v1-release/channels`
- K3s version channel: "stable", "latest", "v1.30", etc.
- Set `Config.Spec.K3s.Version` if not specified
- Channel API returns JSON with `data[].name` and `data[].latest`

## Phase: StageBinaries (`phase/stage_binaries.go`)
- Uses `h.K3sBinaryProvider(target)`
- Binary caching pattern: `BinaryCacher` interface
- Downloads K3s binary from GitHub releases: `github.com/k3s-io/k3s/releases/download/v1.30.0+k3s1/k3s` (or `k3s-arm64`, `k3s-armhf`)
- No `.exe` variant (no Windows)
- Upload vs download distinction: same pattern

### `phase/stage_binaries_test.go`

## Phase: UploadFiles (`phase/uploadfiles.go`)
- Generic file upload to hosts
- Uses `parallelDoUpload` for concurrency control

### `phase/uploadfiles_test.go`

## Phase: InstallBinaries (`phase/install_binaries.go`)
- Calls `h.InstallK3sBinary(path)`
- Verifies installed version: `k3s --version`
- Sets correct permissions (0750)
- Creates binary directory if needed

## Phase: PrepareArm (`phase/arm_prepare.go`)
- ARM preparation
- Sets up ARM-specific requirements (e.g. `ETCD_UNSUPPORTED_ARCH`)

## Phase: ConfigureK3s (`phase/configure_k3s.go`)
- K3s config is a YAML file at `/etc/rancher/k3s/config.yaml`
- K3s config format is flat key-value:
  ```yaml
  write-kubeconfig-mode: "0644"
  tls-san:
    - "my-lb.example.com"
  cluster-init: true
  ```
- Write config file to each server node
- Agent nodes may also need config (for `--server` URL, token, etc.)
- Config diff detection for upgrades

### `phase/configure_k3s_test.go`

## Phase: EnsureJoinTokenWorkaround (`phase/ensure_join_token_workaround.go`)
- Handles token edge cases for specific K3s versions
- May be a no-op for most versions

### `phase/ensure_join_token_workaround_test.go`

## Phase: InitializeK3s (`phase/initialize_k3s.go`)
- Start the first K3s server node (the leader)
- Create systemd service, start it
- For HA: use `--cluster-init` flag on first server (enables embedded etcd)
- Wait for K3s to be ready: check `/readyz` endpoint
- Read and store the join token from `/var/lib/rancher/k3s/server/node-token`

## Phase: InstallServers (`phase/install_servers.go`)
- Join token: read from leader's `/var/lib/rancher/k3s/server/node-token`
- Place binary, create systemd unit, configure, start
- Server flags: `--server https://<leader>:6443 --token <token> --data-dir <dir>`
- Quorum-aware parallel batching: odd number batches
- Service name: `k3s`
- Readiness check: `k3s kubectl get --raw='/readyz'`

## Phase: InstallAgents (`phase/install_agents.go`)
- Join token: same token from leader, for agent use
- Place binary, create systemd unit with `k3s agent` command
- Agent flags: `--server https://<server>:6443 --token <token> --data-dir <dir>`
- Service name: `k3s-agent`
- API validation: check connectivity to server API
- Node ready wait: `k3s kubectl get node` from leader
- No Windows path

## Phase: RunHooks (`phase/run_hooks.go`)
- Generic hook execution
- Runs before/after hooks for apply, install, upgrade, reset actions

## Phase: ApplyManifests (`phase/apply_manifests.go`)
- Uses `k3s kubectl apply` to apply additional manifests after cluster ready

## Phase: DaemonReload (`phase/daemon_reload.go`)
- Runs `systemctl daemon-reload`

## Phase: Downloads (`phase/downloads.go`)
- K3s download URL pattern
