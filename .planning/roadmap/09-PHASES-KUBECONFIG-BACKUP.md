# 09: Kubeconfig, Backup, and Remaining Phases

## Phase: GetKubeconfig (`phase/get_kubeconfig.go`)
- Reads kubeconfig file directly from `/etc/rancher/k3s/k3s.yaml` on server node
- Kubeconfig rewriting logic (cluster name, user, API address)
- Uses `k8s.io/client-go/tools/clientcmd`
- `readKubeconfig` function:
  ```go
  var readKubeconfig = func(h *cluster.Host) (string, error) {
      output, err := h.Sudo().FS().ReadFile("/etc/rancher/k3s/k3s.yaml")
      return string(output), err
  }
  ```
- DryRun handling

### `phase/get_kubeconfig_test.go`

## Phase: Backup (`phase/backup.go`)
- K3s supports etcd snapshots: `k3s etcd-snapshot save --name <name> --data-dir <dir>`
- Options:
  1. **SQLite backup**: if using default SQLite, backup the DB file
  2. **Etcd snapshot**: if using embedded etcd (`--cluster-init`), use `k3s etcd-snapshot save`
  3. **Data-dir backup**: tar up the data directory
- Download snapshot file from server

## Phase: Restore (`phase/restore.go`)
- K3s can restore from etcd snapshot: `k3s server --cluster-reset --cluster-reset-restore-path=<snapshot>`
- Only works with embedded etcd clusters
- Phase sequence: Connect → DetectOS → GatherFacts → GatherK3sFacts → Backup → Disconnect

## Remaining Minor Phases

### `phase/apply_manifests.go`
- Uses `k3s kubectl apply` to apply additional manifests after cluster ready

### `phase/daemon_reload.go`
- Runs `systemctl daemon-reload`

### `phase/downloads.go`
- K3s download URLs

### `phase/lock.go` / `phase/unlock.go`
- Lock file path from configurer (`K3sctlLockFilePath`)
