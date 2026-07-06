package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dhruv/devbox/internal/devbox/model"
	"github.com/dhruv/devbox/internal/devbox/proxy"
	"github.com/dhruv/devbox/internal/devbox/state"
)

func (o *Operator) AddRoute(ctx context.Context, input state.AddRouteInput) (model.Route, error) {
	input.Domain = o.Config.Domain
	route, err := o.Store.AddRoute(input)
	if err != nil {
		return model.Route{}, err
	}
	routes, err := o.Store.ListRoutes()
	if err != nil {
		return model.Route{}, err
	}
	if err := o.ApplyCaddyRoutes(ctx, routes); err != nil {
		return model.Route{}, err
	}
	return route, nil
}

func (o *Operator) ApplyCaddyRoutes(ctx context.Context, routes []model.Route) error {
	if strings.TrimSpace(o.Config.CaddyAdminURL) == "" {
		return nil
	}
	payload, err := proxy.BuildCaddyConfig(routes)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(o.Config.CaddyAdminURL, "/")+"/load", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("caddy returned %s", resp.Status)
	}
	return nil
}
