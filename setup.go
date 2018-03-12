package geoip

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

// GeoIP Comments me
type GeoIP struct {
	Next      httpserver.Handler
	DBHandler *maxminddb.Reader
	Config    Config
}

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("geoip", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}

	dbhandler, err := maxminddb.Open(config.DatabasePath)
	if err != nil {
		return c.Err("geoip: Can't open database: " + config.DatabasePath)
	}
	// Create new middleware
	newMiddleWare := func(next httpserver.Handler) httpserver.Handler {
		return &GeoIP{
			Next:      next,
			DBHandler: dbhandler,
			Config:    config,
		}
	}
	// Add middleware
	cfg := httpserver.GetConfig(c)
	cfg.AddMiddleware(newMiddleWare)

	return nil
}
