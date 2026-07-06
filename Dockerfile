FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/devbox ./cmd/devbox

FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl git docker.io tmux openssh-client \
  && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/devbox /usr/local/bin/devbox
ENV DEVBOX_STATE_PATH=/var/lib/devbox/state.json \
    DEVBOX_WORKSPACE_ROOT=/workspaces \
    DEVBOX_DOMAIN=devbox \
    DEVBOX_PROJECT_IMAGE=devbox-project:local \
    DEVBOX_CADDY_ADMIN_URL=http://caddy:2019
WORKDIR /workspaces
ENTRYPOINT ["devbox"]
CMD ["daemon"]
