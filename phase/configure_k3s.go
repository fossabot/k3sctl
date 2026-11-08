package phase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"strings"

	"log/slog"

	"github.com/nirantarasi/k3sctl/pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster"
	"gopkg.in/yaml.v2"
)

type ConfigureK3s struct {
	GenericPhase
	leader *cluster.Host
}

func (p *ConfigureK3s) Title() string {
	return "Configure K3s"
}

func (p *ConfigureK3s) Prepare() {
	p.leader = p.Config.Spec.K3sLeader()
}

func (p *ConfigureK3s) Run(ctx context.Context) error {
	clusterToken, err := p.resolveClusterToken()
	if err != nil {
		return fmt.Errorf("resolve cluster token: %w", err)
	}

	return p.parallelDo(ctx, p.Config.Spec.Hosts, func(ctx context.Context, h *cluster.Host) error {
		configDir := h.Configurer.K3sConfigDirPath()
		if err := h.Sudo().Exec(fmt.Sprintf("mkdir -p %s", shellQuote(configDir))); err != nil {
			return fmt.Errorf("%s: create config dir: %w", h, err)
		}

		var configData []byte
		if h.IsServer() {
			isLeader := p.leader != nil && h == p.leader
			var joinAddr string
			if !isLeader && p.leader != nil {
				joinAddr = p.leader.Address()
			}
			configData, err = p.Config.Spec.K3s.GenerateServerConfig(h, clusterToken, isLeader, joinAddr)
		} else {
			serverAddr := ""
			if p.leader != nil {
				serverAddr = p.leader.Address()
			}
			configData, err = p.Config.Spec.K3s.GenerateAgentConfig(h, serverAddr, clusterToken)
		}
		if err != nil {
			return fmt.Errorf("%s: generate config: %w", h, err)
		}

		configChanged := h.Metadata.K3sServiceRunning &&
			!configBytesEqual([]byte(h.Metadata.K3sExistingConfig), configData)

		configPath := h.Configurer.K3sConfigPath()
		if err := h.Sudo().FS().WriteFile(configPath, configData, fs.FileMode(0600)); err != nil {
			return fmt.Errorf("%s: write config: %w", h, err)
		}
		slog.Info(fmt.Sprintf("%s: wrote K3s config to %s", h, configPath))

		if configChanged {
			h.Metadata.NeedsConfigRestart = true
			slog.Info(fmt.Sprintf("%s: config.yaml changed, will restart K3s", h))
			p.cleanupStaleSkipFiles(h)
		}

		return nil
	})
}

// resolveClusterToken determines the token every node will use. On an existing
// cluster the token is read from the leader during fact gathering. When the
// cluster is being initialized it is either taken from the user config or
// generated up front, so follower and agent configs can be written with a
// usable token before the leader starts.
func (p *ConfigureK3s) resolveClusterToken() (string, error) {
	if p.leader == nil {
		return "", fmt.Errorf("no leader server configured")
	}
	if p.leader.Metadata.NodeToken != "" {
		return p.leader.Metadata.NodeToken, nil
	}

	if p.Config.Spec.K3s == nil {
		return "", fmt.Errorf("k3s spec is not configured")
	}
	if token, ok := p.Config.Spec.K3s.Config["token"].(string); ok && token != "" {
		p.leader.Metadata.NodeToken = token
		return token, nil
	}

	if p.leader.Metadata.K3sServiceRunning {
		return "", fmt.Errorf("leader %s is running but its node token is unknown", p.leader)
	}

	token, err := generateClusterToken()
	if err != nil {
		return "", err
	}
	p.leader.Metadata.NodeToken = token
	slog.Info(fmt.Sprintf("%s: generated cluster token", p.leader))
	return token, nil
}

func generateClusterToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// configBytesEqual compares two YAML config byte slices for semantic equality
// by normalizing trailing whitespace. This prevents spurious restarts when the
// only difference is a trailing newline.
func configBytesEqual(existing, desired []byte) bool {
	return strings.TrimSpace(string(existing)) == strings.TrimSpace(string(desired))
}

// cleanupStaleSkipFiles removes .skip files in the K3s manifests directory for
// components that were previously disabled but are no longer in the disable
// list. K3s creates these files when --disable is set but does not remove them
// when the flag is dropped, so re-enabled components stay disabled after a
// restart unless the skip files are cleaned up.
func (p *ConfigureK3s) cleanupStaleSkipFiles(h *cluster.Host) {
	oldDisabled := disabledComponentsFromYAML(h.Metadata.K3sExistingConfig)
	newDisabled := make(map[string]bool)
	if d := p.Config.Spec.K3s.Disable; d != nil {
		for _, c := range d.Components {
			newDisabled[c] = true
		}
	}

	manifestsDir := fmt.Sprintf("%s/server/manifests", h.K3sDataDir())
	for _, comp := range oldDisabled {
		if newDisabled[comp] {
			continue
		}
		skipFile := fmt.Sprintf("%s/%s.skip", manifestsDir, comp)
		if h.Sudo().FS().FileExist(skipFile) {
			if err := h.Sudo().FS().Remove(skipFile); err != nil {
				slog.Warn(fmt.Sprintf("%s: failed to remove skip file %s: %v", h, skipFile, err))
			} else {
				slog.Info(fmt.Sprintf("%s: removed %s (component re-enabled)", h, skipFile))
			}
		}
	}
}

func disabledComponentsFromYAML(configYAML string) []string {
	if configYAML == "" {
		return nil
	}
	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(configYAML), &parsed); err != nil {
		return nil
	}
	raw, ok := parsed["disable"]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case string:
		return []string{v}
	}
	return nil
}
