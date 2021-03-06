/*
      Copyright 2021 Google LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/ThomasHabets/tleservice/pkg/proto"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	startTime  = flag.String("time", "2021-10-19 12:07", "Start time to give data for.")
	getTLE     = flag.Bool("get_tle", false, "Get TLE data from the Internet. Off by default because it requires the server to have a login.")
	duration   = flag.Duration("duration", time.Hour, "Duration to get data for")
	period     = flag.Duration("period", time.Minute, "Data periodicity")
)

func printRange(ctx context.Context, client pb.TLEServiceClient, st time.Time, tle1, tle2 string) {
	var tss []int64
	for ts := st; ts.Before(st.Add(*duration)); ts = ts.Add(*period) {
		tss = append(tss, ts.Unix())
	}

	resps, err := client.GetInstant(ctx, &pb.GetInstantRequest{
		Tle: &pb.TLE{
			Tle1: tle1,
			Tle2: tle2,
		},
		Observer: &pb.LLA{
			Latitude:  51.4375,
			Longitude: 0.1250,
			Altitude:  48,
		},
		// Model: pb.Model_WGS72,
		Timestamp: tss,
	})
	if err != nil {
		log.Fatalf("Failed to RPC: %v", err)
	}

	for _, resp := range resps.Instant {
		fmt.Printf("%v deg=(%.2f %.2f) alt=%.2f"+
			" xyz=(%.2f %.2f %.2f)"+
			" ecef=(%.2f %.2f %.2f)"+
			" vel=(%.2f %.2f %.2f)"+
			" azimuth=%.2f elevation=%.2f range=%.2f"+
			" %.2f"+
			"\n",
			time.Unix(resp.Timestamp, 0),
			resp.Lla.Latitude, resp.Lla.LongitudeEw, resp.Lla.Altitude,
			resp.Position.X, resp.Position.Y, resp.Position.Z,
			resp.PositionEcef.X, resp.PositionEcef.Y, resp.PositionEcef.Z,
			resp.Velocity.X, resp.Velocity.Y, resp.Velocity.Z,
			resp.LookAngles.Azimuth, resp.LookAngles.Elevation, resp.LookAngles.Range,
			resp.AngularVelocity,
		)
	}
}

func main() {
	flag.Parse()

	// TLE data for the ISS. By the time you read this it'll be old, though.
	tle1 := "1 25544U 98067A   21307.55056576  .00006301  00000-0  12281-3 0  9999"
	tle2 := "2 25544  51.6446  13.6218 0003585 168.4944 336.9654 15.48910556310190"

	// 2022-06-07 16:57
	tle1 = "1 25544U 98067A   22158.15063898  .00006400  00000+0  12044-3 0  9991"
	tle2 = "2 25544  51.6454  26.3008 0004489 203.5655 299.2429 15.49899681343622"

	// Connect to server.
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewTLEServiceClient(conn)
	ctx := context.Background()

	if *getTLE {
		resp, err := client.GetTLE(ctx, &pb.GetTLERequest{CatId: 25544})
		if err != nil {
			log.Fatal(err)
		}
		tle1 = resp.Tle.Tle1
		tle2 = resp.Tle.Tle2
	}

	st, err := time.Parse("2006-01-02 15:04:05", *startTime)
	if err != nil {
		log.Fatalf("Failed to parse %q as time: %v", *startTime, err)
	}

	printRange(ctx, client, st, tle1, tle2)
}
