package geoip

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// Config specifies configuration parsed for Caddyfile
type Config struct {
	DatabasePath string
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (g *GeoIP) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.NextArg() {
			continue
		}
		g.Config.DatabasePath = d.Val()
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var g GeoIP
	err := g.UnmarshalCaddyfile(h.Dispenser)
	return g, err
}
