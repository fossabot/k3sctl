# Roadmap: K3s-Specific Phases

These phases contain K3s-specific logic for cluster lifecycle management.

## phase/default_k3s_version.go
- Resolves default K3s version via GitHub releases API (`k3s-io/k3s`)
- Sets `Spec.K3s.Version` if not explicitly set
- K3s release naming: `v1.28.2+k3s1` format

## phase/gather_k3s_facts.go
- On each host, runs `k3s --version` to get installed binary version (parses "k3s version v1.x.x+k3s1")
- Checks if k3s/k3s-agent service is running via systemd
- Reads existing config from `/etc/rancher/k3s/config.yaml`
- Checks for existing node-token at `/var/lib/rancher/k3s/server/node-token`

### phase/gather_k3s_facts_test.go
- Tests for K3s version parsing

## phase/stage_binaries.go
- Stages K3s binary using BinaryProvider pattern
- Downloads/uploads binary to temp path on host
- Handles BinaryCacher interface for local upload dedup
- K3s binary URL format: `https://github.com/k3s-io/k3s/releases/download/{version}/k3s` (amd64)
- ARM64: `k3s-arm64`, ARMhf: `k3s-armhf`

### phase/stage_binaries_test.go

## phase/install_binaries.go
- Moves staged binary to install location (`/usr/local/bin/k3s`), sets permissions

## phase/configure_k3s.go
- Generates `/etc/rancher/k3s/config.yaml` (YAML format)
- Creates node-specific config from cluster config
- Server config: cluster-init, tls-san, node-ip, flannel settings, etc.
- Agent config: server URL, token
- File-based configuration (no dynamic config CRD)

### phase/configure_k3s_test.go

## phase/initialize_k3s.go
- First server: `k3s server --cluster-init` (for HA) or plain `k3s server`
- Installs as systemd service
- Waits for node-token to be generated
- Waits for K3s API to be ready

## phase/install_servers.go
- Reads node-token from first server
- Installs K3s server with `--server https://<first-server>:6443 --token <token>`
- Quorum-aware batch installation for multiple servers
- Installs as systemd service
- Waits for API readiness
- CleanUp resets failed installs; After hook invalidates tokens

## phase/install_agents.go
- Reads node-token from server
- Installs K3s agent with `K3S_URL=https://<server>:6443 K3S_TOKEN=<token>`
- Or uses config.yaml with server/token fields
- Installs as systemd service
- Waits for node ready status

## phase/upgrade_servers.go
- Rolling upgrade (one at a time for servers)
- Stops K3s service, replaces binary, restarts, waits for ready

## phase/upgrade_agents.go
- Drain, stop, replace binary, reinstall, start, uncordon
- Batched parallel upgrade respecting WorkerDisruptionPercent

## phase/reinstall.go
- Reinstalls K3s service when install flags have changed
- Detects flag changes by comparing current vs desired

## phase/reset_agents.go
- Stops k3s-agent, removes data dirs, removes binary, removes systemd unit

## phase/reset_servers.go
- Stops K3s, removes data dirs (`/var/lib/rancher/k3s`), removes config, removes binary

## phase/reset_leader.go
- Handles leader reset specially (last server to reset — reset initial server last)

## phase/get_kubeconfig.go
- Reads `/etc/rancher/k3s/k3s.yaml`, rewrites server address/user/cluster
- Uses k8s client-go clientcmd for kubeconfig manipulation

### phase/get_kubeconfig_test.go

## phase/backup.go
- Runs `k3s etcd-snapshot save` or copies SQLite database
- For etcd mode: `k3s etcd-snapshot save --dir /path`
- For SQLite mode: copy `/var/lib/rancher/k3s/server/db/`
- Downloads backup archive from host

## phase/restore.go
- Uploads snapshot, runs `k3s server --cluster-reset --cluster-reset-restore-path=<snapshot>`

## phase/apply_manifests.go
- Applies extra Kubernetes manifests to cluster
- Writes manifests to leader, runs `k3s kubectl apply`

## phase/validate_etcd_members.go
- Validates etcd member list matches expected servers
- Only applicable when K3s uses etcd (`--cluster-init` or `--datastore-endpoint`)
- For SQLite mode: skip entirely

## phase/arm_prepare.go
- ARM-specific preparation (kernel modules, settings)
- K3s supports ARM natively, but may still need cgroup setup
- Check cgroup v2 configuration on ARM hosts

## Implementation Order
1. default_k3s_version.go
2. gather_k3s_facts.go
3. stage_binaries.go + install_binaries.go
4. configure_k3s.go
5. initialize_k3s.go
6. install_servers.go
7. install_agents.go
8. get_kubeconfig.go
9. upgrade_servers.go + upgrade_agents.go
10. reinstall.go
11. reset_agents.go + reset_servers.go + reset_leader.go
12. backup.go + restore.go
13. apply_manifests.go
14. validate_etcd_members.go (conditional)
