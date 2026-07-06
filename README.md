# DevBox

DevBox is an agent-first development control plane for a Linux home server called TheBox. It is designed to run as a Docker container on ZimaOS, manage per-project development containers through the host Docker socket, and expose project services through private LAN/Tailscale URLs.

## Status

This repository contains the v1 implementation scaffold and core CLI behavior. It is intentionally headless: no dashboard, no public routes, no managed HTTPS, and no automatic dev-server discovery.

## Development

```sh
make check
make docker-build
make docker-build-project
```

Useful environment defaults are documented in `.env.example`.

## Install On ZimaOS

Use the published GHCR images:

```sh
mkdir -p /DATA/AppData/devbox
cd /DATA/AppData/devbox
curl -fsSLO https://raw.githubusercontent.com/Illusion47586/DevBox/main/deploy/zimaos-compose.yaml
curl -fsSLO https://raw.githubusercontent.com/Illusion47586/DevBox/main/deploy/Caddyfile
docker compose -f zimaos-compose.yaml up -d
docker compose -f zimaos-compose.yaml exec devbox devbox doctor --json
```

Configure your LAN DNS so `*.devbox` resolves to TheBox. Then agents can run commands through the running control-plane container:

```sh
docker compose -f zimaos-compose.yaml exec devbox devbox project register /workspaces/myapp --name myapp --json
docker compose -f zimaos-compose.yaml exec devbox devbox route add --project myapp --service web --target 3000 --json
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
