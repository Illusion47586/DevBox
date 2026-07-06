FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/devbox ./cmd/devbox

FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl docker.io gh git openssh-client tmux \
  && rm -rf /var/lib/apt/lists/*
RUN curl https://mise.run | sh \
  && ln -s /root/.local/bin/mise /usr/local/bin/mise
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
  && apt-get update \
  && apt-get install -y --no-install-recommends nodejs \
  && npm install -g @openai/codex @anthropic-ai/claude-code opencode-ai skills@1.5.14 \
  && rm -rf /var/lib/apt/lists/*
RUN skills add vercel-labs/agent-skills \
    --skill vercel-react-best-practices \
    --skill web-design-guidelines \
    --skill writing-guidelines \
    -g --agent codex claude-code opencode -y --copy \
  && skills add obra/superpowers \
    --skill test-driven-development \
    --skill systematic-debugging \
    --skill requesting-code-review \
    --skill receiving-code-review \
    --skill verification-before-completion \
    --skill writing-plans \
    -g --agent codex claude-code opencode -y --copy
COPY skills /opt/devbox/skills
RUN mkdir -p /root/.codex/skills /root/.agents/skills /root/.claude/skills \
  && cp -R /root/.agents/skills/. /root/.codex/skills/ \
  && cp -R /opt/devbox/skills/. /root/.codex/skills/ \
  && cp -R /opt/devbox/skills/. /root/.agents/skills/ \
  && cp -R /opt/devbox/skills/. /root/.claude/skills/
COPY --from=build /out/devbox /usr/local/bin/devbox
ENV DEVBOX_STATE_PATH=/var/lib/devbox/state.json \
    DEVBOX_WORKSPACE_ROOT=/workspaces \
    DEVBOX_DOMAIN=devbox \
    DEVBOX_PROJECT_IMAGE=ghcr.io/illusion47586/devbox:latest \
    DEVBOX_CADDY_ADMIN_URL=http://caddy:2019
WORKDIR /workspaces
ENTRYPOINT ["devbox"]
CMD ["daemon"]
