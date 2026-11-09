# 12: Binary Providers — K3s Binary Acquisition

## Overview

Binary providers implement the strategy pattern for acquiring K3s binaries on target hosts.

## `pkg/k3s/provider.go`
- Package: `k3s`
- `BinaryProvider` interface:
  - `NeedsUpgrade() bool`
  - `IsUpload() bool`
  - `Stage(ctx context.Context) (string, error)`
  - `CleanUp(ctx context.Context)`
- `BinaryCacher` interface:
  - `BinaryCacheKey() (string, error)`
  - `EnsureCached(ctx context.Context) error`
- Module path: `github.com/nirantarasi/k3sctl/pkg/k3s`

## `pkg/k3s/binprovider/host.go`
- `Host` interface:
  - `FS()`, `Sudo()`, `Arch()`, `Dir()`, `OSKind()`, `DownloadURL()`, `Touch()`, `DeleteFile()`, `SetFileMode()`
  - `InstalledK3sVersion() *version.Version`
  - `RunningK3sVersion() *version.Version`
- No `IsWindows()` — K3s is Linux-only

## `pkg/k3s/binprovider/helpers.go`
- `stagedFile` struct: holds host + tmpPath
- `CleanUp()`: removes temp binary
- `needsUpgrade()`: compare target vs binary vs running
- `stageTempPath()`: Linux-only (no .exe handling)
- `stageUpload()`: dir creation, upload, chmod 0755
- `stageDownload()`: dir creation, download, chmod

## `pkg/k3s/binprovider/download.go`
- `downloadProvider` struct
- `NewGitHub()` — K3s download URL pattern:
  - `https://github.com/k3s-io/k3s/releases/download/{version}/k3s` for amd64
  - `https://github.com/k3s-io/k3s/releases/download/{version}/k3s-arm64` for arm64
  - `https://github.com/k3s-io/k3s/releases/download/{version}/k3s-armhf` for arm
  - Version format: `v1.30.0+k3s1` (the `+` must be URL-encoded as `%2B`)
  - No Windows binary
- `NewCustomURL()`: user-provided URL resolver
- `urlFor` function: arch-to-binary-suffix mapping

### K3s Download URL Builder
```go
func k3sDownloadURL(version, arch string) string {
    encodedVersion := strings.ReplaceAll(version, "+", "%2B")
    base := fmt.Sprintf("https://github.com/k3s-io/k3s/releases/download/%s/k3s", encodedVersion)
    switch arch {
    case "arm64", "aarch64":
        return base + "-arm64"
    case "arm", "armv7l":
        return base + "-armhf"
    default: // amd64
        return base
    }
}
```

## `pkg/k3s/binprovider/local_upload.go`
- `localUpload` struct
- `downloadHTTPClient`: HTTP client config
- `cacheFilePath()`: `k3sctl/k3s/linux/{arch}/k3s-{version}` — no `.exe`, no Windows
- `BinaryCacheKey()`: cache key logic
- `EnsureCached()`: download URL uses K3s URL builder
- `downloadToFile()`: HTTP download with atomic rename
- `Stage()`: upload from local cache to remote
- `NewLocalUpload()`: constructor

## `pkg/k3s/binprovider/local_file.go`
- `localFile` struct
- `NewLocalFile()`: constructor

## `pkg/k3s/binprovider/existing.go`
- Uses `InstalledK3sVersion()` and `RunningK3sVersion()`
- Reuses existing K3s binary on host

## K3s Version Handling

K3s versions use `+k3s1` suffix (e.g., `v1.30.0+k3s1`). The `+` in the version string requires URL encoding for downloads.

## Test Files

### `pkg/k3s/binprovider/download_test.go`
- Test K3s URL generation for amd64, arm64, armhf
- Test URL encoding of `+` in version strings

### `pkg/k3s/binprovider/helpers_test.go`
- Test `needsUpgrade()` logic
- Test `stageTempPath()` without Windows paths
