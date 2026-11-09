# Roadmap: Entry Point & CLI Layer

## Source Files

### main.go
- Entry point: creates `cli.App` via `cmd.NewK3sctl()`, runs with `os.Args`
- Blank imports for `versioninfo` and `version` packages for side-effect registration

### version/version.go
- Build-time injected vars: `Version`, `K3sVersion`, `HelperVersionOverride`
- LDFLAGS injection pattern

### cmd/root.go
- `NewK3sctl(in, out, errOut)` → returns `*cli.App`
- Commands: version, apply, kubeconfig, init, reset, backup, config (edit/status), completion
- Signal trapping: SIGINT/SIGTERM with graceful shutdown via context cancellation

### cmd/flags.go
- Defines CLI flags: debug, trace, redact, dryRun, force, config, concurrency, concurrentUploads, timeout, retryTimeout, retryInterval
- `initConfig`: reads YAML config, envsubst, glob expansion, manifest parsing
- `initManager`: creates `phase.Manager` from parsed config
- `initLogging`: logrus with screen+file hooks
- `readConfig`: unmarshals Cluster, validates
- Config reader supports stdin (`-`), glob patterns, env variable substitution
- K3s uses file-based config — no ClusterConfig merge needed

### cmd/apply.go
- Apply command with flags: config, concurrency, no-wait, no-drain, restore-from, kubeconfig-out, kubeconfig-api-address, kubeconfig-user, kubeconfig-cluster, disable-downgrade-check, evict-taint, force
- Before chain: initLogging→initConfig→initManager→displayLogo→displayCopyright→warnOldCache→warnRigMigration
- Creates `action.NewApply(opts)` and runs it

### cmd/init.go
- Generates config template YAML (k3sctl.yaml)
- Flags: cluster-name, controller-count, user, key-path
- `buildHosts()` creates host entries from addresses
- `hostFromAddress()` parses `user@address:port`
- Generates template with K3s defaults (server/agent roles, rancher paths)

### cmd/reset.go
- Reset command: uses action.NewReset
- Flags: config, concurrency, force

### cmd/backup.go
- Backup command: uses action.NewBackup
- Uses `k3s etcd-snapshot save` on server node

### cmd/kubeconfig.go
- Kubeconfig retrieval command
- Reads `/etc/rancher/k3s/k3s.yaml` from first server

### cmd/config_edit.go
- File-based config editing on remote host
- Edits `/etc/rancher/k3s/config.yaml`

### cmd/config_status.go
- Config status/validation from remote host

### cmd/version.go
- Displays version info, checks for updates via GitHub API

### cmd/completion.go
- Shell completion generation (bash/zsh/fish/powershell)

### cmd/logo.go
- ASCII logo display

### cmd/flags_test.go, cmd/init_test.go, cmd/apply_test.go
- Unit tests for flag parsing, init generation, apply config

## Implementation Order
1. `main.go` + `version/version.go`
2. `cmd/root.go` (skeleton with all subcommands)
3. `cmd/flags.go` (config loading pipeline)
4. `cmd/init.go` (config template generation)
5. `cmd/apply.go`
6. `cmd/reset.go`
7. `cmd/backup.go`
8. `cmd/kubeconfig.go`
9. `cmd/config_edit.go` + `cmd/config_status.go`
10. `cmd/version.go` + `cmd/completion.go` + `cmd/logo.go`
