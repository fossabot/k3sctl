package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestK3sSetDefaults(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	assert.NotNil(t, k.Config)
}

func TestK3sValidate(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"empty version is valid", "", false},
		{"version with v prefix", "v1.28.3+k3s1", false},
		{"version without v prefix", "1.28.3+k3s1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &K3s{Version: tt.version}
			err := k.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateServerConfig_Leader(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "my-token", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, true, config["cluster-init"])
	assert.Equal(t, "my-token", config["token"])
	_, hasServer := config["server"]
	assert.False(t, hasServer, "leader should not have server field")
}

func TestGenerateServerConfig_JoinExisting(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "my-token", false, "10.0.0.1")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "https://10.0.0.1:6443", config["server"])
	assert.Equal(t, "my-token", config["token"])
	_, hasClusterInit := config["cluster-init"]
	assert.False(t, hasClusterInit, "joining server should not have cluster-init")
}

func TestGenerateServerConfig_PrivateAddress(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{
		Role:           RoleServer,
		PrivateAddress: "192.168.1.10",
	}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "192.168.1.10", config["node-ip"])
}

func TestGenerateServerConfig_NoFlannel(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{
		Role:      RoleServer,
		NoFlannel: true,
	}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "none", config["flannel-backend"])
}

func TestGenerateServerConfig_ExtraArgs(t *testing.T) {
	k := &K3s{
		ExtraArgs: Flags{"--tls-san=10.0.0.100", "--disable=traefik", "--write-kubeconfig-mode=644"},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "10.0.0.100", config["tls-san"])
	assert.Equal(t, "traefik", config["disable"])
	assert.Equal(t, "644", config["write-kubeconfig-mode"])
}

func TestGenerateServerConfig_ExtraArgsBoolFlag(t *testing.T) {
	k := &K3s{
		ExtraArgs: Flags{"--secrets-encryption"},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, true, config["secrets-encryption"])
}

func TestGenerateServerConfig_BaseConfigPassedThrough(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	k.Config["write-kubeconfig-mode"] = "644"
	k.Config["cluster-cidr"] = "10.42.0.0/16"
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "644", config["write-kubeconfig-mode"])
	assert.Equal(t, "10.42.0.0/16", config["cluster-cidr"])
}

func TestGenerateAgentConfig_Basic(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{Role: RoleAgent}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "agent-token")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "https://10.0.0.1:6443", config["server"])
	assert.Equal(t, "agent-token", config["token"])
}

func TestGenerateAgentConfig_PrivateAddress(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{
		Role:           RoleAgent,
		PrivateAddress: "192.168.1.20",
	}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "token")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "192.168.1.20", config["node-ip"])
}

func TestGenerateAgentConfig_FiltersServerOnlyFields(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	k.Config["cluster-init"] = true
	k.Config["cluster-cidr"] = "10.42.0.0/16"
	k.Config["service-cidr"] = "10.43.0.0/16"
	k.Config["cluster-dns"] = "10.43.0.10"
	k.Config["disable"] = "traefik"
	k.Config["tls-san"] = "my.domain.com"
	k.Config["write-kubeconfig-mode"] = "644"
	k.Config["kube-apiserver-arg"] = "some-arg"
	k.Config["kube-controller-manager-arg"] = "some-arg"
	k.Config["kube-scheduler-arg"] = "some-arg"
	k.Config["etcd-expose-metrics"] = true
	k.Config["node-label"] = "foo=bar"
	host := &Host{Role: RoleAgent}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "token")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	serverOnlyKeys := []string{
		"cluster-init", "cluster-cidr", "service-cidr", "cluster-dns",
		"disable", "tls-san", "etcd-expose-metrics",
		"write-kubeconfig-mode", "kube-apiserver-arg",
		"kube-controller-manager-arg", "kube-scheduler-arg",
	}
	for _, key := range serverOnlyKeys {
		_, present := config[key]
		assert.False(t, present, "agent config should not contain server-only field %q", key)
	}

	assert.Equal(t, "foo=bar", config["node-label"], "non-server fields should pass through")
}

func TestGenerateAgentConfig_EmptyToken(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{Role: RoleAgent}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	_, hasToken := config["token"]
	assert.False(t, hasToken, "empty token should be omitted from agent config")
}

