---
name: devbox-go-development
description: Use when changing DevBox Go code, CLI behavior, project/session/route state, Docker orchestration, Caddy routing, or agent-facing JSON output.
---

# DevBox Go Development

## Core Rules

- Keep Cobra wiring in `cmd/devbox`.
- Keep domain structs in `internal/devbox/model`.
- Keep JSON persistence in `internal/devbox/state`.
- Keep host/runtime operations in `internal/devbox/service`.
- Keep Caddy config generation in `internal/devbox/proxy`.
- Write tests before behavior changes. State, URL generation, and JSON responses should be covered without Docker.
- Keep host effects behind small functions or methods so Docker, Caddy, tmux, and Tailscale failures are easy to report.
- `--json` output must be a single stable JSON object with `ok: true` on success.
- Resource names are lowercase letters, numbers, and hyphens.
- Hand-written code files must stay under 400 LoC; run `make check` to enforce this.

## Verification

Run these before reporting completion:

```sh
make check
```

For a focused loop:

```sh
go test ./...
go build ./cmd/devbox
```

## Architecture Preferences

- JSON state is acceptable for v1; do not introduce a database unless a feature needs transactions.
- Routes are explicit and agent-registered; do not add port autodetection.
- Project containers are sibling containers managed through the host Docker socket.
- Zed support is host SSH plus helper links, not a separate remote-server manager.
