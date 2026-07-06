---
name: devbox-operations
description: Use when operating DevBox on TheBox, configuring Docker/ZimaOS installs, registering projects, managing tmux sessions, exposing routes, or diagnosing LAN/Tailscale access.
---

# DevBox Operations

## Install Shape

DevBox runs as a Docker control-plane container with:

- `/var/run/docker.sock` mounted from TheBox.
- A persistent state volume mounted at `/var/lib/devbox`.
- A host workspace directory mounted at `/workspaces`.
- Caddy available for private HTTP routing.

## Common Commands

```sh
devbox doctor --json
devbox project register /workspaces/myapp --name myapp --json
devbox session create myapp --name agent-main -- codex exec
devbox route add --project myapp --service web --target 3000 --json
devbox route list --json
```

## Network Expectations

- LAN DNS should resolve `*.devbox` to TheBox.
- Tailscale should make TheBox reachable from the tailnet.
- DevBox v1 exposes private HTTP routes only; do not configure public ingress as part of normal operation.
