//go:build integration

package lxd_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nirantarasi/k3sctl/integration/environment"
	"github.com/nirantarasi/k3sctl/integration/lxd"
	"github.com/nirantarasi/k3sctl/integration/multinode"
	"github.com/nirantarasi/k3sctl/pkg/apis/k3sctl.k3sproject.io/v1beta1/cluster"
)

// TestLXDConfigReconciliation verifies that changing K3s config.yaml on an
// existing cluster and re-running apply actually restarts K3s so the change
// takes effect. This is the core test for the config reconciliation fix.
//
// Scenario:
//  1. Install a single-server+agent cluster with Traefik enabled (default).
//  2. Verify Traefik is running.
//  3. Re-apply with disable.components: [traefik, servicelb].
//  4. Verify Traefik is no longer running.
//  5. Re-apply with the same config again (idempotency check).
func TestLXDConfigReconciliation(t *testing.T) {
	setup := newLXDTestSetup(t, "reconcile")

	workDir := t.TempDir()
	keyPair, err := multinode.GenerateSSHKey(filepath.Join(workDir, "keys"))
	if err != nil {
		t.Fatalf("GenerateSSHKey() error = %v", err)
	}

	provider, err := lxd.New(context.Background(), lxd.Config{
		RunID:            setup.runID,
		ImageAlias:       "24.04",
		SSHUser:          "ubuntu",
		SSHPort:          22,
		SSHAuthorizedKey: keyPair.AuthorizedKey,
		Nodes: []lxd.NodeSpec{
			{Name: "server", Role: environment.RoleServer},
			{Name: "agent1", Role: environment.RoleAgent},
		},
	})
	if err != nil {
		t.Fatalf("connect to LXD: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := provider.Destroy(cleanupCtx); err != nil {
			t.Logf("fallback LXD cleanup: %v", err)
		}
	})

	sshRunner := &multinode.SSHRunner{
		User:    "ubuntu",
		Port:    22,
		KeyPath: keyPair.PrivateKeyPath,
		Timeout: 30 * time.Second,
	}

	runner, err := multinode.NewRunner(multinode.Options{
		Binary:       setup.binary,
		WorkDir:      workDir,
		KeyPath:      keyPair.PrivateKeyPath,
		ArtifactDir:  setup.artifactDir,
		ClusterName:  "k3sctl-" + setup.runID,
		K3sVersion:   setup.k3sVersion,
		Environment:  provider,
		Remote:       sshRunner,
		ApplyTimeout: 25 * time.Minute,
	})
	if err != nil {
		t.Fatalf("NewRunner() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
	defer cancel()

	// --- Step 1: Initial install with defaults (Traefik enabled) ---
	t.Log("Step 1: initial install with Traefik enabled")
	if err := runner.Run(ctx); err != nil {
		t.Fatalf("initial apply failed: %v", err)
	}

	// Retrieve the node addresses from the LXD containers.
	nodes, err := provider.Create(ctx)
	if err != nil {
		// The containers already exist from Run(), Create returns the existing
		// addresses. If this fails, read them from the provider config.
		t.Fatalf("get node addresses: %v", err)
	}
	serverNode := findNodeByRole(nodes, environment.RoleServer)
	if serverNode == nil {
		t.Fatal("no server node found")
	}

	// --- Step 2: Verify Traefik is running ---
	t.Log("Step 2: verify Traefik is running")
	waitForPodState(t, sshRunner, ctx, serverNode.Address, "traefik", true)

	// --- Step 3: Re-apply with Traefik and ServiceLB disabled ---
	t.Log("Step 3: re-apply with Traefik and ServiceLB disabled")
	runner.SetK3sOverrides(&multinode.K3sOverrides{
		Disable: &cluster.DisableConfig{
			Components: []string{"traefik", "servicelb"},
		},
	})

	configPath, err := runner.WriteConfig(nodes)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	if err := runner.Apply(ctx, configPath); err != nil {
		t.Fatalf("re-apply with disabled components failed: %v", err)
	}

	// --- Step 4: Verify config.yaml has the disable field ---
	t.Log("Step 4: verify config.yaml has disable field")
	verifyK3sConfigKey(t, sshRunner, ctx, serverNode.Address, "disable")

	// --- Step 5: Verify Traefik is no longer running ---
	t.Log("Step 5: verify Traefik is stopped")
	waitForPodState(t, sshRunner, ctx, serverNode.Address, "traefik", false)

	// --- Step 6: Idempotency — re-apply with same config, no unnecessary restart ---
	t.Log("Step 6: idempotency check — re-apply same config")
	if err := runner.Apply(ctx, configPath); err != nil {
		t.Fatalf("idempotent re-apply failed: %v", err)
	}

	// Cluster should still be healthy
	if err := runner.Verify(ctx, nodes, configPath); err != nil {
		t.Fatalf("verification after idempotent re-apply failed: %v", err)
	}
}

// waitForPodState waits for a pod name substring to be present (running=true) or
// absent (running=false) in the pod list. Polls for up to 3 minutes.
func waitForPodState(t *testing.T, ssh *multinode.SSHRunner, ctx context.Context, address, podSubstr string, wantRunning bool) {
	t.Helper()

	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		result, err := ssh.Run(ctx, address,
			"sudo /usr/local/bin/k3s kubectl get pods -A --no-headers 2>/dev/null || true")
		if err != nil {
			t.Logf("get pods failed (will retry): %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		found := false
		for _, line := range strings.Split(result.Stdout, "\n") {
			if strings.Contains(line, podSubstr) && strings.Contains(line, "Running") {
				found = true
				break
			}
		}

		if found == wantRunning {
			return
		}

		time.Sleep(10 * time.Second)
	}

	if wantRunning {
		t.Errorf("timed out waiting for pod %q to be running on %s", podSubstr, address)
	} else {
		t.Errorf("timed out waiting for pod %q to stop on %s", podSubstr, address)
	}
}
