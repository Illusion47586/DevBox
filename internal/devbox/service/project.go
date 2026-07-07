package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/dhruv/devbox/internal/devbox/model"
	"github.com/dhruv/devbox/internal/devbox/state"
)

func (o *Operator) CreateProject(ctx context.Context, gitURL, name string) (model.Project, error) {
	if strings.TrimSpace(gitURL) == "" {
		return model.Project{}, errors.New("git url is required")
	}
	if name == "" {
		name = deriveProjectName(gitURL)
	}
	workspacePath := filepath.Join(o.Config.WorkspaceRoot, name)
	project, err := o.Store.RegisterProject(state.RegisterProjectInput{
		Name:          name,
		WorkspacePath: workspacePath,
		Image:         o.Config.ProjectImage,
	})
	if err != nil {
		return model.Project{}, err
	}
	if err := run(ctx, "", "git", "clone", gitURL, workspacePath); err != nil {
		return model.Project{}, err
	}
	return project, o.EnsureProjectContainer(ctx, project)
}

func (o *Operator) RegisterProject(ctx context.Context, path, name string) (model.Project, error) {
	if name == "" {
		name = filepath.Base(filepath.Clean(path))
	}
	project, err := o.Store.RegisterProject(state.RegisterProjectInput{
		Name:          name,
		WorkspacePath: path,
		Image:         o.Config.ProjectImage,
	})
	if err != nil {
		return model.Project{}, err
	}
	return project, o.EnsureProjectContainer(ctx, project)
}

func (o *Operator) EnsureProjectContainer(ctx context.Context, project model.Project) error {
	if project.Image == "" {
		project.Image = o.Config.ProjectImage
	}
	if project.Image == "" {
		return nil
	}
	if err := run(ctx, "", "docker", "inspect", project.ContainerName); err == nil {
		return run(ctx, "", "docker", "start", project.ContainerName)
	}
	return run(ctx, "", "docker", projectContainerRunArgs(project, o.projectMountSource(project))...)
}

func (o *Operator) projectMountSource(project model.Project) string {
	if o.Config.HostWorkspaceRoot == "" {
		return project.WorkspacePath
	}
	rel, err := filepath.Rel(o.Config.WorkspaceRoot, project.WorkspacePath)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return project.WorkspacePath
	}
	return filepath.Join(o.Config.HostWorkspaceRoot, rel)
}

func projectContainerRunArgs(project model.Project, mountSource string) []string {
	return []string{"run", "-d",
		"--name", project.ContainerName,
		"--hostname", project.Name,
		"--network", "devbox",
		"--network-alias", project.Name,
		"-v", mountSource + ":/workspace",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-w", "/workspace",
		"--entrypoint", "sleep",
		project.Image,
		"infinity",
	}
}
