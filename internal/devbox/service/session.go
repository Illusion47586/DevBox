package service

import (
	"context"
	"strings"

	"github.com/dhruv/devbox/internal/devbox/model"
	"github.com/dhruv/devbox/internal/devbox/state"
)

func (o *Operator) CreateSession(ctx context.Context, input state.CreateSessionInput) (model.Session, error) {
	session, err := o.Store.CreateSession(input)
	if err != nil {
		return model.Session{}, err
	}
	project, err := o.Store.GetProject(input.Project)
	if err != nil {
		return model.Session{}, err
	}
	args := []string{"exec", project.ContainerName, "tmux", "new-session", "-d", "-s", session.TmuxName}
	args = append(args, strings.Join(input.Command, " "))
	return session, run(ctx, "", "docker", args...)
}

func (o *Operator) StopSession(ctx context.Context, projectName, sessionName string) error {
	session, err := o.Store.GetSession(projectName, sessionName)
	if err != nil {
		return err
	}
	project, err := o.Store.GetProject(projectName)
	if err != nil {
		return err
	}
	if err := run(ctx, "", "docker", "exec", project.ContainerName, "tmux", "kill-session", "-t", session.TmuxName); err != nil {
		return err
	}
	return o.Store.RemoveSession(projectName, sessionName)
}

func (o *Operator) AttachSession(ctx context.Context, projectName, sessionName string) error {
	session, err := o.Store.GetSession(projectName, sessionName)
	if err != nil {
		return err
	}
	project, err := o.Store.GetProject(projectName)
	if err != nil {
		return err
	}
	return foreground(ctx, "docker", "exec", "-it", project.ContainerName, "tmux", "attach", "-t", session.TmuxName)
}

func (o *Operator) Shell(ctx context.Context, projectName string) error {
	project, err := o.Store.GetProject(projectName)
	if err != nil {
		return err
	}
	return foreground(ctx, "docker", "exec", "-it", project.ContainerName, "bash")
}
