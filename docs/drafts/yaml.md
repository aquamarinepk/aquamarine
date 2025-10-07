# Aquamarine YAML — Declarative Spec (Draft)

Status: Draft | Updated: 2025-10-07 | Version: 0.1

## Purpose

Define a simple YAML for generating a feature‑first Go monolith with a single executable and centralized assets.

## File

- Name: aquamarine.yaml (repo root)
- Scope: one app (monolith), multiple features under internal/feat/<feature>

## Top‑Level Structure

```yaml path=null start=null
version: <semver>
project:
  name: <string>
  module: <go-module>
runtime:
  http:
    api:
      host: <0.0.0.0|127.0.0.1>
      port: <int>
    web:
      host: <0.0.0.0|127.0.0.1>
      port: <int>
  database:
    engine: <sqlite|mongodb>   # early engines; postgres later
    # dsn: ${ENV_VAR}          # optional
ordering:              # optional (for migrations/seeds across feats)
  requires:            # optional edges for feat ordering
    # - [featA, featB]  # featA depends on featB
feats:
  - { ...feat spec... }
```

Notes:
- Convention over configuration: no need to list asset paths; the generator/runtime look under assets/ with well-known conventions (templates/<feat>/..., migrations/<engine>/<feat>/..., seeds/<engine>/<feat>/...).
- Base dir is fixed to assets (invariant, not configurable).
- ordering is optional; default ordering is alphabetical by feat name.

## Feat

Minimal shape common to all kinds:

```yaml path=null start=null
name: <string>        # feat package name
kind: <atom|domain|web>
```

Kinds and fields:
- atom: single entity use case; repo can equal service implementation behind a small interface.
  - model: single entity (with fields/validations)
  - api.routes: list of REST endpoints
- domain: multiple entities and business rules, small service coordinating a repo
  - models: map of entities with fields/validations
  - service.methods: list of entrypoints (use cases)
  - api.routes: list of REST endpoints
- web: centralized BFF; composes other feats to render HTML/HTMX
  - pages: list of {route, uses: [feats]}

### Field details

- models/model
  - fields: map of fieldName -> {type: string, default?: any, validations?: [..]}
    - validations: [required, min, max, pattern, email, unique, ...] (subset pragmatic)
- service.methods: list of method names to scaffold (transport‑agnostic signatures will be derived)
- api.routes: list of {method: GET|POST|PUT|PATCH|DELETE, path: /path, handler: MethodName}

Notes:
- Assets discovery is by convention under assets/: no need to list migrations/seeds/templates in the YAML.
- Templates are resolved in Go by convention (no templates in YAML): handlers/web decide which template to render.
- Seeding files are timestamp‑ordered; you don’t declare them here.

## Example (full)

```yaml path=null start=null
version: 0.1

project:
  name: myapp
  module: github.com/you/myapp

runtime:
  http:
    api:
      host: 127.0.0.1
      port: 8081
    web:
      host: 127.0.0.1
      port: 8080
  database:
    engine: sqlite

ordering:
  requires:
    # - [dashboard, auth]

feats:
  - name: auth
    kind: domain
    models:
      User:
        fields:
          id:    {type: uuid}
          email: {type: string, validations: [required, email]}
          pass:  {type: string, validations: [required, min: 8]}
      Role: {fields: {name: {type: string, validations: [required, unique]}}}
      Permission: {fields: {code: {type: string, validations: [required, unique]}}}
    service:
      methods: [Register, Login, AssignRole]
    api:
      routes:
        - {method: POST, path: /login, handler: Login}
        - {method: POST, path: /register, handler: Register}

  - name: profile
    kind: atom
    model:
      name: Profile
      fields:
        id:   {type: uuid}
        name: {type: string, validations: [required]}
    api:
      routes:
        - {method: GET, path: /me, handler: GetMe}

  - name: web
    kind: web
    pages:
      - {route: "GET /", uses: [auth, profile]}
      - {route: "GET /dashboard", uses: [auth]}
```

## Conventions

- Naming: feats are lower_snake or simple lowercase; handlers match service methods when relevant.
- Assets: discovered by convention under assets/ without listing them in YAML.
- Migrations: filename‑ordered (timestamp or incremental) per engine under assets/migrations/<engine>/<feat>/.
- Seeds: timestamp‑ordered under assets/seeds/<engine>/<feat>/; prefer idempotent operations (UPSERT by natural keys).

## Validation (informal, initial)

- version: required; semver‑ish string.
- project.name/module: required.
- runtime.http.api.port, runtime.http.web.port: required; int.
- runtime.http.api.host, runtime.http.web.host: optional; defaults to 127.0.0.1 (use 0.0.0.0 to expose).
- runtime.database.engine: one of sqlite|mongodb (postgres later).
- feats[].name: required; unique.
- feats[].kind: one of atom|domain|web.
- api.routes[].method: one of GET|POST|PUT|PATCH|DELETE.

## CLI Mapping (reference)

- aquamarine generate -f aquamarine.yaml  # create/update skeleton
- aquamarine add feature <name>            # optionally, create entry and base files
- aquamarine add endpoint <feature> ...    # add routes/methods

## Notes

- The YAML is declarative and aims to be minimal; omitted fields use sensible defaults.
- The runtime loads assets centrally and features read by well‑known prefixes.
