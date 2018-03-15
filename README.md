[![Build Status](https://travis-ci.org/kodnaplakal/caddy-geoip.svg?branch=master)](https://travis-ci.org/kodnaplakal/caddy-geoip)
## Overview

`geoip` is a Caddy plugin that allow to determine user Geolocation by IP address using MaxMind database.

## Placeholders

The following placeholders are available:

```
  geoip_country_code - Country ISO code, example CY for Cyprus
  geoip_latitude - Latitude, example 34.684100
  geoip_longitude - Longitude, example 33.037900
  geoip_time_zone - Time zone, example Asia/Nicosia
  geoip_country_eu - Return 'true' if country in Europen Union
  geoip_country_name - Full country name
  geoip_city_name - City name
```


## Examples

(1) Set database path and return country code header:

```
geoip /path/to/db/GeoLite2-City.mmdb
header Country-Code {geoip_country_code}
```

(2) Proxy pass headers to backend:

```
localhost
geoip /path/to/db/GeoLite2-City.mmdb
proxy / localhost:3000 {
  header_upstream Country-Name {geoip_country_name}
  header_upstream Country-Code {geoip_country_code}
  header_upstream Country-Eu {geoip_country_eu}
  header_upstream City-Name {geoip_city_name}
  header_upstream Latitude {geoip_latitude}
  header_upstream Longitude {geoip_longitude}
  header_upstream Time-Zone {geoip_time_zone}
}
```

    
