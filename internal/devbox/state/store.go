package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dhruv/devbox/internal/devbox/model"
)

var resourceNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

type Store struct {
	path string
}

type RegisterProjectInput struct {
	Name          string
	WorkspacePath string
	Image         string
}

type AddRouteInput struct {
	Project string
	Service string
	Target  int
	Domain  string
}

type CreateSessionInput struct {
	Project string
	Name    string
	Command []string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Load() (model.State, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return model.NewState(), nil
	}
	if err != nil {
		return model.State{}, err
	}

	var current model.State
	if err := json.Unmarshal(data, &current); err != nil {
		return model.State{}, err
	}
	if current.Version == 0 {
		current.Version = 1
	}
	if current.Projects == nil {
		current.Projects = map[string]model.Project{}
	}
	if current.Routes == nil {
		current.Routes = map[string]model.Route{}
	}
	if current.Sessions == nil {
		current.Sessions = map[string]model.Session{}
	}
	return current, nil
}

func (s *Store) Save(current model.State) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(s.path, data, 0o644)
}

func (s *Store) RegisterProject(input RegisterProjectInput) (model.Project, error) {
	if err := ValidateName("project", input.Name); err != nil {
		return model.Project{}, err
	}
	if strings.TrimSpace(input.WorkspacePath) == "" {
		return model.Project{}, errors.New("workspace path is required")
	}

	current, err := s.Load()
	if err != nil {
		return model.Project{}, err
	}

	now := time.Now().UTC()
	existing := current.Projects[input.Name]
	createdAt := now
	if !existing.CreatedAt.IsZero() {
		createdAt = existing.CreatedAt
	}
	image := input.Image
	if image == "" {
		image = existing.Image
	}
	project := model.Project{
		Name:          input.Name,
		WorkspacePath: input.WorkspacePath,
		ContainerName: "devbox-project-" + input.Name,
		Image:         image,
		CreatedAt:     createdAt,
		UpdatedAt:     now,
	}
	current.Projects[input.Name] = project
	return project, s.Save(current)
}

func (s *Store) ListProjects() ([]model.Project, error) {
	current, err := s.Load()
	if err != nil {
		return nil, err
	}
	projects := make([]model.Project, 0, len(current.Projects))
	for _, project := range current.Projects {
		projects = append(projects, project)
	}
	return projects, nil
}

func (s *Store) GetProject(name string) (model.Project, error) {
	current, err := s.Load()
	if err != nil {
		return model.Project{}, err
	}
	project, ok := current.Projects[name]
	if !ok {
		return model.Project{}, fmt.Errorf("project %q is not registered", name)
	}
	return project, nil
}

func (s *Store) RemoveProject(name string) error {
	current, err := s.Load()
	if err != nil {
		return err
	}
	delete(current.Projects, name)
	for key, route := range current.Routes {
		if route.Project == name {
			delete(current.Routes, key)
		}
	}
	for key, session := range current.Sessions {
		if session.Project == name {
			delete(current.Sessions, key)
		}
	}
	return s.Save(current)
}

func (s *Store) AddRoute(input AddRouteInput) (model.Route, error) {
	if err := ValidateName("project", input.Project); err != nil {
		return model.Route{}, err
	}
	if err := ValidateName("service", input.Service); err != nil {
		return model.Route{}, err
	}
	if input.Target <= 0 || input.Target > 65535 {
		return model.Route{}, errors.New("target must be a valid TCP port")
	}

	current, err := s.Load()
	if err != nil {
		return model.Route{}, err
	}
	if _, ok := current.Projects[input.Project]; !ok {
		return model.Route{}, fmt.Errorf("project %q is not registered", input.Project)
	}

	now := time.Now().UTC()
	key := RouteKey(input.Project, input.Service)
	existing := current.Routes[key]
	createdAt := now
	if !existing.CreatedAt.IsZero() {
		createdAt = existing.CreatedAt
	}
	domain := strings.Trim(strings.TrimSpace(input.Domain), ".")
	if domain == "" {
		domain = "devbox"
	}
	host := input.Service + "." + input.Project + "." + domain
	route := model.Route{
		Project:   input.Project,
		Service:   input.Service,
		Host:      host,
		URL:       "http://" + host,
		Target:    input.Target,
		CreatedAt: createdAt,
		UpdatedAt: now,
	}
	current.Routes[key] = route
	return route, s.Save(current)
}

func (s *Store) ListRoutes() ([]model.Route, error) {
	current, err := s.Load()
	if err != nil {
		return nil, err
	}
	routes := make([]model.Route, 0, len(current.Routes))
	for _, route := range current.Routes {
		routes = append(routes, route)
	}
	return routes, nil
}

func (s *Store) RemoveRoute(project, service string) error {
	current, err := s.Load()
	if err != nil {
		return err
	}
	delete(current.Routes, RouteKey(project, service))
	return s.Save(current)
}

func (s *Store) CreateSession(input CreateSessionInput) (model.Session, error) {
	if err := ValidateName("project", input.Project); err != nil {
		return model.Session{}, err
	}
	if err := ValidateName("session", input.Name); err != nil {
		return model.Session{}, err
	}
	if len(input.Command) == 0 {
		return model.Session{}, errors.New("session command is required")
	}

	current, err := s.Load()
	if err != nil {
		return model.Session{}, err
	}
	if _, ok := current.Projects[input.Project]; !ok {
		return model.Session{}, fmt.Errorf("project %q is not registered", input.Project)
	}

	now := time.Now().UTC()
	key := SessionKey(input.Project, input.Name)
	existing := current.Sessions[key]
	createdAt := now
	if !existing.CreatedAt.IsZero() {
		createdAt = existing.CreatedAt
	}
	session := model.Session{
		Project:       input.Project,
		Name:          input.Name,
		TmuxName:      "devbox-" + input.Project + "-" + input.Name,
		Command:       input.Command,
		AttachCommand: "devbox session attach " + input.Project + " " + input.Name,
		CreatedAt:     createdAt,
		UpdatedAt:     now,
	}
	current.Sessions[key] = session
	return session, s.Save(current)
}

func (s *Store) ListSessions(project string) ([]model.Session, error) {
	current, err := s.Load()
	if err != nil {
		return nil, err
	}
	sessions := make([]model.Session, 0, len(current.Sessions))
	for _, session := range current.Sessions {
		if project == "" || session.Project == project {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (s *Store) GetSession(project, name string) (model.Session, error) {
	current, err := s.Load()
	if err != nil {
		return model.Session{}, err
	}
	session, ok := current.Sessions[SessionKey(project, name)]
	if !ok {
		return model.Session{}, fmt.Errorf("session %q for project %q is not registered", name, project)
	}
	return session, nil
}

func (s *Store) RemoveSession(project, name string) error {
	current, err := s.Load()
	if err != nil {
		return err
	}
	delete(current.Sessions, SessionKey(project, name))
	return s.Save(current)
}

func RouteKey(project, service string) string {
	return project + "/" + service
}

func SessionKey(project, name string) string {
	return project + "/" + name
}

func ValidateName(kind, value string) error {
	if !resourceNamePattern.MatchString(value) {
		return fmt.Errorf("%s name %q must use lowercase letters, numbers, and hyphens", kind, value)
	}
	return nil
}
