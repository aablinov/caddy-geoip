package geoip

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/caddyserver/caddy/v2"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type testResponseRecorder struct {
	*caddyhttp.ResponseWriterWrapper
}

func (testResponseRecorder) CloseNotify() <-chan bool { return nil }

func get(replacer *caddy.Replacer, s string) string {
	val, _ := replacer.GetString(s)
	return val
}

func next(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func TestReplacers(t *testing.T) {
	dbhandler, err := maxminddb.Open("./test-data/GeoLite2-City.mmdb")
	if err != nil {
		t.Errorf("geoip: Can't open database: GeoLite2-City.mmdb")
	}

	buf := new(bytes.Buffer)

	// shouldBuf determines whether to execute templates on this response,
	// since generally we will not want to execute for images or CSS, etc.
	shouldBuf := func(status int, header http.Header) bool {
		return true
	}

	config := Config{}

	l := GeoIP{
		DBHandler: dbhandler,
		Config:    config,
	}

	r := httptest.NewRequest("GET", "/", strings.NewReader(""))
	r.RemoteAddr = "212.50.99.193"

	rr := caddyhttp.NewResponseRecorder(caddyhttp.ResponseWriterWrapper{}, buf, shouldBuf)
	replacer := caddy.NewReplacer()
	ctx := context.WithValue(r.Context(), caddy.ReplacerCtxKey, replacer)
	*r = *r.WithContext(ctx)

	l.ServeHTTP(rr, r, caddyhttp.HandlerFunc(next))

	if got, want := get(replacer, "geoip_country_code"), "CY"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_code} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_country_name"), "Cyprus"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_name} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_country_eu"), "false"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_eu} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_city_name"), "Limassol"; got != want {
		t.Errorf("Expected custom placeholder {geoip_city_name} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_latitude"), "34.684100"; got != want {
		t.Errorf("Expected custom placeholder {geoip_latitude} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_longitude"), "33.037900"; got != want {
		t.Errorf("Expected custom placeholder {geoip_longitude} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_time_zone"), "Asia/Nicosia"; got != want {
		t.Errorf("Expected custom placeholder {geoip_time_zone} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_geohash"), "swpmrf13wbgg"; got != want {
		t.Errorf("Expected custom placeholder {geoip_geohash} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_city_geoname_id"), "146384"; got != want {
		t.Errorf("Expected custom placeholder {geoip_city_geoname_id} to be set (%s), but it wasn't; got: %s", want, got)
	}

	if got, want := get(replacer, "geoip_country_geoname_id"), "146669"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_geoname_id} to be set (%s), but it wasn't; got: %s", want, got)
	}

	//
	// Verify that a request via the loopback interface address results in
	// the expected placeholder values.
	//
	var loopback_placeholders = [][2]string{
		{"geoip_country_code", "**"},
		{"geoip_country_name", "Loopback"},
		{"geoip_city_name", "Loopback"},
		{"geoip_country_geoname_id", "0"},
		{"geoip_city_geoname_id", "0"},
		{"geoip_latitude", "0.000000"},
		{"geoip_longitude", "0.000000"},
		{"geoip_geohash", "s00000000000"},
		{"geoip_time_zone", ""},
	}

	r = httptest.NewRequest("GET", "/", strings.NewReader(""))
	r.RemoteAddr = "127.0.0.1"

	rr = caddyhttp.NewResponseRecorder(caddyhttp.ResponseWriterWrapper{}, buf, shouldBuf)
	replacer = caddy.NewReplacer()
	ctx = context.WithValue(r.Context(), caddy.ReplacerCtxKey, replacer)
	*r = *r.WithContext(ctx)

	l.ServeHTTP(rr, r, caddyhttp.HandlerFunc(next))

	for _, v := range loopback_placeholders {
		if got, want := get(replacer, v[0]), v[1]; got != want {
			t.Errorf("Expected custom placeholder %s to be set (%s), but it wasn't; got: %s", v[0], want, got)
		}
	}

	//
	// Verify that a request via a private address results in the expected
	// placeholder values. Note that the MaxMind DB doesn't include
	// location data for private addresses.
	//
	var private_addr_placeholders = [][2]string{
		{"geoip_country_code", "!!"},
		{"geoip_country_name", "No Country"},
		{"geoip_city_name", "No City"},
		{"geoip_country_geoname_id", "0"},
		{"geoip_city_geoname_id", "0"},
		{"geoip_latitude", "0.000000"},
		{"geoip_longitude", "0.000000"},
		{"geoip_geohash", "s00000000000"},
		{"geoip_time_zone", ""},
	}

	r = httptest.NewRequest("GET", "/", strings.NewReader(""))
	r.RemoteAddr = "192.168.0.1"

	rr = caddyhttp.NewResponseRecorder(caddyhttp.ResponseWriterWrapper{}, buf, shouldBuf)
	replacer = caddy.NewReplacer()
	ctx = context.WithValue(r.Context(), caddy.ReplacerCtxKey, replacer)
	*r = *r.WithContext(ctx)

	l.ServeHTTP(rr, r, caddyhttp.HandlerFunc(next))
	if got, want := get(replacer, "geoip_country_code"), "!!"; got != want {
		t.Errorf("Expected custom placeholder {geoip_country_code} to be set (%s), but it wasn't; got: %s", want, got)
	}

	for _, v := range private_addr_placeholders {
		if got, want := get(replacer, v[0]), v[1]; got != want {
			t.Errorf("Expected custom placeholder %s to be set (%s), but it wasn't; got: %s", v[0], want, got)
		}
	}
}
