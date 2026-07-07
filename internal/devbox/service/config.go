package service

import (
	"os"
	"strings"
)

type RuntimeConfig struct {
	StatePath         string
	WorkspaceRoot     string
	HostWorkspaceRoot string
	Domain            string
	ProjectImage      string
	CaddyAdminURL     string
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		StatePath:         envOrDefault("DEVBOX_STATE_PATH", "/var/lib/devbox/state.json"),
		WorkspaceRoot:     envOrDefault("DEVBOX_WORKSPACE_ROOT", "/workspaces"),
		HostWorkspaceRoot: envOrDefault("DEVBOX_HOST_WORKSPACE_ROOT", ""),
		Domain:            envOrDefault("DEVBOX_DOMAIN", "devbox"),
		ProjectImage:      envOrDefault("DEVBOX_PROJECT_IMAGE", "ghcr.io/illusion47586/devbox:latest"),
		CaddyAdminURL:     envOrDefault("DEVBOX_CADDY_ADMIN_URL", "http://127.0.0.1:2019"),
	}
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
