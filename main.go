package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	auth_pb "github.com/nileshsimaria/jtisim/authentication"
	pb "github.com/nileshsimaria/jtisim/telemetry"
	flag "github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	host = flag.String("host", "", "host name or ip")
	port = flag.Int32("port", 50051, "grpc server port")
)

type server struct {
}
type authServer struct {
}

func (s *authServer) LoginCheck(ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginReply, error) {
	// allow everyone
	rep := &auth_pb.LoginReply{
		Result: true,
	}
	return rep, nil
}

func (s *server) TelemetrySubscribe(req *pb.SubscriptionRequest, stream pb.OpenConfigTelemetry_TelemetrySubscribeServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Println("Client metadata:")
		fmt.Println(md)
	}

	// send metadata to client
	header := metadata.Pairs("jtisim", "yes")
	stream.SendHeader(header)

	plist := req.GetPathList()
	ch := make(chan *pb.OpenConfigData)
	for _, path := range plist {
		pname := path.GetPath()
		switch {
		case strings.HasPrefix(pname, "/interfaces"):
			go streamInterfaces(ch, path)
		case strings.HasPrefix(pname, "/bgp"):
			go streamBGP(ch, path)
		case strings.HasPrefix(pname, "/lldp"):
			go streamLLDP(ch, path)
		default:
			log.Fatalf("Sensor (%s) is not yet supported", pname)
		}
	}

	for {
		select {
		case data := <-ch:
			if err := stream.Send(data); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func (s *server) CancelTelemetrySubscription(ctx context.Context, req *pb.CancelSubscriptionRequest) (*pb.CancelSubscriptionReply, error) {
	return nil, nil
}

func (s *server) GetTelemetrySubscriptions(ctx context.Context, req *pb.GetSubscriptionsRequest) (*pb.GetSubscriptionsReply, error) {
	return nil, nil
}

func (s *server) GetTelemetryOperationalState(ctx context.Context, req *pb.GetOperationalStateRequest) (*pb.GetOperationalStateReply, error) {
	return nil, nil
}

func (s *server) GetDataEncodings(ctx context.Context, req *pb.DataEncodingRequest) (*pb.DataEncodingReply, error) {
	return nil, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authServer := &authServer{}

	auth_pb.RegisterLoginServer(grpcServer, authServer)
	pb.RegisterOpenConfigTelemetryServer(grpcServer, &server{})

	grpcServer.Serve(lis)
}
