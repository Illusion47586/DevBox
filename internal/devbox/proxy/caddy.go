package proxy

import (
	"encoding/json"
	"strconv"

	"github.com/dhruv/devbox/internal/devbox/model"
)

type caddyConfig struct {
	Admin caddyAdmin `json:"admin"`
	Apps  caddyApps  `json:"apps"`
}

type caddyAdmin struct {
	Listen string `json:"listen"`
}

type caddyApps struct {
	HTTP caddyHTTPApp `json:"http"`
}

type caddyHTTPApp struct {
	Servers map[string]caddyServer `json:"servers"`
}

type caddyServer struct {
	Listen []string     `json:"listen"`
	Routes []caddyRoute `json:"routes"`
}

type caddyRoute struct {
	Match  []caddyMatch   `json:"match"`
	Handle []caddyHandler `json:"handle"`
}

type caddyMatch struct {
	Host []string `json:"host"`
}

type caddyHandler struct {
	Handler   string          `json:"handler"`
	Upstreams []caddyUpstream `json:"upstreams"`
}

type caddyUpstream struct {
	Dial string `json:"dial"`
}

func BuildCaddyConfig(routes []model.Route) ([]byte, error) {
	serverRoutes := make([]caddyRoute, 0, len(routes))
	for _, route := range routes {
		serverRoutes = append(serverRoutes, caddyRoute{
			Match: []caddyMatch{{Host: []string{route.Host}}},
			Handle: []caddyHandler{{
				Handler: "reverse_proxy",
				Upstreams: []caddyUpstream{{
					Dial: route.Project + ":" + strconv.Itoa(route.Target),
				}},
			}},
		})
	}

	cfg := caddyConfig{
		Admin: caddyAdmin{Listen: "0.0.0.0:2019"},
		Apps: caddyApps{HTTP: caddyHTTPApp{Servers: map[string]caddyServer{
			"devbox": {
				Listen: []string{":80"},
				Routes: serverRoutes,
			},
		}}},
	}
	return json.MarshalIndent(cfg, "", "  ")
}
