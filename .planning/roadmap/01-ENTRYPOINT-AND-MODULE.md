# 01: Entry Point, Module Setup, and Version

## Files

### `go.mod`
- **Module**: `github.com/nirantarasi/k3sctl`
- **Go version**: 1.27.0
- **Key dependencies**:
  - `github.com/urfave/cli/v2` — CLI framework
  - `github.com/k3sproject/rig/v2` — SSH remote execution
  - `github.com/creasty/defaults` — struct defaults
  - `github.com/jellydator/validation` — validation
  - `github.com/logrusorgru/aurora` — colored output
  - `github.com/AlecAivazis/survey/v2` — interactive prompts
  - `github.com/carlmjohnson/versioninfo` — build version info
  - `k8s.io/client-go` — kubeconfig manipulation
  - `gopkg.in/yaml.v2` — YAML parsing
  - `github.com/mattn/go-isatty` — terminal detection
  - `github.com/a8m/envsubst` — env var substitution
  - `github.com/adrg/xdg` — XDG paths
  - `github.com/bmatcuk/doublestar/v4` — glob patterns
  - `github.com/sergi/go-diff` — diff output
  - `github.com/stretchr/testify` — test assertions
  - `github.com/samber/slog-logrus/v2` — slog bridge for rig

### `main.go`
- Import `github.com/nirantarasi/k3sctl/cmd`
- Import `github.com/nirantarasi/k3sctl/version` (blank)
- `cmd.NewK3sctl(os.Stdin, os.Stdout, os.Stderr)`

### `version/version.go`
- Variables: `Version`, `GitCommit`, `Environment`
- Function: `IsPre() bool` — checks for prerelease
- Version/GitCommit/Environment injected at build time via ldflags

## Structure

```go
// main.go
package main
import (
    "os"
    "github.com/nirantarasi/k3sctl/cmd"
    _ "github.com/carlmjohnson/versioninfo"
    _ "github.com/nirantarasi/k3sctl/version"
)
func main() {
    k3sctl := cmd.NewK3sctl(os.Stdin, os.Stdout, os.Stderr)
    if err := k3sctl.Run(os.Args); err != nil { log.Fatal(err) }
}
```

```go
// version/version.go
package version
var (
    Version     = versioninfo.Version
    GitCommit   = versioninfo.Revision
    Environment = "development"
)
func IsPre() bool { return strings.Contains(Version, "-") }
```