func TestK3sValidate_EtcdSnapshotS3RequiresBucket(t *testing.T) {
	tests := []struct {
		name    string
		k3s     K3s
		wantErr bool
	}{
		{
			"S3 enabled without bucket",
			K3s{EtcdSnapshots: &EtcdSnapshotConfig{S3: &EtcdS3Config{Enabled: true}}},
			true,
		},
		{
			"S3 enabled with bucket",
			K3s{EtcdSnapshots: &EtcdSnapshotConfig{S3: &EtcdS3Config{Enabled: true, Bucket: "my-bucket"}}},
			false,
		},
		{
			"S3 disabled without bucket",
			K3s{EtcdSnapshots: &EtcdSnapshotConfig{S3: &EtcdS3Config{Enabled: false}}},
			false,
		},
		{
			"nil S3",
			K3s{EtcdSnapshots: &EtcdSnapshotConfig{}},
			false,
		},
		{
			"nil etcdSnapshots",
			K3s{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.k3s.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "bucket")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateServerConfig_EtcdSnapshots(t *testing.T) {
	t.Run("scheduled snapshots with cron and retention", func(t *testing.T) {
		k := &K3s{
			EtcdSnapshots: &EtcdSnapshotConfig{
				Scheduled:    true,
				ScheduleCron: "0 */12 * * *",
				Retention:    5,
				Compress:     true,
				Dir:          "/snapshots",
			},
		}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, "0 */12 * * *", config["etcd-snapshot-schedule-cron"])
		assert.Equal(t, 5, config["etcd-snapshot-retention"])
		assert.Equal(t, true, config["etcd-snapshot-compress"])
		assert.Equal(t, "/snapshots", config["etcd-snapshot-dir"])
	})

	t.Run("unscheduled does not set cron fields", func(t *testing.T) {
		k := &K3s{
			EtcdSnapshots: &EtcdSnapshotConfig{
				Scheduled:    false,
				ScheduleCron: "0 */12 * * *",
				Retention:    5,
			},
		}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		_, hasCron := config["etcd-snapshot-schedule-cron"]
		assert.False(t, hasCron, "cron should not be set when scheduled=false")
		_, hasRetention := config["etcd-snapshot-retention"]
		assert.False(t, hasRetention, "retention should not be set when scheduled=false")
	})

	t.Run("S3 config", func(t *testing.T) {
		k := &K3s{
			EtcdSnapshots: &EtcdSnapshotConfig{
				S3: &EtcdS3Config{
					Enabled:       true,
					Endpoint:      "s3.example.com",
					EndpointCA:    "/ca.crt",
					SkipSSLVerify: true,
					AccessKey:     "AKID",
					SecretKey:     "secret",
					Bucket:        "backups",
					Region:        "us-east-1",
					Folder:        "k3s",
					Insecure:      true,
				},
			},
		}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, true, config["etcd-s3"])
		assert.Equal(t, "s3.example.com", config["etcd-s3-endpoint"])
		assert.Equal(t, "/ca.crt", config["etcd-s3-endpoint-ca"])
		assert.Equal(t, true, config["etcd-s3-skip-ssl-verify"])
		assert.Equal(t, "AKID", config["etcd-s3-access-key"])
		assert.Equal(t, "secret", config["etcd-s3-secret-key"])
		assert.Equal(t, "backups", config["etcd-s3-bucket"])
		assert.Equal(t, "us-east-1", config["etcd-s3-region"])
		assert.Equal(t, "k3s", config["etcd-s3-folder"])
		assert.Equal(t, true, config["etcd-s3-insecure"])
	})

	t.Run("S3 disabled produces no S3 keys", func(t *testing.T) {
		k := &K3s{
			EtcdSnapshots: &EtcdSnapshotConfig{
				S3: &EtcdS3Config{
					Enabled: false,
					Bucket:  "backups",
				},
			},
		}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		_, hasS3 := config["etcd-s3"]
		assert.False(t, hasS3, "etcd-s3 should not be set when S3 is disabled")
	})

	t.Run("nil etcdSnapshots produces no snapshot keys", func(t *testing.T) {
		k := &K3s{}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		for key := range config {
			assert.False(t, len(key) > 4 && key[:4] == "etcd", "no etcd keys expected, got %q", key)
		}
	})
}

func TestGenerateServerConfig_SecretsEncrypt(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		k := &K3s{SecretsEncrypt: &SecretsEncryptConfig{Enabled: true}}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, true, config["secrets-encryption"])
	})

	t.Run("disabled", func(t *testing.T) {
		k := &K3s{SecretsEncrypt: &SecretsEncryptConfig{Enabled: false}}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		_, has := config["secrets-encryption"]
		assert.False(t, has)
	})

	t.Run("nil", func(t *testing.T) {
		k := &K3s{}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		_, has := config["secrets-encryption"]
		assert.False(t, has)
	})
}

