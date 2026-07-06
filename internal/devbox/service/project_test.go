package service

import (
	"reflect"
	"testing"

	"github.com/dhruv/devbox/internal/devbox/model"
)

func TestProjectContainerRunArgsOverrideEntrypoint(t *testing.T) {
	project := model.Project{
		Name:          "myapp",
		WorkspacePath: "/workspaces/myapp",
		ContainerName: "devbox-project-myapp",
		Image:         "ghcr.io/illusion47586/devbox:latest",
	}

	args := projectContainerRunArgs(project)
	want := []string{
		"run", "-d",
		"--name", "devbox-project-myapp",
		"--hostname", "myapp",
		"--network", "devbox",
		"--network-alias", "myapp",
		"-v", "/workspaces/myapp:/workspace",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-w", "/workspace",
		"--entrypoint", "sleep",
		"ghcr.io/illusion47586/devbox:latest",
		"infinity",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("run args = %#v, want %#v", args, want)
	}
}
