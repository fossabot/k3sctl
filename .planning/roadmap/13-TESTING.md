# 13: Testing Strategy

## Overview

k3sctl uses hermetic unit tests with mocks/fakes (no real hosts), behavior-named tests, and smoke tests for integration validation.

## Unit Test Files

### Phase Tests

| Test File | Notes |
|---|---|
| `phase/manager_test.go` | Phase manager lifecycle |
| `phase/connect_test.go` | Connection handling |
| `phase/validate_hosts_test.go` | Role names (server/agent) |
| `phase/validate_facts_test.go` | K3s version format |
| `phase/gather_facts_test.go` | Host fact fields |
| `phase/gather_k3s_facts_test.go` | K3s version parsing, service detection |
| `phase/configure_k3s_test.go` | K3s YAML config format |
| `phase/get_kubeconfig_test.go` | File read for kubeconfig |
| `phase/stage_binaries_test.go` | K3s binary URLs |
| `phase/uploadfiles_test.go` | File upload |

### API Type Tests

| Test File | Notes |
|---|---|
| `pkg/apis/.../cluster_test.go` | API version, storage type |
| `pkg/apis/.../cluster/k3s_test.go` | Token format, version parsing |
| `pkg/apis/.../cluster/host_test.go` | Role names, service names, K3s commands |
| `pkg/apis/.../cluster/hosts_test.go` | Role-based methods |
| `pkg/apis/.../cluster/spec_test.go` | K3sLeader, K3s fields |
| `pkg/apis/.../cluster/flags_test.go` | Flag manipulation |
| `pkg/apis/.../cluster/flags_quoting_test.go` | Flag quoting |
| `pkg/apis/.../cluster/uploadfile_test.go` | File upload types |

### Configurer Tests

| Test File | Notes |
|---|---|
| `configurer/linux/linux_test.go` | K3s paths, service names |
| `configurer/linux/alpine_test.go` | Alpine-specific |

### CLI Tests

| Test File | Notes |
|---|---|
| `cmd/apply_test.go` | k3sctl output format |
| `cmd/flags_test.go` | k3sctl.yaml config file |
| `cmd/init_test.go` | K3s default config format |

### Package Tests

| Test File | Notes |
|---|---|
| `pkg/retry/retry_test.go` | Retry logic |
| `pkg/download/download_test.go` | Download utility |
| `pkg/node/statusfunc_test.go` | Node status checks |
| `pkg/manifest/reader_test.go` | Manifest reader |

## K3s-Specific Tests

- `pkg/k3s/version_test.go` — K3s version parsing (`v1.30.0+k3s1` format)
- `pkg/k3s/binprovider/download_test.go` — K3s GitHub URL generation
- `phase/default_k3s_version_test.go` — K3s channel API parsing
- `configurer/linux/k3s_paths_test.go` — K3s-specific paths across distros

## Testing Conventions

1. **Naming**: `TestMethodName_scenario` (e.g., `TestValidateHosts_DuplicateAddresses`)
2. **Assertions**: `github.com/stretchr/testify` (assert and require)
3. **Deep equality**: `github.com/go-test/deep` where needed
4. **Mocking**: Interface-based mocks via counterfeiter, no external mock frameworks
5. **Hermetic**: No real SSH connections, no network calls, no host-specific config
6. **Build tags**: No `_test` build tags needed (tests in same package)

## Smoke Tests

### Phase 1: Unit Tests Only
- Comprehensive unit test coverage
- Mock SSH/rig interactions via interfaces

### Phase 2: Smoke Tests (Future)
- Container-based fake VMs (bootloose)
- Key smoke tests:
  - `smoke-basic` — basic cluster create/destroy
  - `smoke-upgrade` — version upgrade flow
  - `smoke-reset` — cluster reset
  - `smoke-dryrun` — dry-run mode verification
  - `smoke-files` — file upload functionality
  - `smoke-init` — init command output

### Smoke Test Infrastructure
- Uses bootloose (containerized SSH-accessible "VMs")
- Each smoke test: create cluster → verify → cleanup
- Matrix: multiple Linux distros (Alpine, Ubuntu, CentOS, etc.)
