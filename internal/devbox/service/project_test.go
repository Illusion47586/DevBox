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

	args := projectContainerRunArgs(project, project.WorkspacePath)
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

func TestProjectContainerMountSourceUsesHostWorkspaceRoot(t *testing.T) {
	operator := NewOperator(RuntimeConfig{
		WorkspaceRoot:     "/workspaces",
		HostWorkspaceRoot: "/DATA/AppData/devbox/workspaces",
	})
	project := model.Project{
		Name:          "test",
		WorkspacePath: "/workspaces/test",
	}

	source := operator.projectMountSource(project)
	if source != "/DATA/AppData/devbox/workspaces/test" {
		t.Fatalf("mount source = %q, want host workspace path", source)
	}
}
