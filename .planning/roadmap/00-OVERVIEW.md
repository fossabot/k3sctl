# k3sctl Roadmap: Complete Project Overview

## Mission

`k3sctl` is a pure Go cluster management tool for K3s. It uses a phase-based execution model to manage K3s cluster lifecycles over SSH.

## Repository

- **Repo**: `github.com/nirantarasi/k3sctl`
- **Go module**: `github.com/nirantarasi/k3sctl`

## K3s Reference

| Aspect | K3s (k3sctl) |
|--------|-------------|
| Binary path | `/usr/local/bin/k3s` |
| Config path | `/etc/rancher/k3s/config.yaml` |
| Data dir | `/var/lib/rancher/k3s` |
| Token path | `/etc/rancher/k3s/token` (server token file) |
| Kubeconfig | `/etc/rancher/k3s/k3s.yaml` (file on server node) |
| Service names | `k3s` (server), `k3s-agent` (agent) |
| Roles | server, agent |
| Join mechanism | Token-based (server token at `/var/lib/rancher/k3s/server/node-token`) |
| Default storage | Embedded SQLite (or external etcd/MySQL/PostgreSQL) |
| Windows support | No |
| Install command | `k3s server` / `k3s agent` (systemd units) |
| kubectl | `k3s kubectl` |
| Version format | `v1.30.0+k3s1` |
| API port default | 6443 |
| Backup | `k3s etcd-snapshot save` (embedded etcd only) |
| Config mechanism | Config file (`config.yaml`) or CLI flags |
| Kubelet extra args | Via `--kubelet-arg` flag |
| CNI | flannel (default) |

## Architecture

```
┌─────────────────────────────────────────────┐
│                   CLI Layer                  │
│               (cmd/ - urfave/cli)            │
│  apply | reset | init | kubeconfig | config  │
└───────────────────────┬─────────────────────┘
                        │
┌───────────────────────▼─────────────────────┐
│                Action Layer                  │
│              (action/ - workflows)           │
│   Compose phases in correct order            │
└───────────────────────┬─────────────────────┘
                        │
┌───────────────────────▼─────────────────────┐
│          Phase Manager + Phases              │
│              (phase/ - execution)            │
│   Prepare → ShouldRun → Before → Run →      │
│   After → CleanUp (lifecycle)                │
└───────────────────────┬─────────────────────┘
                        │
┌───────────────────────▼─────────────────────┐
│         Configurer Layer                     │
│     (configurer/ - OS abstraction)           │
│   Linux base + distro-specific modules       │
└───────────────────────┬─────────────────────┘
                        │
┌───────────────────────▼─────────────────────┐
│            Remote Execution                  │
│        (rig v2 - SSH/localhost)              │
│   Connect, exec, file transfer, sudo         │
└─────────────────────────────────────────────┘
```

## Implementation Phases (Development Order)

1. **Foundation** (Documents 01-03): Go module, core types, configurer interface
2. **Execution Engine** (Documents 04-05): Phase manager, generic phase, retry
3. **API Types** (Document 06): Cluster, Spec, Host, K3s config, Options
4. **Phases** (Documents 07-09): All ~35 phases for K3s
5. **Actions** (Document 10): Apply, Reset, Kubeconfig, Init workflows
6. **CLI** (Document 11): Commands, flags, config loading
7. **Supporting Packages** (Document 12): Binary providers, download, node status
8. **Testing** (Document 13): Unit tests for all packages
9. **Build & CI** (Document 14): Taskfile, goreleaser, GitHub Actions

## Files Not Needed

- `configurer/windows.go` — K3s has no Windows support
- `configurer/windows/windows.go` — K3s has no Windows support
- `configurer/windows/windows_test.go` — K3s has no Windows support
- `configurer/windows_test.go` — K3s has no Windows support
- All Windows-specific code paths in `host.go`, `install_workers.go`, etc.

## Key Design Decisions

- Roles: server/agent (not controller/worker)
- Service names: `k3s`/`k3s-agent`
- Token: read from `/var/lib/rancher/k3s/server/node-token` on leader
- Kubeconfig: read file directly from `/etc/rancher/k3s/k3s.yaml`
- Config: flat YAML at `/etc/rancher/k3s/config.yaml`
- Install: place binary + create systemd unit (no install script)
- Backup: `k3s etcd-snapshot save` for embedded etcd clusters
- Reset: uninstall scripts (`k3s-uninstall.sh`, `k3s-agent-uninstall.sh`) with manual fallback
