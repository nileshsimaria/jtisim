package main

import (
	"fmt"
	pb "github.com/nileshsimaria/jtimon/telemetry"
	flag "github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

var (
	host = flag.String("host", "127.0.0.1", "host name or ip")
	port = flag.Int32("port", 50051, "grpc server port")
)

type server struct{}

func (s *server) TelemetrySubscribe(req *pb.SubscriptionRequest, stream pb.OpenConfigTelemetry_TelemetrySubscribeServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Println("Received metadata")
		fmt.Println(md)
	}

	// send metadata to client
	header := metadata.Pairs("jtisim", "yes")
	stream.SendHeader(header)

	plist := req.GetPathList()
	for _, path := range plist {
		pname := path.GetPath()
		freq := path.GetSampleFrequency()
		fmt.Println(pname, freq)
	}

	seq := uint64(0)

	for {
		kv := []*pb.KeyValue{
			{Key: "__prefix__", Value: &pb.KeyValue_StrValue{StrValue: "/interfaces/interface[name='xe-1/2/0']/"}},
			{Key: "state/mtu", Value: &pb.KeyValue_UintValue{UintValue: 1514}},
		}

		d := &pb.OpenConfigData{
			SystemId:       "jvsim",
			ComponentId:    1212,
			Timestamp:      1510946604929,
			SequenceNumber: seq,
			Kv:             kv,
		}
		stream.Send(d)
		seq++
	}
	return nil
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
	pb.RegisterOpenConfigTelemetryServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
