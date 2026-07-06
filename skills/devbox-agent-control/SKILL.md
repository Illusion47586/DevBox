---
name: devbox-agent-control
description: Use when a coding agent needs to control DevBox projects, tmux sessions, service routes, Zed links, or JSON CLI workflows from inside the DevBox image.
---

# DevBox Agent Control

## Rules

- Use `devbox ... --json` for state-changing or inspection commands.
- Register routes explicitly; do not infer ports from running processes.
- Create project sessions with `devbox session create <project> --name <name> -- <command>`.
- Keep long-running work inside project tmux sessions.
- Use `devbox shell <project>` for interactive project access.
- Use `devbox zed <project> --json` when a human needs a Zed SSH link.

## Common Flow

```sh
devbox project register /workspaces/myapp --name myapp --json
devbox session create myapp --name agent-main -- codex exec
devbox route add --project myapp --service web --target 3000 --json
devbox route list --json
```

## Network Expectations

DevBox v1 exposes private HTTP routes only. Expected route shape:

```text
http://service.project.devbox
```
