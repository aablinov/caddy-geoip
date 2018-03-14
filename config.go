package geoip

import (
	"github.com/mholt/caddy"
)

// Config specifies configuration parsed for Caddyfile
type Config struct {
	DatabasePath string
}

func parseConfig(c *caddy.Controller) (Config, error) {
	var config = Config{}
	for c.Next() {
		value := c.Val()
		switch value {
		case "geoip":
			if !c.NextArg() {
				continue
			}
			config.DatabasePath = c.Val()
		}
	}
	return config, nil
}
