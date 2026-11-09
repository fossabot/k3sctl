# k3sctl Development Roadmap

Architecture documentation and implementation plan for k3sctl.

## Overview

| File | Contents |
|------|----------|
| [00-OVERVIEW.md](00-OVERVIEW.md) | Project overview, K3s specifics, dependencies, file counts |

## Detailed Analysis

| File | Contents |
|------|----------|
| [01-ENTRYPOINT-AND-CLI.md](01-ENTRYPOINT-AND-CLI.md) | Entry point, CLI layer: main.go, cmd/ |
| [02-ACTIONS.md](02-ACTIONS.md) | Action layer: apply, reset, backup, kubeconfig |
| [03-PHASE-MANAGER.md](03-PHASE-MANAGER.md) | Phase manager framework |
| [04-PHASES-INFRASTRUCTURE.md](04-PHASES-INFRASTRUCTURE.md) | Reusable phases: connect, detect_os, facts, hooks |
| [05-PHASES-K3S-SPECIFIC.md](05-PHASES-K3S-SPECIFIC.md) | K3s phases: install, upgrade, reset, backup |
| [06-CONFIGURER.md](06-CONFIGURER.md) | Configurer layer: interface, registry, Linux distros |
| [07-API-TYPES.md](07-API-TYPES.md) | Config data model: Cluster, Host, K3s, Spec |
| [08-PACKAGES.md](08-PACKAGES.md) | Supporting packages: retry, download, binprovider, node |

## Implementation Plan

| File | Contents |
|------|----------|
| [09-FILE-MAPPING.md](09-FILE-MAPPING.md) | Complete file inventory |
| [10-IMPLEMENTATION-PLAN.md](10-IMPLEMENTATION-PLAN.md) | Detailed implementation plan with task checklist |

## Quick Stats

- **~109 Go files** (82 source + 27 test)
- **~40 direct implementations** (infrastructure/utility code)
- **~60 K3s-specific implementations** (K3s lifecycle logic)
- **~3 additional files** (K3s-specific additions)
