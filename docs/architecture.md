# DevBox Architecture

## Runtime Model

DevBox is a Go CLI packaged into a Docker image. The DevBox container mounts the host Docker socket, a persistent state volume, and a host workspace directory. It creates sibling project containers from the same image rather than running projects inside the control-plane container.

## Components

- `cmd/devbox`: Cobra command tree and JSON output shaping.
- `internal/devbox/model`: shared domain types.
- `internal/devbox/state`: JSON persistence and state-level validation.
- `internal/devbox/service`: project, route, session, Zed, and doctor operations.
- `internal/devbox/proxy`: Caddy config generation.
- `Dockerfile`: single DevBox image for the control plane and project containers.
- `skills/`: repo-local skills that guide future coding agents.

The image installs selected `skills.sh` registry skills with the `skills` CLI, then installs local `skills/` into `/root/.codex/skills`, `/root/.agents/skills`, and `/root/.claude/skills` so bundled agents can discover DevBox-specific operating guidance from inside project containers.

## State

State is stored as JSON in `DEVBOX_STATE_PATH`. The v1 state tracks projects, routes, and tmux sessions. This keeps the system inspectable and easy to back up on a home server.

## Networking

Caddy handles dynamic HTTP reverse proxying. Agents explicitly register routes through the CLI. DevBox generates hosts in the form `service.project.devbox`; existing LAN DNS is responsible for resolving `*.devbox` to TheBox.

## Sessions

Agent sessions are tmux sessions inside project containers. DevBox records session metadata and provides attach/stop/list commands so agents do not need to remember Docker container names.
