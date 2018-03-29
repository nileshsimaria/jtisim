package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	pb "github.com/nileshsimaria/jtimon/telemetry"
	flag "github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	host = flag.String("host", "127.0.0.1", "host name or ip")
	port = flag.Int32("port", 50051, "grpc server port")
)

func streamBGP(ch chan *pb.OpenConfigData, path *pb.Path) {
	pname := path.GetPath()
	freq := path.GetSampleFrequency()
	fmt.Println(pname, freq)

	seq := uint64(0)
	for {
		kv := []*pb.KeyValue{
			{Key: "__prefix__", Value: &pb.KeyValue_StrValue{StrValue: "/bgp"}},
			{Key: "state/foo", Value: &pb.KeyValue_UintValue{UintValue: 1111}},
		}

		d := &pb.OpenConfigData{
			SystemId:       "jvsim",
			ComponentId:    2,
			Timestamp:      uint64(MakeMSTimestamp()),
			SequenceNumber: seq,
			Kv:             kv,
		}
		ch <- d
		time.Sleep(time.Duration(freq) * time.Millisecond)
		seq++
	}
}

func streamLLDP(ch chan *pb.OpenConfigData, path *pb.Path) {
	pname := path.GetPath()
	freq := path.GetSampleFrequency()
	fmt.Println(pname, freq)

	seq := uint64(0)
	for {
		kv := []*pb.KeyValue{
			{Key: "__prefix__", Value: &pb.KeyValue_StrValue{StrValue: "/lldp"}},
			{Key: "state/foo", Value: &pb.KeyValue_UintValue{UintValue: 2222}},
		}

		d := &pb.OpenConfigData{
			SystemId:       "jvsim",
			ComponentId:    3,
			Timestamp:      uint64(MakeMSTimestamp()),
			SequenceNumber: seq,
			Kv:             kv,
		}
		ch <- d
		time.Sleep(time.Duration(freq) * time.Millisecond)
		seq++
	}
}

type server struct{}

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
			stream.Send(data)
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
	pb.RegisterOpenConfigTelemetryServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
