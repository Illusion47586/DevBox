package proxy

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dhruv/devbox/internal/devbox/model"
)

func TestCaddyConfigIncludesRoutes(t *testing.T) {
	routes := []model.Route{{
		Project: "myapp",
		Service: "web",
		Host:    "web.myapp.devbox",
		Target:  3000,
	}}

	payload, err := BuildCaddyConfig(routes)
	if err != nil {
		t.Fatalf("build caddy config: %v", err)
	}

	var cfg map[string]any
	if err := json.Unmarshal(payload, &cfg); err != nil {
		t.Fatalf("unmarshal caddy config: %v", err)
	}

	apps := cfg["apps"].(map[string]any)
	admin := cfg["admin"].(map[string]any)
	if admin["listen"] != "0.0.0.0:2019" {
		t.Fatalf("admin listen = %v, want 0.0.0.0:2019", admin["listen"])
	}

	httpApp := apps["http"].(map[string]any)
	servers := httpApp["servers"].(map[string]any)
	devbox := servers["devbox"].(map[string]any)
	if devbox["listen"].([]any)[0] != ":80" {
		t.Fatalf("listen = %v, want :80", devbox["listen"])
	}

	encoded := string(payload)
	if !strings.Contains(encoded, "web.myapp.devbox") {
		t.Fatalf("config did not contain route host: %s", encoded)
	}
	if !strings.Contains(encoded, "myapp:3000") {
		t.Fatalf("config did not contain upstream target: %s", encoded)
	}
}
