package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/dhruv/devbox/internal/devbox/model"
	"github.com/dhruv/devbox/internal/devbox/service"
)

func TestProjectRegisterJSON(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	var stdout bytes.Buffer

	cmd := NewRootCommand(service.RuntimeConfig{
		StatePath: statePath,
		Domain:    "devbox",
	}, &stdout)
	cmd.SetArgs([]string{"project", "register", "/workspaces/myapp", "--name", "myapp", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout.String())
	}
	if response["ok"] != true {
		t.Fatalf("ok = %v", response["ok"])
	}
	project := response["project"].(map[string]any)
	if project["name"] != "myapp" {
		t.Fatalf("project name = %v", project["name"])
	}
}

func TestRouteAddJSON(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	var stdout bytes.Buffer

	cmd := NewRootCommand(service.RuntimeConfig{
		StatePath: statePath,
		Domain:    "devbox",
	}, &stdout)
	cmd.SetArgs([]string{"project", "register", "/workspaces/myapp", "--name", "myapp", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("register: %v", err)
	}

	stdout.Reset()
	cmd = NewRootCommand(service.RuntimeConfig{StatePath: statePath, Domain: "devbox"}, &stdout)
	cmd.SetArgs([]string{"route", "add", "--project", "myapp", "--service", "web", "--target", "3000", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("route add: %v", err)
	}

	var response struct {
		OK    bool        `json:"ok"`
		Route model.Route `json:"route"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout.String())
	}
	if !response.OK {
		t.Fatal("ok = false")
	}
	if response.Route.URL != "http://web.myapp.devbox" {
		t.Fatalf("route url = %q", response.Route.URL)
	}
}
