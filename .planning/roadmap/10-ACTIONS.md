# 10: Action Layer — Workflow Composition

## Overview

Actions compose phases into workflows. Each action creates an ordered list of phases and runs them through the Manager.

## `action/apply.go`
- `ApplyOptions` struct
- `NewApply()` phase list:

```go
Phases: phase.Phases{
    &phase.DefaultK3sVersion{},
    &phase.Connect{},
    &phase.DetectOS{},
    &phase.PrepareHosts{},
    &phase.GatherFacts{},
    &phase.ValidateHosts{},
    &phase.GatherK3sFacts{},
    &phase.ValidateFacts{SkipDowngradeCheck: opts.DisableDowngradeCheck},
    &phase.StageBinaries{},
    &phase.UploadFiles{},
    &phase.InstallBinaries{},
    &phase.PrepareArm{},
    &phase.ConfigureK3s{},
    &phase.Restore{RestoreFrom: opts.RestoreFrom},
    &phase.RunHooks{Stage: "before", Action: "apply"},
    &phase.InitializeK3s{},
    &phase.InstallServers{},
    &phase.InstallAgents{},
    &phase.UpgradeServers{},
    &phase.UpgradeAgents{NoDrain: opts.NoDrain},
    &phase.Reinstall{},
    &phase.ResetAgents{NoDrain: opts.NoDrain},
    &phase.ResetServers{NoDrain: opts.NoDrain},
    &phase.RunHooks{Stage: "after", Action: "apply"},
    &phase.ApplyManifests{},
}
```

- Success message: "k3s cluster version %s is now installed"
- Kubeconfig tip: "k3sctl kubeconfig"
- ConfigPaths default: "k3sctl.yaml"

## `action/reset.go`
- Confirmation message: "Going to reset all of the hosts..."
- Phase list:
  ```go
  &phase.DefaultK3sVersion{},
  &phase.Connect{},
  &phase.DetectOS{},
  lockPhase,
  &phase.PrepareHosts{},
  &phase.GatherFacts{SkipMachineIDs: true},
  &phase.GatherK3sFacts{},
  &phase.ResetAgents{NoDrain: true, NoDelete: true},
  &phase.ResetServers{NoDrain: true, NoDelete: true, NoLeave: true},
  &phase.ResetLeader{},
  &phase.DaemonReload{},
  &phase.Unlock{Cancel: lockPhase.Cancel},
  &phase.Disconnect{},
  ```

## `action/kubeconfig.go`
- `K3sLeader()` for leader identification
- Phases: Connect → DetectOS → GetKubeconfig → Disconnect
- Only connects to leader server node

## `action/config_edit.go`
- K3s config is a file on disk, not a Kubernetes resource
- SSH to leader, edit `/etc/rancher/k3s/config.yaml` directly
- Open in `$EDITOR`, diff, write back, restart K3s service

## `action/config_status.go`
- Shows: K3s version, cluster status, node status
- Reads `/etc/rancher/k3s/config.yaml` and displays status
- Can also run `k3s kubectl get nodes -o wide`
