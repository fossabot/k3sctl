# 08: Upgrade and Reset Phases

## Phase: UpgradeServers (`phase/upgrade_servers.go`)
- Detect servers needing upgrade: compare K3s versions
- Stop K3s service: `systemctl stop k3s`
- Replace binary: same install binary flow
- Restart service: `systemctl start k3s`
- Wait for readiness: `/readyz` endpoint
- Drain if needed: `k3s kubectl drain`
- Uncordon after: `k3s kubectl uncordon`
- Sequential upgrade for servers (one at a time for HA)
- Evict taint support
- Reinstall service if flags changed
- Handle `--cluster-init` → `--server` transition

## Phase: UpgradeAgents (`phase/upgrade_agents.go`)
- Detect agents needing upgrade
- Drain node (with disruption budget): `k3s kubectl drain` from leader
- Stop agent service: `systemctl stop k3s-agent`
- Replace binary
- Restart agent: `systemctl start k3s-agent`
- Wait for node ready
- Uncordon: `k3s kubectl uncordon`
- Parallel with disruption percentage limits
- NoDrain option support
- Evict taint support

## Phase: Reinstall (`phase/reinstall.go`)
- Detect hosts where install flags have changed
- Stop service, reinstall with new flags, restart
- Update systemd unit file with new arguments
- Same drain/uncordon pattern

## Phase: ResetAgents (`phase/reset_agents.go`)
- Stop agent service: `systemctl stop k3s-agent`
- Delete node from cluster: `k3s kubectl delete node` from leader
- Drain before delete (optional)
- Run K3s agent uninstall: `/usr/local/bin/k3s-agent-uninstall.sh` (if exists) or manual cleanup
- Manual cleanup: remove data dir, config, systemd unit
- NoDrain, NoDelete options

## Phase: ResetServers (`phase/reset_servers.go`)
- Stop server service: `systemctl stop k3s`
- Leave etcd cluster if embedded etcd
- Delete node from cluster
- Run K3s uninstall: `/usr/local/bin/k3s-uninstall.sh` (if exists) or manual cleanup
- Manual cleanup: remove data dir, config, token files, systemd units
- NoLeave, NoDrain, NoDelete options

## Phase: ResetLeader (`phase/reset_leader.go`)
- Special handling for the last server (leader reset)
- Must be last to leave the cluster
- Run uninstall script or manual cleanup
- Remove all K3s state: data dir, config, binaries

## K3s Reset Details

1. K3s provides uninstall scripts: `k3s-uninstall.sh` and `k3s-agent-uninstall.sh`
2. These scripts handle: stopping service, removing systemd units, removing data, removing binary
3. k3sctl prefers these scripts when available, with manual fallback

## Upgrade Strategy

K3s upgrade flow:
1. Drain node (if agent)
2. Stop K3s service
3. Replace `/usr/local/bin/k3s` binary with new version
4. Start K3s service
5. Wait for node ready
6. Uncordon (if drained)

Server upgrade must be sequential in HA clusters. Agent upgrade can be batched with disruption budget.
