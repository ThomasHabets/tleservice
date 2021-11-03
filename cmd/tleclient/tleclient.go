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
)

func main() {
	flag.Parse()

	//tle1 := "1 25544U 98067A   21290.96791059  .00007152  00000-0  13913-3 0  9995"
	//tle2 := "2 25544  51.6432  95.6210 0004029 117.2302  27.3327 15.48732973307628"
	tle1 := "1 25544U 98067A   21307.55056576  .00006301  00000-0  12281-3 0  9999"
	tle2 := "2 25544  51.6446  13.6218 0003585 168.4944 336.9654 15.48910556310190"
	// models: wgs72, wgs84, wgs72old

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTLEServiceClient(conn)
	ctx := context.Background()
	st := time.Now()
	if true {
		for i := 0; i < 180; i++ {
			ts := st.Add(time.Duration(i) * time.Minute)
			resp, err := client.GetInstant(ctx, &pb.GetInstantRequest{
				Tle: &pb.TLE{
					Tle1: tle1,
					Tle2: tle2,
				},
				Observer: &pb.LLA{
					Latitude:  51.76,
					Longitude: 0,
					Altitude:  0,
				},
				Timestamp: ts.Unix(),
			})
			if err != nil {
				log.Fatalf("Failed to RPC: %v", err)
			}
			if true {
				fmt.Printf("%.2f %.2f %.2f"+
					" %.2f %.2f %.2f"+
					" %.2f %.2f %.2f"+
					" %.2f %.2f %.2f"+
					" %.2f %.2f %.2f"+
					" %.2f"+
					"\n",
					resp.Lla.Latitude, resp.Lla.LongitudeEw, resp.Lla.Altitude,
					resp.Position.X, resp.Position.Y, resp.Position.Z,
					resp.PositionEcef.X, resp.PositionEcef.Y, resp.PositionEcef.Z,
					resp.Velocity.X, resp.Velocity.Y, resp.Velocity.Z,
					resp.LookAngles.Azimuth, resp.LookAngles.Elevation, resp.LookAngles.Range,
					resp.AngularVelocity,
				)
			} else {
				fmt.Printf("%.2f %.2f %.2f\n",
					resp.LookAngles.Azimuth, resp.LookAngles.Elevation, resp.LookAngles.Range,
				)
			}
		}
	} else {
		ts := time.Now()
		ts, err := time.Parse("2006-01-02 15:04:05", "2021-10-20 09:51:06")
		if err != nil {
			log.Fatalf("Failed to parse timestamp: %v", err)
		}
		resp, err := client.GetInstant(ctx, &pb.GetInstantRequest{
			Tle: &pb.TLE{
				Tle1: tle1,
				Tle2: tle2,
			},
			Timestamp: ts.Unix(),
			Observer: &pb.LLA{
				Latitude:  51.76,
				Longitude: 0,
				Altitude:  0,
			},
		})
		if err != nil {
			log.Fatalf("Failed to RPC: %v", err)
		}
		fmt.Printf("%.2f %.2f %.2f\n", resp.Lla.Latitude, resp.Lla.LongitudeEw, resp.Lla.Altitude)
		fmt.Printf("%.2f %.2f %.2f\n", resp.LookAngles.Azimuth, resp.LookAngles.Elevation, resp.LookAngles.Range)
	}
}
