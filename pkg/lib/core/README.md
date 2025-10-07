# pkg/lib/core

Purpose
- Shared runtime primitives consumed by generated projects.
- Keep it small and boring: config, logging hooks, HTTP helpers, validation, composite FS utilities.

Guidelines
- No magic, no reflection-heavy tricks.
- Prefer explicit funcs over global state.
- Avoid coupling to specific web/db frameworks.
