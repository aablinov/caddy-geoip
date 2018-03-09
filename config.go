package geoip

import (
	"github.com/mholt/caddy"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

// Config specifies configuration parsed for Caddyfile
type Config struct {
	DBHandler *maxminddb.Reader // Database's handler if it gets opened.

	// Yout can set returned header names in config
	// Country
	HeaderNameCountryCode string
	HeaderNameCountryIsEU string
	HeaderNameCountryName string

	// City
	HeaderNameCityName string

	// Location
	HeaderNameLocationLat      string
	HeaderNameLocationLon      string
	HeaderNameLocationTimeZone string
}

// NewConfig initialize new Config with default values
func NewConfig() Config {
	c := Config{}

	c.HeaderNameCountryCode = "X-Geoip-Country-Code"
	c.HeaderNameCountryIsEU = "X-Geoip-Country-Eu"
	c.HeaderNameCountryName = "X-Geoip-Country-Name"

	c.HeaderNameCityName = "X-Geoip-City-Name"

	c.HeaderNameLocationLat = "X-Geoip-Location-Lat"
	c.HeaderNameLocationLon = "X-Geoip-Location-Lon"
	c.HeaderNameLocationTimeZone = "X-Geoip-Location-Tz"
	return c
}

func parseConfig(c *caddy.Controller) (Config, error) {
	var config = NewConfig()
	for c.Next() {
		for c.NextBlock() {
			value := c.Val()

			switch value {
			case "database":
				if !c.NextArg() {
					continue
				}
				// Check if a database has already been opened
				if config.DBHandler != nil {
					continue
				}

				database := c.Val()

				// Open the database.
				var err error
				config.DBHandler, err = maxminddb.Open(database)
				if err != nil {
					return config, c.Err("geoip: Can't open database: " + database)
				}
			case "set_header_country_code":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameCountryCode = c.Val()
			case "set_header_country_name":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameCountryName = c.Val()
			case "set_header_country_eu":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameCountryIsEU = c.Val()
			case "set_header_city_name":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameCityName = c.Val()
			case "set_header_location_lat":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameLocationLat = c.Val()
			case "set_header_location_lon":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameLocationLon = c.Val()
			case "set_header_location_tz":
				if !c.NextArg() {
					continue
				}
				config.HeaderNameLocationTimeZone = c.Val()
			}
		}
	}
	return config, nil
}