func TestGenerateServerConfig_TokenConfig(t *testing.T) {
	t.Run("both tokens", func(t *testing.T) {
		k := &K3s{Token: &TokenConfig{
			ServerToken: "server-secret",
			AgentToken:  "agent-secret",
		}}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, "server-secret", config["token"])
		assert.Equal(t, "agent-secret", config["agent-token"])
	})

	t.Run("server token overrides cluster token", func(t *testing.T) {
		k := &K3s{Token: &TokenConfig{ServerToken: "from-yaml"}}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "from-cluster", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, "from-yaml", config["token"], "token config should override clusterToken")
	})

	t.Run("agent token only", func(t *testing.T) {
		k := &K3s{Token: &TokenConfig{AgentToken: "agent-only"}}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "cluster-tok", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, "cluster-tok", config["token"], "cluster token should remain when no server token in config")
		assert.Equal(t, "agent-only", config["agent-token"])
	})

	t.Run("nil token config", func(t *testing.T) {
		k := &K3s{}
		k.SetDefaults()
		host := &Host{Role: RoleServer}

		data, err := k.GenerateServerConfig(host, "cluster-tok", true, "")
		require.NoError(t, err)

		var config map[string]any
		require.NoError(t, yaml.Unmarshal(data, &config))

		assert.Equal(t, "cluster-tok", config["token"])
		_, hasAgent := config["agent-token"]
		assert.False(t, hasAgent)
	})
}

func TestGenerateServerConfig_LabelsAndTaints(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{
		Role: RoleServer,
		Labels: map[string]string{
			"environment": "production",
			"tier":        "control-plane",
		},
		Taints: []string{
			"node-role.kubernetes.io/master:NoSchedule",
		},
	}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	labels, ok := config["node-label"].([]interface{})
	require.True(t, ok, "node-label should be a list")
	assert.Len(t, labels, 2)

	taints, ok := config["node-taint"].([]interface{})
	require.True(t, ok, "node-taint should be a list")
	assert.Len(t, taints, 1)
	assert.Equal(t, "node-role.kubernetes.io/master:NoSchedule", taints[0])
}

func TestGenerateAgentConfig_LabelsAndTaints(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{
		Role: RoleAgent,
		Labels: map[string]string{
			"pool": "workers",
		},
		Taints: []string{
			"dedicated=gpu:NoSchedule",
		},
	}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "token")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	labels, ok := config["node-label"].([]interface{})
	require.True(t, ok, "node-label should be a list")
	assert.Len(t, labels, 1)
	assert.Equal(t, "pool=workers", labels[0])

	taints, ok := config["node-taint"].([]interface{})
	require.True(t, ok, "node-taint should be a list")
	assert.Len(t, taints, 1)
	assert.Equal(t, "dedicated=gpu:NoSchedule", taints[0])
}

