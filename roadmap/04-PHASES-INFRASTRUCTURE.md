# Roadmap: Infrastructure Phases (Reusable)

These phases are infrastructure and require minimal K3s-specific changes.

## phase/connect.go
- Connects to all hosts in parallel with retry (10min timeout)
- Handles auth rejection limit (12 attempts before abort)
- Uses `rig.NewClient` with retry disabled (phases handle retries)
- Configures reconnect, keep-alive

### phase/connect_test.go
- Tests for connection setup

## phase/detect_os.go
- Resolves OS configurer per host via `configurer.ResolveOSModule`
- Falls back to `ID_LIKE` entries from `/etc/os-release`
- Sets `K3sBinaryPath` via configurer

## phase/disconnect.go
- Disconnects from all hosts

## phase/lock.go
- Creates lock file on each host to prevent concurrent k3sctl operations
- Uses `/var/lock/k3sctl.lock` or configurer-provided path

## phase/unlock.go
- Removes lock file from each host

## phase/prepare_hosts.go
- Validates host connectivity, sets up host metadata
- Prepares download tracker per host

## phase/gather_facts.go
- Detects per-host: architecture, machine ID, hostname, private interface/address
- Handles hostname override
- Checks for existing K3s binary

### phase/gather_facts_test.go
- Tests for fact gathering logic

## phase/validate_hosts.go
- Validates gathered host facts (arch support, connectivity)

### phase/validate_hosts_test.go

## phase/validate_facts.go
- Validates cluster-level facts after gathering
- Checks for version compatibility, downgrade prevention

### phase/validate_facts_test.go

## phase/run_hooks.go
- Executes user-defined hooks (before/after actions)
- Hooks defined in host config as command strings
- Hooks run shell commands on remote hosts via rig exec (user-defined, not application logic)

## phase/daemon_reload.go
- Runs `systemctl daemon-reload` on each host

## phase/downloads.go
- Downloads extra files defined in host upload config (URL sources)
- Uses download.Tracker for caching/deduplication

## phase/uploadfiles.go
- Uploads local files to remote hosts
- Handles glob patterns, permissions, ownership

### phase/uploadfiles_test.go

## Implementation Order
1. connect.go + disconnect.go (core connectivity)
2. detect_os.go (OS detection)
3. lock.go + unlock.go
4. prepare_hosts.go
5. gather_facts.go
6. validate_hosts.go + validate_facts.go
7. run_hooks.go + daemon_reload.go
8. downloads.go + uploadfiles.go
