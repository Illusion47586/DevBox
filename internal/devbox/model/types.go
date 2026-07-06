package model

import "time"

type State struct {
	Version  int                `json:"version"`
	Projects map[string]Project `json:"projects"`
	Routes   map[string]Route   `json:"routes"`
	Sessions map[string]Session `json:"sessions"`
}

type Project struct {
	Name          string    `json:"name"`
	WorkspacePath string    `json:"workspace_path"`
	ContainerName string    `json:"container_name"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Route struct {
	Project   string    `json:"project"`
	Service   string    `json:"service"`
	Host      string    `json:"host"`
	URL       string    `json:"url"`
	Target    int       `json:"target"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	Project       string    `json:"project"`
	Name          string    `json:"name"`
	TmuxName      string    `json:"tmux_name"`
	Command       []string  `json:"command"`
	AttachCommand string    `json:"attach_command"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewState() State {
	return State{
		Version:  1,
		Projects: map[string]Project{},
		Routes:   map[string]Route{},
		Sessions: map[string]Session{},
	}
}