func TestK3sValidate_SecretsEncryptRequiresEnabled(t *testing.T) {
	tests := []struct {
		name    string
		k3s     K3s
		wantErr bool
	}{
		{
			"rotateKeys without enabled",
			K3s{SecretsEncrypt: &SecretsEncryptConfig{RotateKeys: true}},
			true,
		},
		{
			"reencrypt without enabled",
			K3s{SecretsEncrypt: &SecretsEncryptConfig{Reencrypt: true}},
			true,
		},
		{
			"rotateKeys with enabled",
			K3s{SecretsEncrypt: &SecretsEncryptConfig{Enabled: true, RotateKeys: true}},
			false,
		},
		{
			"reencrypt with enabled",
			K3s{SecretsEncrypt: &SecretsEncryptConfig{Enabled: true, Reencrypt: true}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.k3s.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestK3sValidate_CertRotateCARequiresEnabled(t *testing.T) {
	k := K3s{CertRotate: &CertRotateConfig{RotateCA: true}}
	assert.Error(t, k.Validate())

	k2 := K3s{CertRotate: &CertRotateConfig{Enabled: true, RotateCA: true}}
	assert.NoError(t, k2.Validate())
}

func TestK3sValidate_EtcdPruneRequiresRetention(t *testing.T) {
	k := K3s{EtcdSnapshots: &EtcdSnapshotConfig{
		Prune: &EtcdPruneConfig{Enabled: true, Retention: 0},
	}}
	assert.Error(t, k.Validate())

	k2 := K3s{EtcdSnapshots: &EtcdSnapshotConfig{
		Prune: &EtcdPruneConfig{Enabled: true, Retention: 5},
	}}
	assert.NoError(t, k2.Validate())
}

func TestK3sYAMLRoundTrip_NewFields(t *testing.T) {
	input := `
version: v1.30.0+k3s1
secretsEncrypt:
  enabled: true
  rotateKeys: true
  reencrypt: true
certRotate:
  enabled: true
  rotateCA: true
token:
  serverToken: srv
  rotate: true
etcdSnapshots:
  scheduled: true
  prune:
    enabled: true
    retention: 10
    name: on-demand
registries:
  mirrors:
    docker.io:
      endpoint:
        - https://registry.example.com
  configs:
    registry.example.com:
      auth:
        username: user
        password: pass
      tls:
        caFile: /ca.pem
        insecureSkipVerify: true
healthCheck:
  enabled: true
  timeout: 120s
`
	var k K3s
	require.NoError(t, yaml.Unmarshal([]byte(input), &k))

	require.NotNil(t, k.SecretsEncrypt)
	assert.True(t, k.SecretsEncrypt.RotateKeys)
	assert.True(t, k.SecretsEncrypt.Reencrypt)

	require.NotNil(t, k.CertRotate)
	assert.True(t, k.CertRotate.RotateCA)

	require.NotNil(t, k.Token)
	assert.True(t, k.Token.Rotate)

	require.NotNil(t, k.EtcdSnapshots)
	require.NotNil(t, k.EtcdSnapshots.Prune)
	assert.True(t, k.EtcdSnapshots.Prune.Enabled)
	assert.Equal(t, 10, k.EtcdSnapshots.Prune.Retention)
	assert.Equal(t, "on-demand", k.EtcdSnapshots.Prune.Name)

	require.NotNil(t, k.Registries)
	require.Len(t, k.Registries.Mirrors, 1)
	mirror := k.Registries.Mirrors["docker.io"]
	assert.Equal(t, []string{"https://registry.example.com"}, mirror.Endpoint)

	require.Len(t, k.Registries.Configs, 1)
	regCfg := k.Registries.Configs["registry.example.com"]
	require.NotNil(t, regCfg.Auth)
	assert.Equal(t, "user", regCfg.Auth.Username)
	assert.Equal(t, "pass", regCfg.Auth.Password)
	require.NotNil(t, regCfg.TLS)
	assert.Equal(t, "/ca.pem", regCfg.TLS.CAFile)
	assert.True(t, regCfg.TLS.InsecureSkipVerify)

	require.NotNil(t, k.HealthCheck)
	assert.True(t, k.HealthCheck.Enabled)
	assert.Equal(t, "120s", k.HealthCheck.Timeout)
}

func TestK3sYAMLRoundTrip(t *testing.T) {
	input := `
version: v1.30.0+k3s1
etcdSnapshots:
  scheduled: true
  scheduleCron: "0 */6 * * *"
  retention: 3
  s3:
    enabled: true
    bucket: my-bucket
    region: us-west-2
secretsEncrypt:
  enabled: true
certRotate:
  enabled: true
token:
  serverToken: srv-tok
  agentToken: agt-tok
preflight:
  enabled: true
  checkCgroups: true
  checkDisk: true
  minDiskGB: 10
`
	var k K3s
	require.NoError(t, yaml.Unmarshal([]byte(input), &k))

	assert.Equal(t, "v1.30.0+k3s1", k.Version)

	require.NotNil(t, k.EtcdSnapshots)
	assert.True(t, k.EtcdSnapshots.Scheduled)
	assert.Equal(t, "0 */6 * * *", k.EtcdSnapshots.ScheduleCron)
	assert.Equal(t, 3, k.EtcdSnapshots.Retention)
	require.NotNil(t, k.EtcdSnapshots.S3)
	assert.True(t, k.EtcdSnapshots.S3.Enabled)
	assert.Equal(t, "my-bucket", k.EtcdSnapshots.S3.Bucket)
	assert.Equal(t, "us-west-2", k.EtcdSnapshots.S3.Region)

	require.NotNil(t, k.SecretsEncrypt)
	assert.True(t, k.SecretsEncrypt.Enabled)

	require.NotNil(t, k.CertRotate)
	assert.True(t, k.CertRotate.Enabled)

	require.NotNil(t, k.Token)
	assert.Equal(t, "srv-tok", k.Token.ServerToken)
	assert.Equal(t, "agt-tok", k.Token.AgentToken)

	require.NotNil(t, k.Preflight)
	assert.True(t, k.Preflight.Enabled)
	assert.True(t, k.Preflight.CheckCgroups)
	assert.True(t, k.Preflight.CheckDisk)
	assert.Equal(t, 10, k.Preflight.MinDiskGB)
}

func TestGenerateServerConfig_NetworkConfig(t *testing.T) {
	k := &K3s{
		Network: &NetworkConfig{
			ClusterCIDR:        []string{"10.42.0.0/16", "fd42::/48"},
			ServiceCIDR:        []string{"10.43.0.0/16", "fd43::/112"},
			ClusterDNS:         []string{"10.43.0.10"},
			ClusterDomain:      "my.cluster",
			NodePortRange:      "30000-40000",
			FlannelBackend:     "wireguard-native",
			FlannelIPv6Masq:    true,
			FlannelExternalIP:  true,
			EgressSelectorMode: "cluster",
			ServiceLBNamespace: "lb-ns",
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "10.42.0.0/16,fd42::/48", config["cluster-cidr"])
	assert.Equal(t, "10.43.0.0/16,fd43::/112", config["service-cidr"])
	assert.Equal(t, "10.43.0.10", config["cluster-dns"])
	assert.Equal(t, "my.cluster", config["cluster-domain"])
	assert.Equal(t, "30000-40000", config["service-node-port-range"])
	assert.Equal(t, "wireguard-native", config["flannel-backend"])
	assert.Equal(t, true, config["flannel-ipv6-masq"])
	assert.Equal(t, true, config["flannel-external-ip"])
	assert.Equal(t, "cluster", config["egress-selector-mode"])
	assert.Equal(t, "lb-ns", config["servicelb-namespace"])
}

func TestGenerateServerConfig_NetworkConfig_Nil(t *testing.T) {
	k := &K3s{}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	_, has := config["cluster-cidr"]
	assert.False(t, has)
}

func TestGenerateServerConfig_DisableConfig(t *testing.T) {
	k := &K3s{
		Disable: &DisableConfig{
			Components:      []string{"traefik", "servicelb", "metrics-server"},
			Scheduler:       true,
			CloudController: true,
			KubeProxy:       true,
			NetworkPolicy:   true,
			HelmController:  true,
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	components, ok := config["disable"].([]interface{})
	require.True(t, ok)
	assert.Len(t, components, 3)
	assert.Equal(t, true, config["disable-scheduler"])
	assert.Equal(t, true, config["disable-cloud-controller"])
	assert.Equal(t, true, config["disable-kube-proxy"])
	assert.Equal(t, true, config["disable-network-policy"])
	assert.Equal(t, true, config["disable-helm-controller"])
}

func TestGenerateServerConfig_TLSConfig(t *testing.T) {
	k := &K3s{
		TLS: &TLSConfig{
			SANs:        []string{"lb.example.com", "10.0.0.100"},
			SANSecurity: true,
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	sans, ok := config["tls-san"].([]interface{})
	require.True(t, ok)
	assert.Len(t, sans, 2)
	assert.Equal(t, "lb.example.com", sans[0])
	assert.Equal(t, "10.0.0.100", sans[1])
	assert.Equal(t, true, config["tls-san-security"])
}

func TestGenerateServerConfig_KubeconfigConfig(t *testing.T) {
	k := &K3s{
		Kubeconfig: &KubeconfigConfig{
			Mode:  "644",
			Group: "admin",
			Name:  "my-cluster",
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "644", config["write-kubeconfig-mode"])
	assert.Equal(t, "admin", config["write-kubeconfig-group"])
	assert.Equal(t, "my-cluster", config["write-kubeconfig-name"])
}

func TestGenerateServerConfig_DatastoreConfig(t *testing.T) {
	k := &K3s{
		Datastore: &DatastoreConfig{
			Endpoint: "postgres://user:pass@db:5432/k3s",
			CAFile:   "/etc/ssl/ca.pem",
			CertFile: "/etc/ssl/cert.pem",
			KeyFile:  "/etc/ssl/key.pem",
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "postgres://user:pass@db:5432/k3s", config["datastore-endpoint"])
	assert.Equal(t, "/etc/ssl/ca.pem", config["datastore-cafile"])
	assert.Equal(t, "/etc/ssl/cert.pem", config["datastore-certfile"])
	assert.Equal(t, "/etc/ssl/key.pem", config["datastore-keyfile"])
}

func TestGenerateServerConfig_ServerListenerConfig(t *testing.T) {
	k := &K3s{
		Server: &ServerConfig{
			AdvertiseAddress: "10.0.0.50",
			AdvertisePort:    6443,
			BindAddress:      "0.0.0.0",
			HTTPSPort:        6443,
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "10.0.0.50", config["advertise-address"])
	assert.Equal(t, 6443, config["advertise-port"])
	assert.Equal(t, "0.0.0.0", config["bind-address"])
	assert.Equal(t, 6443, config["https-listen-port"])
}

func TestGenerateServerConfig_SecurityConfig(t *testing.T) {
	k := &K3s{
		Security: &SecurityConfig{
			SELinux:               true,
			ProtectKernelDefaults: true,
			EncryptionProvider:    "secretbox",
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, true, config["selinux"])
	assert.Equal(t, true, config["protect-kernel-defaults"])
	assert.Equal(t, "secretbox", config["secrets-encryption-provider"])
}

func TestGenerateServerConfig_ComponentArgs(t *testing.T) {
	k := &K3s{
		ComponentArgs: &ComponentArgsConfig{
			KubeAPIServer:              []string{"--audit-log-path=/var/log/audit.log"},
			KubeControllerManager:      []string{"--node-cidr-mask-size=24"},
			KubeScheduler:              []string{"--config=/etc/scheduler.yaml"},
			KubeCloudControllerManager: []string{"--cloud-provider=external"},
			Kubelet:                    []string{"--max-pods=110"},
			KubeProxy:                  []string{"--proxy-mode=ipvs"},
			Etcd:                       []string{"--quota-backend-bytes=8589934592"},
			HelmController:             []string{"--default-job-image=rancher/klipper-helm:v0.8.0"},
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	apiArgs, ok := config["kube-apiserver-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--audit-log-path=/var/log/audit.log", apiArgs[0])

	ctrlArgs, ok := config["kube-controller-manager-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--node-cidr-mask-size=24", ctrlArgs[0])

	schedArgs, ok := config["kube-scheduler-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--config=/etc/scheduler.yaml", schedArgs[0])

	ccmArgs, ok := config["kube-cloud-controller-manager-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--cloud-provider=external", ccmArgs[0])

	kubeletArgs, ok := config["kubelet-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--max-pods=110", kubeletArgs[0])

	proxyArgs, ok := config["kube-proxy-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--proxy-mode=ipvs", proxyArgs[0])

	etcdArgs, ok := config["etcd-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--quota-backend-bytes=8589934592", etcdArgs[0])

	helmArgs, ok := config["helm-controller-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--default-job-image=rancher/klipper-helm:v0.8.0", helmArgs[0])
}

func TestGenerateAgentConfig_ComponentArgs(t *testing.T) {
	k := &K3s{
		ComponentArgs: &ComponentArgsConfig{
			KubeAPIServer: []string{"--audit-log-path=/var/log/audit.log"},
			Kubelet:       []string{"--max-pods=110"},
			KubeProxy:     []string{"--proxy-mode=ipvs"},
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleAgent}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "token")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	// Agent should get kubelet and kube-proxy args, but NOT server-only args
	kubeletArgs, ok := config["kubelet-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--max-pods=110", kubeletArgs[0])

	proxyArgs, ok := config["kube-proxy-arg"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "--proxy-mode=ipvs", proxyArgs[0])

	_, hasAPI := config["kube-apiserver-arg"]
	assert.False(t, hasAPI, "agent should not have kube-apiserver-arg")
}

func TestGenerateAgentConfig_SecurityConfig(t *testing.T) {
	k := &K3s{
		Security: &SecurityConfig{
			SELinux:               true,
			ProtectKernelDefaults: true,
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleAgent}

	data, err := k.GenerateAgentConfig(host, "10.0.0.1", "token")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, true, config["selinux"])
	assert.Equal(t, true, config["protect-kernel-defaults"])
}

func TestGenerateServerConfig_RuntimeConfig(t *testing.T) {
	k := &K3s{
		Runtime: &RuntimeConfig{
			SystemDefaultRegistry: "registry.example.com",
			PauseImage:            "registry.example.com/pause:3.9",
			Snapshotter:           "stargz",
			DefaultRuntime:        "nvidia",
		},
	}
	k.SetDefaults()
	host := &Host{Role: RoleServer}

	data, err := k.GenerateServerConfig(host, "", true, "")
	require.NoError(t, err)

	var config map[string]any
	require.NoError(t, yaml.Unmarshal(data, &config))

	assert.Equal(t, "registry.example.com", config["system-default-registry"])
	assert.Equal(t, "registry.example.com/pause:3.9", config["pause-image"])
	assert.Equal(t, "stargz", config["snapshotter"])
	assert.Equal(t, "nvidia", config["default-runtime"])
}

func TestK3sValidate_NetworkConfig(t *testing.T) {
	tests := []struct {
		name    string
		k3s     K3s
		wantErr bool
		errMsg  string
	}{
		{
			"valid flannel backend",
			K3s{Network: &NetworkConfig{FlannelBackend: "vxlan"}},
			false, "",
		},
		{
			"invalid flannel backend",
			K3s{Network: &NetworkConfig{FlannelBackend: "invalid"}},
			true, "flannelBackend",
		},
		{
			"valid egress mode",
			K3s{Network: &NetworkConfig{EgressSelectorMode: "cluster"}},
			false, "",
		},
		{
			"invalid egress mode",
			K3s{Network: &NetworkConfig{EgressSelectorMode: "invalid"}},
			true, "egressSelectorMode",
		},
		{
			"nil network is valid",
			K3s{},
			false, "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.k3s.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestK3sValidate_DatastoreConfig(t *testing.T) {
	t.Run("TLS files without endpoint", func(t *testing.T) {
		k := K3s{Datastore: &DatastoreConfig{CAFile: "/ca.pem"}}
		assert.Error(t, k.Validate())
	})

	t.Run("endpoint with TLS files", func(t *testing.T) {
		k := K3s{Datastore: &DatastoreConfig{
			Endpoint: "postgres://db:5432/k3s",
			CAFile:   "/ca.pem",
		}}
		assert.NoError(t, k.Validate())
	})

	t.Run("endpoint without TLS", func(t *testing.T) {
		k := K3s{Datastore: &DatastoreConfig{Endpoint: "postgres://db:5432/k3s"}}
		assert.NoError(t, k.Validate())
	})
}

func TestK3sValidate_SecurityConfig(t *testing.T) {
	t.Run("valid provider aescbc", func(t *testing.T) {
		k := K3s{Security: &SecurityConfig{EncryptionProvider: "aescbc"}}
		assert.NoError(t, k.Validate())
	})

	t.Run("valid provider secretbox", func(t *testing.T) {
		k := K3s{Security: &SecurityConfig{EncryptionProvider: "secretbox"}}
		assert.NoError(t, k.Validate())
	})

	t.Run("invalid provider", func(t *testing.T) {
		k := K3s{Security: &SecurityConfig{EncryptionProvider: "unknown"}}
		err := k.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encryptionProvider")
	})
}

func TestK3sYAMLRoundTrip_NewFeatures(t *testing.T) {
	input := `
version: v1.31.0+k3s1
network:
  clusterCIDR:
    - 10.42.0.0/16
    - fd42::/48
  serviceCIDR:
    - 10.43.0.0/16
  clusterDNS:
    - 10.43.0.10
  clusterDomain: my.cluster
  nodePortRange: "30000-40000"
  flannelBackend: wireguard-native
  flannelIPv6Masq: true
  flannelExternalIP: true
  egressSelectorMode: cluster
  serviceLBNamespace: lb-ns
disable:
  components:
    - traefik
    - servicelb
  scheduler: true
  kubeProxy: true
  helmController: true
tls:
  sans:
    - lb.example.com
    - 10.0.0.100
  sanSecurity: true
kubeconfig:
  mode: "644"
  group: admin
  name: my-cluster
datastore:
  endpoint: "postgres://db:5432/k3s"
  caFile: /ca.pem
  certFile: /cert.pem
  keyFile: /key.pem
server:
  advertiseAddress: 10.0.0.50
  advertisePort: 6443
  bindAddress: 0.0.0.0
  httpsPort: 6443
security:
  selinux: true
  protectKernelDefaults: true
  encryptionProvider: secretbox
componentArgs:
  kubeAPIServer:
    - "--audit-log-path=/var/log/audit.log"
  kubelet:
    - "--max-pods=110"
  etcd:
    - "--quota-backend-bytes=8589934592"
runtime:
  systemDefaultRegistry: registry.example.com
  pauseImage: registry.example.com/pause:3.9
  snapshotter: stargz
  defaultRuntime: nvidia
`
	var k K3s
	require.NoError(t, yaml.Unmarshal([]byte(input), &k))

	// Network
	require.NotNil(t, k.Network)
	assert.Equal(t, []string{"10.42.0.0/16", "fd42::/48"}, k.Network.ClusterCIDR)
	assert.Equal(t, []string{"10.43.0.0/16"}, k.Network.ServiceCIDR)
	assert.Equal(t, "my.cluster", k.Network.ClusterDomain)
	assert.Equal(t, "wireguard-native", k.Network.FlannelBackend)
	assert.True(t, k.Network.FlannelIPv6Masq)
	assert.Equal(t, "cluster", k.Network.EgressSelectorMode)

	// Disable
	require.NotNil(t, k.Disable)
	assert.Equal(t, []string{"traefik", "servicelb"}, k.Disable.Components)
	assert.True(t, k.Disable.Scheduler)
	assert.True(t, k.Disable.KubeProxy)
	assert.True(t, k.Disable.HelmController)

	// TLS
	require.NotNil(t, k.TLS)
	assert.Equal(t, []string{"lb.example.com", "10.0.0.100"}, k.TLS.SANs)
	assert.True(t, k.TLS.SANSecurity)

	// Kubeconfig
	require.NotNil(t, k.Kubeconfig)
	assert.Equal(t, "644", k.Kubeconfig.Mode)
	assert.Equal(t, "admin", k.Kubeconfig.Group)
	assert.Equal(t, "my-cluster", k.Kubeconfig.Name)

	// Datastore
	require.NotNil(t, k.Datastore)
	assert.Equal(t, "postgres://db:5432/k3s", k.Datastore.Endpoint)
	assert.Equal(t, "/ca.pem", k.Datastore.CAFile)

	// Server
	require.NotNil(t, k.Server)
	assert.Equal(t, "10.0.0.50", k.Server.AdvertiseAddress)
	assert.Equal(t, 6443, k.Server.AdvertisePort)

	// Security
	require.NotNil(t, k.Security)
	assert.True(t, k.Security.SELinux)
	assert.True(t, k.Security.ProtectKernelDefaults)
	assert.Equal(t, "secretbox", k.Security.EncryptionProvider)

	// ComponentArgs
	require.NotNil(t, k.ComponentArgs)
	assert.Equal(t, []string{"--audit-log-path=/var/log/audit.log"}, k.ComponentArgs.KubeAPIServer)
	assert.Equal(t, []string{"--max-pods=110"}, k.ComponentArgs.Kubelet)
	assert.Equal(t, []string{"--quota-backend-bytes=8589934592"}, k.ComponentArgs.Etcd)

	// Runtime
	require.NotNil(t, k.Runtime)
	assert.Equal(t, "registry.example.com", k.Runtime.SystemDefaultRegistry)
	assert.Equal(t, "registry.example.com/pause:3.9", k.Runtime.PauseImage)
	assert.Equal(t, "stargz", k.Runtime.Snapshotter)
	assert.Equal(t, "nvidia", k.Runtime.DefaultRuntime)
}
