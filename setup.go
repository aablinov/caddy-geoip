package geoip

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// GeoIP Comments me
type GeoIP struct {
	Next   httpserver.Handler
	Config Config
}

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("geoip", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	ifconfig, err := parseConfig(c)
	if err != nil {
		return err
	}

	// Create new middleware
	newMiddleWare := func(next httpserver.Handler) httpserver.Handler {
		return &GeoIP{
			Next:   next,
			Config: ifconfig,
		}
	}
	// Add middleware
	cfg := httpserver.GetConfig(c)
	cfg.AddMiddleware(newMiddleWare)

	return nil
}
