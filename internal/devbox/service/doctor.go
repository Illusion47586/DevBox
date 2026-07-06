package service

import (
	"context"
	"os"
	"path/filepath"
)

type DoctorCheck struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (o *Operator) Doctor(ctx context.Context) []DoctorCheck {
	return []DoctorCheck{
		checkPath("state_directory", filepath.Dir(o.Config.StatePath)),
		checkPath("workspace_root", o.Config.WorkspaceRoot),
		checkCommand(ctx, "docker", "docker", "version"),
		checkCommand(ctx, "git", "git", "--version"),
		checkCommand(ctx, "tailscale", "tailscale", "status"),
	}
}

func checkPath(name, path string) DoctorCheck {
	if _, err := os.Stat(path); err != nil {
		return DoctorCheck{Name: name, OK: false, Message: err.Error()}
	}
	return DoctorCheck{Name: name, OK: true, Message: path}
}

func checkCommand(ctx context.Context, name string, command string, args ...string) DoctorCheck {
	if err := run(ctx, "", command, args...); err != nil {
		return DoctorCheck{Name: name, OK: false, Message: err.Error()}
	}
	return DoctorCheck{Name: name, OK: true, Message: "available"}
}
