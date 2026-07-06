#!/usr/bin/env sh
set -eu

dockerfile="Dockerfile"

for expected in \
  "npm install -g @openai/codex @anthropic-ai/claude-code opencode-ai skills@" \
  "skills add vercel-labs/agent-skills" \
  "--skill vercel-react-best-practices" \
  "--skill web-design-guidelines" \
  "--skill writing-guidelines" \
  "skills add obra/superpowers" \
  "--skill test-driven-development" \
  "--skill systematic-debugging" \
  "--skill requesting-code-review" \
  "--skill receiving-code-review" \
  "--skill verification-before-completion" \
  "--skill writing-plans" \
  "COPY skills /opt/devbox/skills" \
  "/root/.codex/skills" \
  "/root/.agents/skills" \
  "/root/.claude/skills"; do
  if ! grep -Fq -- "$expected" "$dockerfile"; then
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
