# Aquamarine: A Feature‑First Monolith Generator for Go

Status: Draft | Updated: 2025-10-07 | Version: 0.1

## The Idea

Aquamarine is a generator plus a small runtime library that produces feature‑first Go monoliths. You describe your app in YAML (and/or via CLI commands), and get a single executable that contains code and embedded assets (migrations, queries, seeds, templates, static files).

## Why

We want spinning up a cohesive monolith to be fast, explicit, and ergonomic:

- Feature-first layout that’s easy to navigate and remove wholesale.
- Minimal ceremony, small interfaces, no hidden magic.
- A single binary to run the entire app.


## The Philosophy

### Less Framework, More Generated Code
Generate explicit Go you can read, edit, and own. Keep the runtime thin and predictable.

### Feature‑First, Not Folders‑by‑Type
Each feature lives in a single package with its domain types, service, and handlers side‑by‑side.

### Single Binary with Embedded Assets
All assets are bundled via go:embed and loaded from a composite FS at runtime.

### In‑Process Composition by Default
Web (HTML/HTMX) and API layers call feature services directly in‑process. HTTP is optional for cross‑process scenarios.

### Small, Purposeful Interfaces
No generic or bloated interfaces. Prefer concrete structs for models and keep service/repo interfaces tight.  
In simple atomic features (single-entity CRUDs), the repository implementation itself can satisfy the service interface and be injected into the handler, ensuring consistency with more complex domains.


## Feature Model
Aquamarine structures a monolith as a set of features:

- **feature** (default): a business module with domain types, a small service, local API handlers, and assets.  
- **web**: a web layer that orchestrates multiple features to render a cohesive HTML interface (optionally with HTMX). In the future, a separate BFF may also be implemented to aggregate responses from different services and expose them as a consolidated API for external clients.  
- **support** (optional): jobs, schedulers, subscribers — same packaging rules, if needed.


## The DSL: aquamarine.yaml

Everything can be described declaratively and/or created with CLI. Example:

```yaml path=null start=null
version: 0.1
project:
  name: myapp
  module: github.com/example/myapp
runtime:
  http:
    api: { host: 127.0.0.1, port: 8081 }
    web: { host: 127.0.0.1, port: 8080 }
  database:
    engine: sqlite
feats:
  - name: auth
    service:
      methods: [Login, Register]
    api:
      routes:
        - {method: POST, path: "/login", handler: Login}
        - {method: POST, path: "/register", handler: Register}
  - name: web
    kind: web
    pages:
      - {route: "GET /dashboard", uses: [auth]}
```

## What It Generates

### File Structure

```text path=null start=null
.
├── main.go                 # assembly (early stages; can move to cmd/ later)
├── assets/                 # centralized assets tree
│   ├── migrations/
│   │   └── postgres/
│   │       └── 0001_auth_init.sql
│   ├── web/
│   │   └── templates/
│   │       └── auth/
│   │           └── login.html
│   └── static/
│       └── css/
│           └── app.css
├── internal/
│   ├── platform/           # config, router, db, log, errors, metrics, composite FS
│   ├── web/                # BFF/handlers/views (templates global to web)
│   └── feat/
│       └── auth/           # feature package (flat)
│           ├── feature.go  # domain structs, errors
│           ├── service.go  # interface + concrete implementation
│           ├── repo.go     # minimal repo interface (+in‑memory fake optional)
│           └── http.go     # JSON handlers + RegisterAPIRoutes
├── docs/
│   └── drafts/overview.md
└── go.mod
```

### Clean and Consistent Code
- Feature‑local API handlers in http.go with predictable registration.
- Small services; concrete structs for models; no empty interfaces.
- In‑memory fakes for repos/services to enable fast tests and CLI usage.
- Web (and BFF) composes view models using services in‑process; templates loaded from centralized assets.
- Assets organized by type in centralized `assets/` tree with feature-based subdirectories where appropriate.

## Design Decisions

- Feature‑first flat packages: no nested "api/adapters/repo" folders per feature.
- In‑process service calls for API and Web; optional HTTP providers if you need process boundaries.
- Centralized assets tree embedded by main and accessed by features via conventional paths.
- A generated "sync" file aggregates features: registers routes, exposes services to Web.
- Minimal runtime primitives: config, router/middleware, db connectors, logging, validation, embedded FS.

## Assets Strategy

Aquamarine uses a centralized `assets/` tree at the repository root. This tree is embedded by main via go:embed and features access assets through conventional paths.

**Key principles:**
- Clean separation between code (in `internal/feat/`) and assets (in `assets/`)
- Assets organized by type: migrations, web templates, static files
- Feature-specific assets use subdirectories within their asset type
- Zero "sync" step needed when modifying assets—go:embed handles changes on rebuild

**Asset Organization:**
```
assets/
├── migrations/           # database migrations by engine
│   └── postgres/
│       └── 0001_auth_init.sql
├── web/                 # web-related assets
│   └── templates/
│       └── auth/           # feature-specific templates
│           └── login.html
└── static/              # static files (CSS, JS, images)
    └── css/
        └── app.css
```

**Migrations and Seeds:**
- Ordering: filename‑based (timestamp or incremental). Within an engine, apply lexicographically; across features: default alphabetical, with optional requires for topological ordering.
- Integrity: checksum tracking for applied migrations; edits to past files fail with a clear message (create a new migration instead).
- Seeding: prefer idempotent UPSERTs or code‑based seeds via in‑process services; optional phases (10‑system, 20‑foundation, 30‑feature, 90‑demo); cross‑feature dependencies via requires or stable natural keys.

**Developer Experience:**
- Asset changes are automatically picked up by go:embed on rebuild
- Adding/removing features requires no asset synchronization step
- Templates and static files follow predictable paths for easy discovery

---

**TL;DR**: Aquamarine generates feature‑first Go monoliths from YAML/CLI, with embedded assets and in‑process composition. One binary runs the app.
