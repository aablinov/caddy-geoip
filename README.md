[![Build Status](https://travis-ci.org/kodnaplakal/caddy-geoip.svg?branch=master)](https://travis-ci.org/kodnaplakal/caddy-geoip)
## Overview

`geoip` is a Caddy plugin that allow to determine user Geolocation by IP address using MaxMind database.

## Headers

`geoip` set this headers:

```
  X-Geoip-Country-Code - Country ISO code, example CY for Cyprus
  X-Geoip-Location-Lat - Latitude, example 34.684100
  X-Geoip-Location-Lon - Longitude, example 33.037900
  X-Geoip-Location-Tz - Time zone, example Asia/Nicosia
  X-Geoip-Country-Eu - Return 'true' if country in Europen Union
  X-Geoip-Country-Name - Full country name
  X-Geoip-City-Name - City name
```


## Examples

(1) Set database path:

```
geoip {
  database /path/to/db/GeoLite2-City.mmdb
}
```


(2) Set custom header names.

```
geoip {
  database path/to/maxmind/db
  set_header_country_code Code
  set_header_country_name CountryName
  set_header_country_eu Eu
  set_header_city_name CityName
  set_header_location_lat Lat
  set_header_location_lon Lon
  set_header_location_tz TZ
}
```

    
