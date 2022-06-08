package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/joshuaferrara/go-satellite"
)

func main() {

	// Ground station.
	location := satellite.LatLong{
		Latitude:  51.4375 * math.Pi / 180,
		Longitude: 0.1250 * math.Pi / 180,
	}
	altitude := 48.0

	// Initial test. I no longer have what gpredict said for these.
	tle1 := "1 25544U 98067A   21290.96791059  .00007152  00000-0  13913-3 0  9995"
	tle2 := "2 25544  51.6432  95.6210 0004029 117.2302  27.3327 15.48732973307628"
	ts, err := time.Parse("2006-01-02 15:04:05", "2021-10-20 09:51:06")
	if err != nil {
		log.Fatal(err)
	}

	// Test case 1: gpredict says az 282.85, elevation 0.00.
	if false {
		// 2022-06-07 16:57
		tle1 = "1 25544U 98067A   22158.15063898  .00006400  00000+0  12044-3 0  9991"
		tle2 = "2 25544  51.6454  26.3008 0004489 203.5655 299.2429 15.49899681343622"

		ts, err = time.Parse("2006-01-02 15:04:05", "2022-06-07 17:20:30")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Test case 2. gpredict says az 215.26, elevation 22.24.
	if true {
		// 2022-06-07 16:57
		tle1 = "1 25544U 98067A   22158.15063898  .00006400  00000+0  12044-3 0  9991"
		tle2 = "2 25544  51.6454  26.3008 0004489 203.5655 299.2429 15.49899681343622"

		ts, err = time.Parse("2006-01-02 15:04:05", "2022-06-07 17:25:33")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Test case 3. gpredict says az 201.14, el 21.49.
	if false {
		// 2022-06-07 16:57
		tle1 = "1 25544U 98067A   22158.15063898  .00006400  00000+0  12044-3 0  9991"
		tle2 = "2 25544  51.6454  26.3008 0004489 203.5655 299.2429 15.49899681343622"

		ts, err = time.Parse("2006-01-02 15:04:05", "2022-06-07 17:26:03")
		if err != nil {
			log.Fatal(err)
		}
	}
	//ts = time.Now()
	//ts = ts.UTC()
	year, month, day, hour, minute, second := ts.Year(),
		int(ts.Month()),
		ts.Day(),
		ts.Hour(),
		ts.Minute(),
		ts.Second()
	// initialize satellite
	sat := satellite.TLEToSat(tle1, tle2, "wgs84")
	// get the satellite position
	position, _ := satellite.Propagate(
		sat, year, month, day,
		hour, minute, second,
	)
	// convert the current time to Galileo system time (GST)
	gst := satellite.GSTimeFromDate(
		year, month, day,
		hour, minute, second,
	)
	// get satellite coordinates in radians
	_, _, latlng := satellite.ECIToLLA(position, gst)
	// declare my current location, altitude
	// get my observation angles in radian
	if false {
		fmt.Println(position)
		fmt.Println(location)
		fmt.Println(altitude)
	}
	obs := satellite.ECIToLookAngles(
		position, location, altitude,
		// get Julian date
		satellite.JDay(
			year, month, day,
			hour, minute, second,
		),
	)
	// print satellite coordinates in angles
	fmt.Printf("Sat pos: %f %f\n", latlng.Latitude*180/math.Pi, latlng.Longitude*180/math.Pi)
	// print my observation azimuth in angle
	fmt.Printf("Azimuth:  %.2f\n", obs.Az*180/math.Pi)
	fmt.Printf("Elevation %.2f\n", obs.El*180/math.Pi)
}
