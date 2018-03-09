package geoip

import (
	"reflect"
	"testing"

	"github.com/mholt/caddy"
)

func TestParseConfig(t *testing.T) {
	controller := caddy.NewTestController("http", `
		localhost:8080
		geoip {
			set_header_country_code "Code"
			set_header_country_name "CountryName"
			set_header_country_eu "Eu"
			set_header_city_name "CityName"
			set_header_location_lat "Lat"
			set_header_location_lon "Lon"
			set_header_location_tz "TZ"
		}
	`)
	actual, err := parseConfig(controller)
	if err != nil {
		t.Errorf("parseConfig return err: %v", err)
	}
	expected := Config{
		HeaderNameCountryCode:      "Code",
		HeaderNameCountryName:      "CountryName",
		HeaderNameCountryIsEU:      "Eu",
		HeaderNameCityName:         "CityName",
		HeaderNameLocationLat:      "Lat",
		HeaderNameLocationLon:      "Lon",
		HeaderNameLocationTimeZone: "TZ",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v actual %v", expected, actual)
	}
}
