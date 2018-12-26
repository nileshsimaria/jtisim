package jtisim

import (
	"fmt"
	"log"
	"net"
	"strings"

	auth_pb "github.com/nileshsimaria/jtisim/authentication"
	pb "github.com/nileshsimaria/jtisim/telemetry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// JTISim is JTI Simulator
type JTISim struct {
	host    string
	port    int32
	random  bool
	descDir string
}

// NewJTISim to create new jti simulator
func NewJTISim(host string, port int32, random bool, descDir string) *JTISim {
	return &JTISim{
		host:    host,
		port:    port,
		random:  random,
		descDir: descDir,
	}
}

// Start the simulator
func (s *JTISim) Start() error {
	if lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port)); err == nil {
		grpcServer := grpc.NewServer()
		authServer := &authServer{}

		auth_pb.RegisterLoginServer(grpcServer, authServer)
		pb.RegisterOpenConfigTelemetryServer(grpcServer, &server{s})

		grpcServer.Serve(lis)
	} else {
		return err
	}
	return nil
}

type server struct {
	jtisim *JTISim
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
		log.Println("Client metadata:")
		log.Println(md)
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
			go s.streamInterfaces(ch, path)
		case strings.HasPrefix(pname, "/bgp"):
			go s.streamBGP(ch, path)
		case strings.HasPrefix(pname, "/lldp"):
			go s.streamLLDP(ch, path)
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
