# 04: Retry Package and Utility Packages

## `pkg/retry/retry.go`
- Generic retry utility
- **Exported**:
  - `DefaultTimeout = 2 * time.Minute`
  - `Interval = 5 * time.Second`
  - `ErrAbort` — sentinel for non-retryable errors
  - `Context(ctx, f)` — retry at Interval until success or context cancel
  - `Timeout(ctx, timeout, f)` — Context with timeout
  - `WithDefaultTimeout(ctx, f)` — Timeout with DefaultTimeout
  - `Times(ctx, times, f)` — retry up to N times

### `pkg/retry/retry_test.go`

## `pkg/node/statusfunc.go`
- Node status check functions:
  - `ServiceRunningFunc(h, serviceName)` — checks if systemd service is running
  - `HTTPStatusFunc(h, url, expected...)` — checks HTTP status codes
  - `KubeNodeReadyFunc(h)` — checks if node is Ready via kubectl
- Service names: `k3s` (server) and `k3s-agent` (agent)
- kubectl: uses `k3s kubectl`

### `pkg/node/statusfunc_test.go`

## `pkg/download/download.go`
- Generic download utility
- Provides HTTP download with progress, hash verification

### `pkg/download/download_test.go`

## `pkg/manifest/reader.go`
- Reads multi-document YAML, validates apiVersion/kind
- API version: `k3sctl.k3sproject.io/v1beta1`

### `pkg/manifest/reader_test.go`

## `integration/github/github.go`
- K3s GitHub releases integration
  - K3s releases: `github.com/k3s-io/k3s/releases`
  - Version channel: `update.k3s.io/v1-release/channels`
