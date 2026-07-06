#!/usr/bin/env sh
set -eu

dockerfile="Dockerfile"

for expected in \
  "COPY skills /opt/devbox/skills" \
  "/root/.codex/skills" \
  "/root/.agents/skills" \
  "/root/.claude/skills"; do
  if ! grep -Fq "$expected" "$dockerfile"; then
    printf 'Dockerfile must include %s\n' "$expected"
    exit 1
  fi
done

for skill in \
  skills/devbox-go-development/SKILL.md \
  skills/devbox-operations/SKILL.md \
  skills/devbox-agent-control/SKILL.md; do
  if [ ! -f "$skill" ]; then
    printf 'missing image skill source: %s\n' "$skill"
    exit 1
  fi
done
