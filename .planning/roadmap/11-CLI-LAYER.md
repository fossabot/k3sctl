# 11: CLI Layer — Commands and Flags

## `cmd/root.go`
- `NewK3sctl(in, out, errOut) *cli.App`
- App name: `k3sctl`
- Usage: `k3s cluster management tool`
- Commands: version, apply, kubeconfig, init, reset, config(edit/status), completion
- Global flags: debug, trace, no-redact
- Signal trapping: SIGINT/SIGTERM graceful shutdown
- Context management with cancel

## `cmd/flags.go`
- `configFlag` default: `k3sctl.yaml`
- `initConfig()` — multi-doc YAML reading with glob support
- `initManager()` — creates phase.Manager with config
- `initLogging()` — logrus setup with file logging
- `displayLogo()` — k3sctl logo
- `displayCopyright()` — copyright info
- `warnOldCache()` — XDG cache path for k3sctl
- `warnRigMigration()` — rig v2 migration warning
- Flag definitions: concurrency, timeout, dry-run, etc.
- Config loading: parse `k3sctl.k3sproject.io/v1beta1` API version

## `cmd/apply.go`
- Before: actions(initLogging, initConfig, initManager, displayLogo, displayCopyright, warnOldCache, warnRigMigration)
- Flags: config, concurrency, dry-run, no-wait, no-drain, restore-from, kubeconfig-out, etc.
- `evict-taint` flag support
- Creates `action.NewApply(opts)` and runs it

## `cmd/reset.go`
- action.Reset{}

## `cmd/kubeconfig.go`
- action.Kubeconfig{}

## `cmd/init.go`
- Default K3s YAML config:
  ```yaml
  write-kubeconfig-mode: "0644"
  tls-san: []
  ```
- API version: `k3sctl.k3sproject.io/v1beta1`
- Default roles: server, agent
- Cluster name default: `k3s-cluster`
- Address parsing, host generation

### `cmd/init_test.go`

## `cmd/config_edit.go`
- action.ConfigEdit{}

## `cmd/config_status.go`
- action.ConfigStatus{}

## `cmd/version.go`
- Display k3sctl version info

## `cmd/completion.go`
- Shell completion for k3sctl

## `cmd/logo.go`
- k3sctl ASCII art logo

## Test Files

### `cmd/apply_test.go`
### `cmd/flags_test.go`
