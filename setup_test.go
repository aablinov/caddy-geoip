package geoip

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type testResponseRecorder struct {
	*httpserver.ResponseWriterWrapper
}

func (testResponseRecorder) CloseNotify() <-chan bool { return nil }

func TestReplacers(t *testing.T) {
	dbhandler, err := maxminddb.Open("./test-data/GeoLite2-City.mmdb")
	if err != nil {
		t.Errorf("geoip: Can't open database: GeoLite2-City.mmdb")
	}

	config := Config{}

	l := GeoIP{
		Next: httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
			return 0, nil
		}),
		DBHandler: dbhandler,
		Config:    config,
	}

	r := httptest.NewRequest("GET", "/", strings.NewReader(""))
	r.RemoteAddr = "212.50.99.193"
	rr := httpserver.NewResponseRecorder(testResponseRecorder{
		ResponseWriterWrapper: &httpserver.ResponseWriterWrapper{ResponseWriter: httptest.NewRecorder()},
	})

	rr.Replacer = httpserver.NewReplacer(r, rr, "-")

	l.ServeHTTP(rr, r)

	if got, want := rr.Replacer.Replace("{geoip_country_code}"), "CY"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_code} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_country_name}"), "Cyprus"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_name} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_country_eu}"), "false"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_eu} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_city_name}"), "Limassol"; got != want {
		t.Errorf("Expected custom placeholder {geoip_city_name} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_latitude}"), "34.684100"; got != want {
		t.Errorf("Expected custom placeholder {geoip_latitude} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_longitude}"), "33.037900"; got != want {
		t.Errorf("Expected custom placeholder {geoip_longitude} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_time_zone}"), "Asia/Nicosia"; got != want {
		t.Errorf("Expected custom placeholder {geoip_time_zone} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_geohash}"), "swpmrf13wbgg"; got != want {
		t.Errorf("Expected custom placeholder {geoip_geohash} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_city_geoname_id}"), "146384"; got != want {
		t.Errorf("Expected custom placeholder {geoip_city_geoname_id} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := rr.Replacer.Replace("{geoip_country_geoname_id}"), "146669"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_geoname_id} to be set (%s), but it wasn't; got: %s", want, got)
	}
}
