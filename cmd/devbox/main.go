package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dhruv/devbox/internal/devbox/service"
	"github.com/dhruv/devbox/internal/devbox/state"
	"github.com/spf13/cobra"
)

func main() {
	cmd := NewRootCommand(service.DefaultRuntimeConfig(), os.Stdout)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewRootCommand(config service.RuntimeConfig, stdout io.Writer) *cobra.Command {
	var jsonOut bool
	operator := service.NewOperator(config)

	root := &cobra.Command{
		Use:           "devbox",
		Short:         "Agent-first development control plane for TheBox",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().BoolVar(&jsonOut, "json", false, "write machine-readable JSON")

	root.AddCommand(projectCommand(operator, stdout, &jsonOut))
	root.AddCommand(routeCommand(operator, stdout, &jsonOut))
	root.AddCommand(sessionCommand(operator, stdout, &jsonOut))
	root.AddCommand(shellCommand(operator))
	root.AddCommand(zedCommand(operator, stdout, &jsonOut))
	root.AddCommand(doctorCommand(operator, stdout, &jsonOut))
	root.AddCommand(daemonCommand(stdout))

	return root
}

func daemonCommand(stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Keep the DevBox control-plane container running",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(stdout, "devbox control plane ready")
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-cmd.Context().Done():
				return cmd.Context().Err()
			case <-signals:
				return nil
			}
		},
	}
}

func projectCommand(operator *service.Operator, stdout io.Writer, jsonOut *bool) *cobra.Command {
	cmd := &cobra.Command{Use: "project", Short: "Manage DevBox projects"}

	var createName string
	create := &cobra.Command{
		Use:   "create <git-url>",
		Short: "Clone and register a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := operator.CreateProject(cmd.Context(), args[0], createName)
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "project", project, "registered "+project.Name)
		},
	}
	create.Flags().StringVar(&createName, "name", "", "project name")

	var registerName string
	register := &cobra.Command{
		Use:   "register <path>",
		Short: "Register an existing workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := operator.RegisterProject(cmd.Context(), args[0], registerName)
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "project", project, "registered "+project.Name)
		},
	}
	register.Flags().StringVar(&registerName, "name", "", "project name")

	list := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := operator.Store.ListProjects()
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "projects", projects, fmt.Sprintf("%d projects", len(projects)))
		},
	}

	status := &cobra.Command{
		Use:   "status <project>",
		Short: "Show project status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := operator.Store.GetProject(args[0])
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "project", project, project.Name)
		},
	}

	remove := &cobra.Command{
		Use:   "remove <project>",
		Short: "Remove project state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := operator.Store.RemoveProject(args[0]); err != nil {
				return err
			}
			return write(stdout, *jsonOut, "removed", args[0], "removed "+args[0])
		},
	}

	cmd.AddCommand(create, register, list, status, remove)
	return cmd
}

func routeCommand(operator *service.Operator, stdout io.Writer, jsonOut *bool) *cobra.Command {
	cmd := &cobra.Command{Use: "route", Short: "Manage private development URLs"}

	var project, service string
	var target int
	add := &cobra.Command{
		Use:   "add",
		Short: "Register a service route",
		RunE: func(cmd *cobra.Command, args []string) error {
			route, err := operator.AddRoute(cmd.Context(), state.AddRouteInput{
				Project: project,
				Service: service,
				Target:  target,
			})
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "route", route, route.URL)
		},
	}
	add.Flags().StringVar(&project, "project", "", "project name")
	add.Flags().StringVar(&service, "service", "", "service name")
	add.Flags().IntVar(&target, "target", 0, "target port")
	_ = add.MarkFlagRequired("project")
	_ = add.MarkFlagRequired("service")
	_ = add.MarkFlagRequired("target")

	list := &cobra.Command{
		Use:   "list",
		Short: "List routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			routes, err := operator.Store.ListRoutes()
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "routes", routes, fmt.Sprintf("%d routes", len(routes)))
		},
	}

	remove := &cobra.Command{
		Use:   "remove <project> <service>",
		Short: "Remove a route",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := operator.Store.RemoveRoute(args[0], args[1]); err != nil {
				return err
			}
			routes, err := operator.Store.ListRoutes()
			if err != nil {
				return err
			}
			if err := operator.ApplyCaddyRoutes(cmd.Context(), routes); err != nil {
				return err
			}
			return write(stdout, *jsonOut, "removed", args[0]+"/"+args[1], "removed route")
		},
	}

	cmd.AddCommand(add, list, remove)
	return cmd
}

func sessionCommand(operator *service.Operator, stdout io.Writer, jsonOut *bool) *cobra.Command {
	cmd := &cobra.Command{Use: "session", Short: "Manage project tmux sessions"}

	var name string
	create := &cobra.Command{
		Use:   "create <project> --name <name> -- <command...>",
		Short: "Create a tmux session inside a project container",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := operator.CreateSession(cmd.Context(), state.CreateSessionInput{
				Project: args[0],
				Name:    name,
				Command: args[1:],
			})
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "session", session, session.AttachCommand)
		},
	}
	create.Flags().StringVar(&name, "name", "", "session name")
	_ = create.MarkFlagRequired("name")

	var projectFilter string
	list := &cobra.Command{
		Use:   "list",
		Short: "List sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, err := operator.Store.ListSessions(projectFilter)
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "sessions", sessions, fmt.Sprintf("%d sessions", len(sessions)))
		},
	}
	list.Flags().StringVar(&projectFilter, "project", "", "filter by project")

	attach := &cobra.Command{
		Use:   "attach <project> <session>",
		Short: "Attach to a tmux session",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return operator.AttachSession(cmd.Context(), args[0], args[1])
		},
	}

	stop := &cobra.Command{
		Use:   "stop <project> <session>",
		Short: "Stop a tmux session",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := operator.StopSession(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			return write(stdout, *jsonOut, "stopped", args[0]+"/"+args[1], "stopped session")
		},
	}

	status := &cobra.Command{
		Use:   "status <project> <session>",
		Short: "Show session status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := operator.Store.GetSession(args[0], args[1])
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "session", session, session.AttachCommand)
		},
	}

	cmd.AddCommand(create, list, attach, stop, status)
	return cmd
}

func shellCommand(operator *service.Operator) *cobra.Command {
	return &cobra.Command{
		Use:   "shell <project>",
		Short: "Open a shell inside a project container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return operator.Shell(cmd.Context(), args[0])
		},
	}
}

func zedCommand(operator *service.Operator, stdout io.Writer, jsonOut *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "zed <project>",
		Short: "Print a Zed SSH URL for a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := operator.ZedURL(args[0])
			if err != nil {
				return err
			}
			return write(stdout, *jsonOut, "url", url, url)
		},
	}
}

func doctorCommand(operator *service.Operator, stdout io.Writer, jsonOut *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check DevBox runtime dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			checks := operator.Doctor(context.Background())
			return write(stdout, *jsonOut, "checks", checks, strconv.Itoa(len(checks))+" checks")
		},
	}
}

func write(stdout io.Writer, jsonOut bool, key string, value any, text string) error {
	if !jsonOut {
		_, err := fmt.Fprintln(stdout, text)
		return err
	}
	response := map[string]any{
		"ok": true,
		key:  value,
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}
