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
	"math"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/ThomasHabets/tleservice/pkg/proto"
	"github.com/joshuaferrara/go-satellite"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 10000, "The server port")
	user = flag.String("user", "", "Username for spacetrack")
	pass = flag.String("password", "", "Password for spacetrack") // TODO: put in env or something.
)

func ymd(ts time.Time) (int, int, int, int, int, int) {
	ts = ts.UTC()
	return ts.Year(), int(ts.Month()), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()
}

func gstimeFromDate(ts time.Time) float64 {
	return satellite.GSTimeFromDate(ymd(ts))
}

func fixLong(in float64) float64 {
	for in < 0 {
		in += 2 * math.Pi
	}
	for in > 2*math.Pi {
		in -= 2 * math.Pi
	}
	return rad2deg(in)
}

func deg2rad(in float64) float64 {
	return in / (180.0 / math.Pi)
}

func rad2deg(in float64) float64 {
	return in * 180.0 / math.Pi
}

func modelEnumString(m pb.Model) (string, error) {
	switch m {
	case pb.Model_DEFAULT, pb.Model_WGS84:
		return "wgs84", nil
	case pb.Model_WGS72:
		return "wgs72", nil
	default:
		return "", fmt.Errorf("invalid model")
	}
}

type Server struct {
	pb.UnimplementedTLEServiceServer

	spacetrack *satellite.Spacetrack
}

func (s *Server) GetTLE(ctx context.Context, req *pb.GetTLERequest) (*pb.GetTLEResponse, error) {
	m, err := modelEnumString(req.Model)
	if err != nil {
		return nil, err
	}
	sat, err := s.spacetrack.GetTLE(uint64(req.CatId), time.Now(), m)
	if err != nil {
		return nil, err
	}
	return &pb.GetTLEResponse{
		Tle: &pb.TLE{
			Tle1: sat.Line1,
			Tle2: sat.Line2,
		},
	}, nil
}

func (*Server) GetInstant(ctx context.Context, req *pb.GetInstantRequest) (*pb.GetInstantResponse, error) {
	model, err := modelEnumString(req.Model)
	if err != nil {
		return nil, err
	}

	sat := satellite.TLEToSat(req.Tle.Tle1, req.Tle.Tle2, model)
	resp := pb.GetInstantResponse{}
	for _, timepoint := range req.Timestamp {
		ts := time.Unix(timepoint, 0)
		y, mm, d, h, m, s := ymd(ts)

		jday := satellite.JDay(y, mm, d, h, m, s)
		pos, vel := satellite.Propagate(sat, y, mm, d, h, m, s)

		gst := satellite.GSTimeFromDate(ymd(ts))
		ecef := satellite.ECIToECEF(pos, gst)

		alt, angularVel, latlong := satellite.ECIToLLA(pos, gst)
		lng := fixLong(latlong.Longitude)
		data := &pb.InstantData{
			Timestamp: timepoint,
			Lla: &pb.LLA{
				Latitude:  rad2deg(latlong.Latitude),
				Longitude: lng,
				LongitudeEw: func() float64 {
					if lng > 180 {
						return lng - 360
					}
					return lng
				}(),
				Altitude: alt,
			},
			Position: &pb.Vector3{
				X: pos.X,
				Y: pos.Y,
				Z: pos.Z,
			},
			PositionEcef: &pb.Vector3{
				X: ecef.X,
				Y: ecef.Y,
				Z: ecef.Z,
			},
			Velocity: &pb.Vector3{
				X: vel.X,
				Y: vel.Y,
				Z: vel.Z,
			},
			AngularVelocity: angularVel,
		}
		if req.Observer != nil {
			obsPos := satellite.LatLong{
				Latitude:  deg2rad(req.Observer.Latitude),
				Longitude: deg2rad(req.Observer.Longitude),
			}
			ang := satellite.ECIToLookAngles(pos, obsPos, req.Observer.Altitude, jday)
			data.LookAngles = &pb.LookAngles{
				Azimuth:   rad2deg(ang.Az),
				Elevation: rad2deg(ang.El),
				Range:     ang.Rg,
			}
		}
		resp.Instant = append(resp.Instant, data)
	}
	return &resp, nil
}

func main() {
	flag.Parse()

	srv := Server{
		spacetrack: satellite.NewSpacetrack(*user, *pass),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTLEServiceServer(grpcServer, &srv)
	log.Infof("Running…")
	grpcServer.Serve(lis)
}
