# Aquamarine Generator - Commands

Status: Draft | Updated: 2025-10-07 | Version: 0.1

This document sketches the generator CLI and core flows. Each section is intentionally lightweight and marked with [TODO] to evolve iteratively as needs arise.

## Overview

- Scope: feature-first monolith generator + thin runtime primitives
- Config: aquamarine.yaml (centralized assets; convention over configuration)
- Outputs: project scaffolding, feature packages, API/web wiring, optional aggregator
- Non-goals: heavy frameworks, deep DI, internal HTTP between layers

## Command reference

### new [TODO]
- Purpose: scaffold a fresh repo structure.
- Inputs: project.name, project.module; defaults from current directory if missing.
- Behavior [TODO]:
  - Create base dirs: internal/feat/, internal/web/, internal/platform/, assets/
  - Add initial files: README, go.mod (optional), basic router/server stubs [TBD]
  - Write starter aquamarine.yaml if not present
- Flags [TODO]:
  - --module <path>
  - --skip-go-mod
- Open questions [TODO]:
  - What minimal platform stubs to include by default?

### add feature <name> [TODO]
- Purpose: add a new feature package (flat, feature-first).
- Behavior [TODO]:
  - Create internal/feat/<name> with files: feature.go, service.go, repo.go, http.go, *_test.go (optional)
  - Do not create assets entries in YAML (discovered by convention)
- Flags [TODO]:
  - --kind <atom|domain> (default: domain)
  - --models <comma-separated> (domain)
- YAML impact [TODO]:
  - Append to feats: [{ name, kind, ... }]
- Open questions [TODO]:
  - Should atom scaffold collapse repo==service by default?

### add endpoint <feat> <METHOD> <path> [TODO]
- Purpose: add an API route mapped to a feature service method.
- Behavior [TODO]:
  - Ensure service method exists; create stub if missing
  - Append handler in internal/feat/<feat>/http.go
  - Add to YAML feats[feat].api.routes
- Flags [TODO]:
  - --method-name <Name>
  - --body <json-schema|fields?> [TBD]
- Open questions [TODO]:
  - Where to define validation rules for requests? (model/DTO annotations vs functions)

### generate -f aquamarine.yaml [TODO]
- Purpose: idempotently sync filesystem to YAML (create or update scaffolding).
- Behavior [TODO]:
  - Parse YAML; diff existing code; create missing pieces without clobbering user edits
  - Respect “user land” markers if needed [TBD]
- Flags [TODO]:
  - --dry-run
  - --force (dangerous; avoid by default)
- Open questions [TODO]:
  - Strategy for safe updates without overwriting custom code

### sync [TODO]
- Purpose: regenerate aggregator wiring (imports, route registration, service exposure) when feats are added/removed.
- Behavior [TODO]:
  - Scan internal/feat/* packages
  - Generate a single aggregator file (e.g., internal/feat/_aggregate/enable.go)
  - No need to touch assets (centralized)
- Flags [TODO]:
  - --verify (CI mode)
- Open questions [TODO]:
  - Exact shape of Provide()/registration contracts

## Runtime (app binary) - Related tasks (not generator)

### serve [TODO]
- Starts HTTP servers for web (HTML/HTMX) and API (JSON) using configured hosts/ports.

### migrate [TODO]
- Applies database migrations by engine (sqlite|mongodb initially).
- Ordering by filename (timestamp/incremental). Checksums in prod [TBD].

### seed [TODO]
- Applies idempotent seeds by engine, optionally phased.
- Cross-feature ordering via requires (optional) [TBD].

## Conventions (summary)
- Assets base dir: assets
- Templates resolved by convention in Go (no template path in YAML)
- Features are flat packages; handlers + service + domain in one place
- Web layer is server-side composition using feature services in-process

## Open design topics [TODO]
- Validation DSL: minimal set of validators and codegen strategy
- Test scaffolding defaults (service tests, handler httptest)
- Optional HTTP providers that implement service interfaces (later)
- Aggregator strategy and code ownership boundaries
