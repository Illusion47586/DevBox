# DevBox Agent Guide

DevBox is a Go CLI/control-plane project. Keep the CLI deterministic and agent-friendly: every command that changes or reports state must support `--json`, return useful errors, and avoid hidden interactive prompts unless the command is explicitly an attach/shell command.

## Development Workflow

- Use `make check` before claiming work is complete.
- Use `go test ./...` for the fast verification loop.
- Keep Cobra command wiring in `cmd/devbox`.
- Keep domain structs in `internal/devbox/model`.
- Keep JSON persistence in `internal/devbox/state`.
- Keep host/runtime operations in `internal/devbox/service`.
- Keep reverse-proxy config generation in `internal/devbox/proxy`.
- Prefer small interfaces around host effects. State logic should stay testable without Docker, Caddy, Tailscale, or tmux running.
- Use JSON files for v1 state unless a feature clearly requires transactional storage.

## Architecture Rules

- DevBox runs as a control-plane container and manages sibling project containers through the host Docker socket.
- Project commands and agent tmux sessions run inside per-project containers.
- Zed connects to TheBox host SSH; DevBox provides wrapper commands and URLs, not a custom Zed protocol.
- Routes are explicit. Do not auto-detect listening ports in v1.
- Private network only. Do not add public exposure or managed HTTPS without a new design decision.

## Command Design

- Prefer `devbox noun verb` command groups.
- Use lowercase resource names with digits and hyphens.
- On `--json`, write exactly one JSON object to stdout.
- Human output may be concise text, but should never be required by agents.
- Long-running interactive commands are limited to `shell` and `session attach`.
- Hand-written code files must stay under 400 LoC. `make check` enforces this through `scripts/check-loc.sh`.
