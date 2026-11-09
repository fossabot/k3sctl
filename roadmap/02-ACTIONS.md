# Roadmap: Action Layer

## Source Files

### action/apply.go
- `ApplyOptions` struct: Manager, DisableDowngradeCheck, NoDrain, NoWait, RestoreFrom, KubeconfigOut, KubeconfigAPIAddress, KubeconfigUser, KubeconfigCluster, EvictTaint
- `NewApply(opts)` creates phase sequence:
  1. DefaultK3sVersion
  2. Connect
  3. DetectOS
  4. PrepareHosts
  5. GatherFacts
  6. ValidateHosts
  7. GatherK3sFacts
  8. ValidateFacts
  9. StageBinaries
  10. UploadFiles
  11. InstallBinaries
  12. ConfigureK3s
  13. Restore
  14. RunHooks (before, apply)
  15. InitializeK3s
  16. InstallServers
  17. InstallAgents
  18. UpgradeServers
  19. UpgradeAgents
  20. Reinstall
  21. ResetAgents
  22. ResetServers
  23. RunHooks (after, apply)
  24. ApplyManifests
  25. [GetKubeconfig] (optional)
  26. Disconnect
- `Run()`: executes phases, prints duration, suggests kubeconfig command

### action/reset.go
- Confirms with user (survey prompt or --force)
- Marks all hosts Reset=true
- Phase sequence: DefaultK3sVersion→Connect→DetectOS→Lock→PrepareHosts→GatherFacts→GatherK3sFacts→ResetAgents→ResetServers→ResetLeader→DaemonReload→Unlock→Disconnect
- Reset stops the K3s service, removes data dirs, and removes the binary

### action/backup.go
- Phase sequence: DefaultK3sVersion→Connect→DetectOS→Lock→PrepareHosts→GatherFacts→GatherK3sFacts→Backup→Unlock→Disconnect
- Uses `k3s etcd-snapshot save` on server node

### action/kubeconfig.go
- Filters hosts to only leader server
- Phase sequence: Connect→DetectOS→GetKubeconfig→Disconnect
- Reads `/etc/rancher/k3s/k3s.yaml` from first server

### action/config_edit.go
- Edits `/etc/rancher/k3s/config.yaml` on remote host

### action/config_status.go
- Shows current config from remote host

## Implementation Order
1. `action/apply.go` (core workflow)
2. `action/reset.go`
3. `action/kubeconfig.go`
4. `action/backup.go`
5. `action/config_edit.go` + `action/config_status.go`
