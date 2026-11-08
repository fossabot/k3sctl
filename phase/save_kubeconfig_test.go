package phase

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nirantarasi/k3sctl/pkg/apis/k3sctl.k3sproject.io/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveKubeconfig_ShouldRun(t *testing.T) {
	t.Run("true when kubeconfig is set", func(t *testing.T) {
		p := &SaveKubeconfig{}
		p.Config = &v1beta1.Cluster{
			Metadata: v1beta1.ClusterMetadata{Kubeconfig: testKubeconfig},
		}
		assert.True(t, p.ShouldRun())
	})

	t.Run("false when kubeconfig is empty", func(t *testing.T) {
		p := &SaveKubeconfig{}
		p.Config = &v1beta1.Cluster{}
		assert.False(t, p.ShouldRun())
	})
}

func TestSaveKubeconfig_Run(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".kube", "config")

	p := &SaveKubeconfig{Path: path}
	p.Config = &v1beta1.Cluster{
		Metadata: v1beta1.ClusterMetadata{Kubeconfig: testKubeconfig},
	}

	err := p.Run(t.Context())
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "https://127.0.0.1:6443")
}

func TestSaveKubeconfig_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "config")

	p := &SaveKubeconfig{Path: path}
	p.Config = &v1beta1.Cluster{
		Metadata: v1beta1.ClusterMetadata{Kubeconfig: testKubeconfig},
	}

	err := p.Run(t.Context())
	require.NoError(t, err)

	info, err := os.Stat(filepath.Dir(path))
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
