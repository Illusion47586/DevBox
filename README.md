# DevBox

DevBox is an agent-first development control plane for a Linux home server called TheBox. It is designed to run as a Docker container on ZimaOS, manage per-project development containers through the host Docker socket, and expose project services through private LAN/Tailscale URLs.

## Status

This repository contains the v1 implementation scaffold and core CLI behavior. It is intentionally headless: no dashboard, no public routes, no managed HTTPS, and no automatic dev-server discovery.

## Development

```sh
make check
make docker-build
```

Useful environment defaults are documented in `.env.example`.

DevBox publishes one image:

```text
ghcr.io/illusion47586/devbox:latest
```

The same image runs the control plane and acts as the base image for per-project containers. Project containers override the image entrypoint and stay alive with `sleep infinity`.

The image also installs DevBox skills and selected registry skills from `skills.sh` globally for Codex-style, Claude Code-style, and generic agent skill directories:

```text
/root/.codex/skills
/root/.agents/skills
/root/.claude/skills
```

Registry skills installed by default:

- `vercel-labs/agent-skills`: `vercel-react-best-practices`, `web-design-guidelines`, `writing-guidelines`
- `obra/superpowers`: `test-driven-development`, `systematic-debugging`, `requesting-code-review`, `receiving-code-review`, `verification-before-completion`, `writing-plans`

## Install On ZimaOS

Use the published GHCR image:

```sh
mkdir -p /DATA/AppData/devbox
cd /DATA/AppData/devbox
curl -fsSLO https://raw.githubusercontent.com/Illusion47586/DevBox/main/deploy/compose.yaml
curl -fsSLO https://raw.githubusercontent.com/Illusion47586/DevBox/main/deploy/Caddyfile
docker compose -f compose.yaml up -d
docker compose -f compose.yaml exec devbox devbox doctor --json
```

Configure your LAN DNS so `*.devbox` resolves to TheBox. Then agents can run commands through the running control-plane container:

```sh
docker compose -f compose.yaml exec devbox devbox project register /workspaces/myapp --name myapp --json
docker compose -f compose.yaml exec devbox devbox route add --project myapp --service web --target 3000 --json
```

## Core Commands

```sh
devbox project create https://github.com/example/myapp.git --json
devbox project register /workspaces/myapp --name myapp --json
devbox session create myapp --name agent-main -- codex exec
devbox route add --project myapp --service web --target 3000 --json
devbox zed myapp --json
devbox doctor --json
```

Canonical service URLs use:

```text
http://service.project.devbox
```

Your LAN DNS should resolve `*.devbox` to TheBox. Tailscale/MagicDNS support is private-network-first and verified through `devbox doctor`.
