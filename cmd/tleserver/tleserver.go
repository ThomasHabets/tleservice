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
)

func ymd(ts time.Time) (int, int, int, int, int, int) {
	ts = ts.UTC()
	y, mm, d := ts.Date()
	h, m, s := ts.Hour(), ts.Minute(), ts.Second()
	return y, int(mm), d, h, m, s
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

type Server struct {
	pb.UnimplementedTLEServiceServer
}

func (*Server) GetInstant(ctx context.Context, req *pb.GetInstantRequest) (*pb.GetInstantResponse, error) {
	var model string
	switch req.Model {
	case pb.Model_DEFAULT, pb.Model_WGS84:
		model = "wgs84"
	case pb.Model_WGS72:
		model = "wgs72"
	default:
		return nil, fmt.Errorf("invalid model")
	}
	sat := satellite.TLEToSat(req.Tle.Tle1, req.Tle.Tle2, model)
	ts := time.Unix(req.Timestamp, 0)
	y, mm, d, h, m, s := ymd(ts)

	jday := satellite.JDay(y, mm, d, h, m, s)
	pos, vel := satellite.Propagate(sat, y, mm, d, h, m, s)

	gmst := satellite.GSTimeFromDate(ymd(ts))
	ecef := satellite.ECIToECEF(pos, gmst)

	alt, angularVel, latlong := satellite.ECIToLLA(pos, gmst)
	lng := fixLong(latlong.Longitude)
	resp := &pb.GetInstantResponse{
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
		fmt.Println(req.Observer.Latitude, req.Observer.Longitude)
		ang := satellite.ECIToLookAngles(pos, satellite.LatLong{
			//Latitude:  latlong.Latitude,
			Latitude: deg2rad(req.Observer.Latitude),
			//Longitude: latlong.Longitude,
			Longitude: deg2rad(req.Observer.Longitude),
		}, req.Observer.Altitude, jday-2451545.0)
		resp.LookAngles = &pb.LookAngles{
			Azimuth:   rad2deg(ang.Az),
			Elevation: rad2deg(ang.El),
			Range:     ang.Rg,
		}
	}

	return resp, nil
}

func main() {
	flag.Parse()

	srv := Server{}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTLEServiceServer(grpcServer, &srv)
	log.Infof("Running…")
	grpcServer.Serve(lis)
}
