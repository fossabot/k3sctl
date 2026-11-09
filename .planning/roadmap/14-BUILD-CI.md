# 14: Build System, CI, and Release

## Taskfile.yml

- Binary name: `k3sctl`
- LD_FLAGS: `github.com/nirantarasi/k3sctl/version.Environment`
- Build targets: cross-platform (no Windows for K3s)
- Keep: linux-amd64, linux-arm64, linux-arm, linux-riscv64, darwin-amd64, darwin-arm64
- Lint: golangci-lint integration
- Test: `go test -v ./...`

### k3sctl Task Targets
```
build:       # Build for local platform
build-all:   # Cross-compile all platforms
test:        # Run unit tests
test-race:   # Run tests with race detector
lint:        # Run golangci-lint
vet:         # Run go vet
fmt:         # Format code
check-fmt:   # Check formatting
check:       # check-fmt + lint
clean:       # Remove build artifacts
install:     # Install to /usr/local/bin
generate:    # Run go generate
```

## Build Flags
```
LD_FLAGS = -s -w \
  -X github.com/nirantarasi/k3sctl/version.Environment=$(ENVIRONMENT) \
  -X github.com/carlmjohnson/versioninfo.Revision=$(GIT_COMMIT) \
  -X github.com/carlmjohnson/versioninfo.Version=$(TAG_NAME)

BUILD_FLAGS = -trimpath -a \
  -tags "netgo,osusergo,static_build" \
  -installsuffix netgo \
  -ldflags "$(LD_FLAGS) -extldflags '-static'"
```

## GitHub Actions Workflows

### `.github/workflows/ci.yml` — Unit Tests & Build
- Trigger: PRs touching `**.go`, `go.mod`, `go.sum`, `Taskfile.yml`
- Matrix: ubuntu-latest, macos-latest (no Windows)
- Steps: checkout → setup-go → `task check-fmt` → `task vet` → `task test-race` → `task build` → `task build-all`

### `.github/workflows/cli-smoke.yml` — CLI Smoke Tests
- Builds binary and asserts it works on ubuntu-latest and macos-latest
- Tests: version, help, completion, init output, error handling

### `.github/workflows/real-install.yml` — Real Install Test
- On ubuntu-24.04, builds CLI and runs full apply/reset cycle
- Nightly + manual dispatch + path-filtered pushes

### `.github/workflows/release.yml` — Release
- Trigger: push tags `v*`
- Build: `task build-all`
- Upload artifacts: `bin/k3sctl-*`
- GitHub release: draft with checksums

### `.github/workflows/codeql.yml` — Security
- CodeQL analysis for Go

## Dockerfile
- Multi-stage build
- Binary: `k3sctl`
- Build args: ENVIRONMENT, GIT_COMMIT, TAG_NAME
- Entrypoint: `["/k3sctl"]`

## `.golangci.yml`
- golangci-lint v2 config

## Cross-Platform Build Matrix

| Target | GOOS | GOARCH | Binary Name |
|---|---|---|---|
| Linux amd64 | linux | amd64 | k3sctl-linux-amd64 |
| Linux arm64 | linux | arm64 | k3sctl-linux-arm64 |
| Linux arm | linux | arm | k3sctl-linux-arm |
| Linux riscv64 | linux | riscv64 | k3sctl-linux-riscv64 |
| macOS amd64 | darwin | amd64 | k3sctl-darwin-amd64 |
| macOS arm64 | darwin | arm64 | k3sctl-darwin-arm64 |

Note: No Windows build — K3s is Linux-only for the server/agent, and the CLI tool targets Linux/macOS users.
