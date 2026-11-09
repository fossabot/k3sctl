# Roadmap: Supporting Packages

## pkg/retry/retry.go
- Context-based retry helpers
- `Context(ctx, f)` — retry at `Interval` (5s) until context done
- `Timeout(ctx, timeout, f)` — retry with deadline
- `WithDefaultTimeout(ctx, f)` — retry with 2 minute default
- `Times(ctx, n, f)` — retry N times
- `ErrAbort` — sentinel to stop retrying
- `DefaultTimeout` = 2min, `Interval` = 5s

### pkg/retry/retry_test.go
- Retry logic tests

## pkg/download/download.go
- URL download tracker with caching/deduplication
- `Tracker` struct: manages downloads per host
- `Record` struct: tracks downloaded files (URLDigest, Destination, Size, ModTime, ETag, LastModified, Sha256)
- Key methods:
  - `NewTracker(files, user, label)` — create tracker per host
  - `Needed(ctx, url, dst, wantSum)` — check if download needed (returns Verdict)
  - `Remember(ctx, url, dst, sum)` — record an existing file
  - `Fetch(ctx, url, dst, wantSum)` — download with tracking
- Features: resume support, checksum verification, HTTP HEAD comparison, ETag/Last-Modified tracking
- Records stored in user cache dir (`~/.cache/k3sctl/downloads/`)

### pkg/download/download_test.go
- Download tracker tests

## pkg/k3s/provider.go
- `BinaryProvider` interface:
  - `NeedsUpgrade() bool`
  - `IsUpload() bool`
  - `Stage(ctx) (string, error)`
  - `CleanUp(ctx)`
- `BinaryCacher` interface (optional):
  - `BinaryCacheKey() (string, error)`
  - `EnsureCached(ctx) error`

## pkg/k3s/binprovider/

### binprovider/host.go
- `Host` interface for binary providers:
  - `FS()`, `Sudo()`, `IsWindows()`
  - `Arch()`, `InstalledK3sVersion()`, `RunningK3sVersion()`
  - `Dir()`, `OSKind()`, `DownloadURL()`, `Touch()`, `DeleteFile()`, `SetFileMode()`

### binprovider/helpers.go
- `stagedFile` struct: embedded by providers for temp file management
- `needsUpgrade()` — shared upgrade check logic
- `stageTempPath()` — generate temp file path (preserves .exe on Windows)
- `stageUpload()` — upload local file to temp path on remote host
- `stageDownload()` — download URL to temp path on remote host

### binprovider/download.go
- `downloadProvider` struct: downloads binary directly on remote host
- `NewGitHub(host, installPath, target)` — fetches from K3s GitHub releases
  - URL format: `https://github.com/k3s-io/k3s/releases/download/{version}/k3s` (amd64)
  - ARM64: `k3s-arm64`, ARM: `k3s-armhf`, s390x: `k3s-s390x`
- `NewCustomURL(host, installPath, urlFor, target)` — downloads from custom URL

### binprovider/local_upload.go
- `localUpload` struct: downloads to local cache, then uploads to host
- Uses `http.Client` with configured timeouts
- Cache path: `~/.cache/k3sctl/k3s/{os}/{arch}/k3s-{version}`
- `downloadToFile()` — HTTP GET with temp file + atomic rename
- Implements `BinaryCacher` interface

### binprovider/local_file.go
- `localFile` struct: uploads developer-supplied binary from local machine
- `NewLocalFile(host, localPath, installPath, changed)` — uses `changed()` callback for upgrade detection

### binprovider/existing.go
- `existing` struct: uses already-present binary on host
- `NeedsUpgrade()`: compares installed vs running version
- `Stage()`: returns ("", nil) — no staging needed

## pkg/manifest/reader.go
- Multi-document YAML reader using `k8s.io/apimachinery/pkg/util/yaml`
- `ResourceDefinition` struct: APIVersion, Kind, Metadata.Name, Origin, Raw
- `Reader` struct with methods:
  - `Parse(io.Reader)`, `ParseString()`, `ParseBytes()`
  - `Resources()` — all parsed resources
  - `FilterResources(filter)` — filtered resources
  - `GetResources(apiVersion, kind)` — resources by type

### pkg/manifest/reader_test.go
- Manifest reader tests

## pkg/node/statusfunc.go
- Status-checking functions used with retry:
  - `KubeNodeReadyFunc(host)` — checks node Ready condition via `k3s kubectl get node -o json`
  - `ScheduledEventsAfterFunc(host, since)` — checks for Scheduled events after time
  - `HTTPStatusFunc(host, url, expected...)` — HTTP status check
  - `ServiceRunningFunc(host, service)` — systemd service running check
  - `ServiceStoppedFunc(host, service)` — systemd service stopped check

### pkg/node/statusfunc_test.go
- Status function tests

## integration/github/github.go
- GitHub API integration for version checks:
  - `Release` struct: URL, TagName, PreRelease, Assets
  - `LatestK3sBinaryURL(arch, osKind, preok)` — from k3s-io/k3s releases
  - `LatestK3sVersion(preok)` — latest K3s version string
  - `LatestRelease(preok)` — latest k3sctl release from nirantarasi/k3sctl
  - `unmarshalURLBody()` — generic HTTP GET + JSON unmarshal
  - K3s release tags: `v1.28.2+k3s1` format
  - K3s binary naming: `k3s` (amd64), `k3s-arm64`, `k3s-armhf`

## Implementation Order
1. `pkg/retry/` (utility)
2. `pkg/download/` (with cache path)
3. `pkg/manifest/` (manifest reading)
4. `pkg/k3s/provider.go` (interface)
5. `pkg/k3s/binprovider/host.go` (interface)
6. `pkg/k3s/binprovider/helpers.go`
7. `pkg/k3s/binprovider/download.go` (K3s URLs)
8. `pkg/k3s/binprovider/local_upload.go` (K3s cache)
9. `pkg/k3s/binprovider/local_file.go`
10. `pkg/k3s/binprovider/existing.go`
11. `pkg/node/statusfunc.go` (K3s adaptations)
12. `integration/github/github.go` (K3s GitHub API)
