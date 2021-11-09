# tleservice

TLEService is a microservice for calculating satellite orbits.

This is not an official Google product.

## Intro

[TLE][tle] is a data format for describing something that orbits the Earth.

It's then up to the user of the TLE data to turn this into "well, where is the
object right now?"

## Howto

Start the server:

```
go run ./cmd/tleserver/
```

Example client gets the lat/long/alt of the ISS in the next 10 minutes:

```
./tleclient  | awk '{print $1 " " $2 " " $3}' | head -10
-15.75 36.00 426.39
-12.76 38.30 425.58
-9.74 40.54 424.85
-6.71 42.74 424.19
-3.66 44.91 423.61
-0.61 47.07 423.12
2.44 49.22 422.72
5.50 51.38 422.42
8.54 53.57 422.20
11.56 55.79 422.08
```

(the TLE data for the ISS is currently hardcoded in the client)

## Use cases

* Schedule and aim communication with satellites, such as the amateur radio
  repeaters on the ISS.
* Know where in the sky to look, and when, to see the ISS go by.
* Make beautiful 2D and 3D plots of orbits.


[tle]: https://en.wikipedia.org/wiki/Two-line_element_set
