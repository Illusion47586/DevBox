package state

import (
	"path/filepath"
	"testing"
)

func TestRegisterProjectPersistsProject(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(filepath.Join(dir, "state.json"))

	project, err := store.RegisterProject(RegisterProjectInput{
		Name:          "myapp",
		WorkspacePath: "/workspaces/myapp",
	})
	if err != nil {
		t.Fatalf("register project: %v", err)
	}
	if project.Name != "myapp" {
		t.Fatalf("project name = %q, want myapp", project.Name)
	}

	reopened := NewStore(filepath.Join(dir, "state.json"))
	current, err := reopened.Load()
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if len(current.Projects) != 1 {
		t.Fatalf("projects len = %d, want 1", len(current.Projects))
	}
	if current.Projects["myapp"].WorkspacePath != "/workspaces/myapp" {
		t.Fatalf("workspace path = %q", current.Projects["myapp"].WorkspacePath)
	}
}

func TestRegisterProjectRejectsInvalidName(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "state.json"))

	_, err := store.RegisterProject(RegisterProjectInput{
		Name:          "Bad Name",
		WorkspacePath: "/workspaces/bad",
	})
	if err == nil {
		t.Fatal("expected invalid project name error")
	}
}

func TestAddRouteBuildsCanonicalURL(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "state.json"))
	_, err := store.RegisterProject(RegisterProjectInput{Name: "myapp", WorkspacePath: "/workspaces/myapp"})
	if err != nil {
		t.Fatalf("register project: %v", err)
	}

	route, err := store.AddRoute(AddRouteInput{
		Project: "myapp",
		Service: "web",
		Target:  3000,
		Domain:  "devbox",
	})
	if err != nil {
		t.Fatalf("add route: %v", err)
	}
	if route.URL != "http://web.myapp.devbox" {
		t.Fatalf("url = %q, want http://web.myapp.devbox", route.URL)
	}
	if route.Target != 3000 {
		t.Fatalf("target = %d, want 3000", route.Target)
	}
}

func TestCreateSessionTracksTmuxSession(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "state.json"))
	_, err := store.RegisterProject(RegisterProjectInput{Name: "myapp", WorkspacePath: "/workspaces/myapp"})
	if err != nil {
		t.Fatalf("register project: %v", err)
	}

	session, err := store.CreateSession(CreateSessionInput{
		Project: "myapp",
		Name:    "agent-main",
		Command: []string{"codex", "exec"},
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if session.TmuxName != "devbox-myapp-agent-main" {
		t.Fatalf("tmux name = %q", session.TmuxName)
	}
	if session.AttachCommand != "devbox session attach myapp agent-main" {
		t.Fatalf("attach command = %q", session.AttachCommand)
	}
}
