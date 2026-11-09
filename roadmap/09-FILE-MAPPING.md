# Roadmap: Complete File Inventory

## Legend
- **CORE**: Core infrastructure, generic utilities
- **K3S**: K3s-specific logic and paths
- **NEW**: New file for K3s-specific functionality

## Entry Point & Build
| File | Category |
|------|----------|
| main.go | CORE |
| version/version.go | CORE |
| go.mod | NEW (k3sctl module path) |

## CLI Layer (cmd/)
| File | Category |
|------|----------|
| cmd/root.go | CORE |
| cmd/flags.go | K3S (k3sctl config format) |
| cmd/apply.go | CORE |
| cmd/init.go | K3S (K3s config template) |
| cmd/reset.go | CORE |
| cmd/backup.go | CORE |
| cmd/kubeconfig.go | CORE |
| cmd/config_edit.go | K3S (file-based config) |
| cmd/config_status.go | K3S (file-based config) |
| cmd/version.go | K3S (k3sctl releases) |
| cmd/completion.go | CORE |
| cmd/logo.go | NEW (k3sctl logo) |
| cmd/flags_test.go | CORE |
| cmd/init_test.go | K3S |
| cmd/apply_test.go | CORE |

## Action Layer (action/)
| File | Category |
|------|----------|
| action/apply.go | K3S (K3s phase sequence) |
| action/reset.go | K3S (K3s reset) |
| action/backup.go | K3S (K3s backup) |
| action/kubeconfig.go | CORE |
| action/config_edit.go | K3S |
| action/config_status.go | K3S |

## Phase Layer (phase/)
| File | Category |
|------|----------|
| phase/manager.go | CORE |
| phase/manager_test.go | CORE |
| phase/generic_phase.go | CORE |
| phase/connect.go | CORE |
| phase/connect_test.go | CORE |
| phase/disconnect.go | CORE |
| phase/detect_os.go | CORE |
| phase/lock.go | K3S (k3sctl lock path) |
| phase/unlock.go | K3S |
| phase/prepare_hosts.go | K3S |
| phase/gather_facts.go | K3S |
| phase/gather_facts_test.go | K3S |
| phase/validate_hosts.go | K3S |
| phase/validate_hosts_test.go | K3S |
| phase/validate_facts.go | K3S |
| phase/validate_facts_test.go | K3S |
| phase/run_hooks.go | CORE |
| phase/daemon_reload.go | CORE |
| phase/downloads.go | CORE |
| phase/uploadfiles.go | CORE |
| phase/uploadfiles_test.go | CORE |
| phase/default_k3s_version.go | K3S (K3s version resolution) |
| phase/gather_k3s_facts.go | K3S (K3s fact gathering) |
| phase/gather_k3s_facts_test.go | K3S |
| phase/stage_binaries.go | CORE |
| phase/stage_binaries_test.go | CORE |
| phase/upload_k3s_test.go | K3S |
| phase/install_binaries.go | K3S |
| phase/configure_k3s.go | K3S (K3s config format) |
| phase/configure_k3s_test.go | K3S |
| phase/initialize_k3s.go | K3S (K3s server init) |
| phase/install_servers.go | K3S (K3s server join) |
| phase/install_agents.go | K3S (K3s agent join) |
| phase/upgrade_servers.go | K3S |
| phase/upgrade_agents.go | K3S |
| phase/reinstall.go | K3S |
| phase/reset_agents.go | K3S |
| phase/reset_servers.go | K3S |
| phase/reset_leader.go | K3S |
| phase/get_kubeconfig.go | K3S (read k3s.yaml) |
| phase/get_kubeconfig_test.go | K3S |
| phase/backup.go | K3S (K3s snapshot) |
| phase/restore.go | K3S (K3s restore) |
| phase/apply_manifests.go | K3S (k3s kubectl) |
| phase/validate_etcd_members.go | K3S (conditional) |
| phase/arm_prepare.go | K3S (cgroup check) |

## Configurer Layer (configurer/)
| File | Category |
|------|----------|
| configurer/interface.go | K3S (K3s methods) |
| configurer/host.go | CORE |
| configurer/registry.go | CORE |
| configurer/linux.go | K3S (K3s paths) |
| configurer/linux/alpine.go | K3S |
| configurer/linux/alpine_test.go | K3S |
| configurer/linux/debian.go | K3S |
| configurer/linux/ubuntu.go | K3S |
| configurer/linux/enterpriselinux.go | K3S |
| configurer/linux/enterpriselinux/rhel.go | K3S |
| configurer/linux/enterpriselinux/centos.go | K3S |
| configurer/linux/enterpriselinux/fedora.go | K3S |
| configurer/linux/enterpriselinux/rocky.go | K3S |
| configurer/linux/enterpriselinux/almalinux.go | K3S |
| configurer/linux/enterpriselinux/oracle.go | K3S |
| configurer/linux/enterpriselinux/amazon.go | K3S |
| configurer/linux/archlinux.go | K3S |
| configurer/linux/opensuse.go | K3S |
| configurer/linux/sles.go | K3S |
| configurer/linux/flatcar.go | K3S |
| configurer/linux/coreos.go | K3S |
| configurer/linux/linux_test.go | K3S |

## API Types (pkg/apis/)
| File | Category |
|------|----------|
| pkg/apis/.../cluster.go | K3S (K3s API version) |
| pkg/apis/.../cluster_test.go | K3S |
| .../cluster/spec.go | K3S (K3s spec) |
| .../cluster/spec_test.go | K3S |
| .../cluster/host.go | K3S (K3s host methods) |
| .../cluster/host_test.go | K3S |
| .../cluster/host_posix_test.go | K3S |
| .../cluster/hosts.go | K3S (Servers/Agents) |
| .../cluster/hosts_test.go | K3S |
| .../cluster/k3s.go | K3S (K3s config) |
| .../cluster/k3s_test.go | K3S |
| .../cluster/flags.go | CORE |
| .../cluster/flags_test.go | CORE |
| .../cluster/flags_quoting_test.go | CORE |
| .../cluster/options.go | CORE |
| .../cluster/hook.go | CORE |
| .../cluster/uploadfile.go | CORE |
| .../cluster/uploadfile_test.go | CORE |
| .../cluster/binary_providers.go | K3S |
| .../cluster/quote.go | CORE |

## Supporting Packages (pkg/)
| File | Category |
|------|----------|
| pkg/retry/retry.go | CORE |
| pkg/retry/retry_test.go | CORE |
| pkg/download/download.go | K3S (cache path) |
| pkg/download/download_test.go | K3S |
| pkg/k3s/provider.go | CORE |
| pkg/k3s/binprovider/host.go | K3S |
| pkg/k3s/binprovider/helpers.go | CORE |
| pkg/k3s/binprovider/download.go | K3S (K3s URLs) |
| pkg/k3s/binprovider/local_upload.go | K3S (K3s cache) |
| pkg/k3s/binprovider/local_file.go | CORE |
| pkg/k3s/binprovider/existing.go | CORE |
| pkg/manifest/reader.go | CORE |
| pkg/manifest/reader_test.go | CORE |
| pkg/node/statusfunc.go | K3S |
| pkg/node/statusfunc_test.go | K3S |

## Integration
| File | Category |
|------|----------|
| integration/github/github.go | K3S (K3s releases) |

## Summary Count
- **CORE** (generic infrastructure): ~40 files
- **K3S** (K3s-specific logic): ~60 files
- **NEW** (brand new): ~3 files
- **Total**: ~103 files
