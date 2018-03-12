package geoip

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

func TestToResolveGeoip(t *testing.T) {
	dbhandler, err := maxminddb.Open("./db/GeoLite2-City.mmdb")
	if err != nil {
		t.Errorf("geoip: Can't open database: GeoLite2-City.mmdb")
	}

	config := Config{}

	config.HeaderNameCountryCode = "X-Geoip-Country-Code"
	config.HeaderNameCountryIsEU = "X-Geoip-Country-Eu"
	config.HeaderNameCountryName = "X-Geoip-Country-Name"

	config.HeaderNameCityName = "X-Geoip-City-Name"

	config.HeaderNameLocationLat = "X-Geoip-Location-Lat"
	config.HeaderNameLocationLon = "X-Geoip-Location-Lon"
	config.HeaderNameLocationTimeZone = "X-Geoip-Location-Tz"

	var (
		gotHeaders      http.Header
		expectedHeaders = http.Header{
			"X-Geoip-Country-Code": []string{"CY"},
			"X-Geoip-Location-Lat": []string{"34.684100"},
			"X-Geoip-Location-Lon": []string{"33.037900"},
			"X-Geoip-Location-Tz":  []string{"Asia/Nicosia"},
			"X-Geoip-Country-Eu":   []string{"false"},
			"X-Geoip-Country-Name": []string{"Cyprus"},
			"X-Geoip-City-Name":    []string{"Limassol"},
		}
	)
	l := GeoIP{
		Next: httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
			gotHeaders = r.Header
			return 0, nil
		}),
		DBHandler: dbhandler,
		Config:    config,
	}

	r := httptest.NewRequest("GET", "/", strings.NewReader(""))
	r.RemoteAddr = "212.50.99.193"
	l.ServeHTTP(httptest.NewRecorder(), r)

	if !reflect.DeepEqual(expectedHeaders, gotHeaders) {
		t.Errorf("Expected %v actual %v", expectedHeaders, gotHeaders)
	}
}
