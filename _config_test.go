package geoip

import (
	"reflect"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func _TestParseConfig(t *testing.T) {
	controller := caddyfile.New.NewTestController("http", `
		localhost:8080
		geoip path/to/maxmind/db
	`)
	actual, err := parseCaddyfile(controller)
	if err != nil {
		t.Errorf("parseConfig return err: %v", err)
	}
	expected := Config{
		DatabasePath: "path/to/maxmind/db",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v actual %v", expected, actual)
	}
}
