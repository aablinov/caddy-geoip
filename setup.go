package geoip

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/mmcloughlin/geohash"
	"github.com/oschwald/maxminddb-golang"
)

// Init initializes the module
func init() {
	caddy.RegisterModule(GeoIP{})
	httpcaddyfile.RegisterHandlerDirective("geoip", parseCaddyfile)
}

// GeoIP represents a middleware instance
type GeoIP struct {
	DBHandler *maxminddb.Reader
	Config    Config
}

type geoIPRecord struct {
	Country struct {
		ISOCode           string            `maxminddb:"iso_code"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		Names             map[string]string `maxminddb:"names"`
		GeoNameID         uint64            `maxminddb:"geoname_id"`
	} `maxminddb:"country"`

	City struct {
		Names     map[string]string `maxminddb:"names"`
		GeoNameID uint64            `maxminddb:"geoname_id"`
	} `maxminddb:"city"`

	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

// CaddyModule returns the Caddy module information.
func (GeoIP) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.geoip",
		New: func() caddy.Module { return new(GeoIP) },
	}
}

// Provision implements caddy.Provisioner.
func (g *GeoIP) Provision(ctx caddy.Context) error {
	dbPath := g.Config.DatabasePath
	if dbPath == "" {
		return fmt.Errorf("a db path is required")
	}
	dbhandler, err := maxminddb.Open(dbPath)
	if err != nil {
		return fmt.Errorf("geoip: Can't open database: " + dbPath)
	}
	g.DBHandler = dbhandler
	return nil
}

// Validate implements caddy.Validator.
func (g *GeoIP) Validate() error {
	if g.DBHandler == nil {
		return fmt.Errorf("no db")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (g GeoIP) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	g.lookupLocation(w, r)
	return next.ServeHTTP(w, r)
}

// Interface guards
var (
	_ caddy.Provisioner           = (*GeoIP)(nil)
	_ caddy.Validator             = (*GeoIP)(nil)
	_ caddyhttp.MiddlewareHandler = (*GeoIP)(nil)
	_ caddyfile.Unmarshaler       = (*GeoIP)(nil)
)

func (g GeoIP) lookupLocation(w http.ResponseWriter, r *http.Request) {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	record := g.fetchGeoipData(r)

	repl.Set("geoip_country_code", record.Country.ISOCode)
	repl.Set("geoip_country_name", record.Country.Names["en"])
	repl.Set("geoip_country_eu", strconv.FormatBool(record.Country.IsInEuropeanUnion))
	repl.Set("geoip_country_geoname_id", strconv.FormatUint(record.Country.GeoNameID, 10))
	repl.Set("geoip_city_name", record.City.Names["en"])
	repl.Set("geoip_city_geoname_id", strconv.FormatUint(record.City.GeoNameID, 10))
	repl.Set("geoip_latitude", strconv.FormatFloat(record.Location.Latitude, 'f', 6, 64))
	repl.Set("geoip_longitude", strconv.FormatFloat(record.Location.Longitude, 'f', 6, 64))
	repl.Set("geoip_geohash", geohash.Encode(record.Location.Latitude, record.Location.Longitude))
	repl.Set("geoip_time_zone", record.Location.TimeZone)
}

func (g GeoIP) fetchGeoipData(r *http.Request) geoIPRecord {
	clientIP, _ := getClientIP(r, true)

	var record = geoIPRecord{}
	err := g.DBHandler.Lookup(clientIP, &record)
	if err != nil {
		log.Println(err)
	}

	if record.Country.ISOCode == "" {
		record.Country.Names = make(map[string]string)
		record.City.Names = make(map[string]string)
		if clientIP.IsLoopback() {
			record.Country.ISOCode = "**"
			record.Country.Names["en"] = "Loopback"
			record.City.Names["en"] = "Loopback"
		} else {
			record.Country.ISOCode = "!!"
			record.Country.Names["en"] = "No Country"
			record.City.Names["en"] = "No City"
		}
	}

	return record
}

func getClientIP(r *http.Request, strict bool) (net.IP, error) {
	var ip string

	// Use the client ip from the 'X-Forwarded-For' header, if available.
	if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" && !strict {
		ips := strings.Split(fwdFor, ", ")
		ip = ips[0]
	} else {
		// Otherwise, get the client ip from the request remote address.
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			if serr, ok := err.(*net.AddrError); ok && serr.Err == "missing port in address" { // It's not critical try parse
				ip = r.RemoteAddr
			} else {
				log.Printf("Error when SplitHostPort: %v", serr.Err)
				return nil, err
			}
		}
	}

	// Parse the ip address string into a net.IP.
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, errors.New("unable to parse address")
	}

	return parsedIP, nil
}
