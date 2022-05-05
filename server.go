package jtisim

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	apb "github.com/nileshsimaria/jtimon/authentication"
	dialoutpb "github.com/nileshsimaria/jtimon/gnmi/dialout"
	gnmi "github.com/nileshsimaria/jtimon/gnmi/gnmi"
	gnmipb "github.com/nileshsimaria/jtimon/gnmi/gnmi"
	tpb "github.com/nileshsimaria/jtimon/telemetry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// server size compression
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

// JTISim is JTI Simulator
type JTISim struct {
	host                              string
	port                              int32
	random                            bool
	descDir                           string
	dialOut                           bool
	skipVerify                        bool
	deviceCert, deviceKey, serverName string
	CACert                            string
	SSLServerName                     string
}

// NewJTISim to create new jti simulator
func NewJTISim(host string, port int32, random bool, descDir string, dialOut, skipVerify bool, deviceCert, deviceKey string, CACert string, SSLServerName string, serverName string) *JTISim {
	return &JTISim{
		host:          host,
		port:          port,
		random:        random,
		descDir:       descDir,
		dialOut:       dialOut,
		skipVerify:    skipVerify,
		serverName:    serverName,
		deviceCert:    deviceCert,
		deviceKey:     deviceKey,
		CACert:        CACert,
		SSLServerName: SSLServerName,
	}
}

// Start the simulator
func (s *JTISim) Start() error {
	if !s.dialOut {
		if lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port)); err == nil {
			grpcServer := grpc.NewServer()
			authServer := &authServer{}

			apb.RegisterLoginServer(grpcServer, authServer)
			tpb.RegisterOpenConfigTelemetryServer(grpcServer, &server{s})
			gnmipb.RegisterGNMIServer(grpcServer, &server{s})

			grpcServer.Serve(lis)
		} else {
			return err
		}
	} else {
		dialOutToGnmiCollector(&server{s}, s.host, s.port)
	}
	return nil
}

type server struct {
	jtisim *JTISim
}
type authServer struct {
}

func (s *authServer) LoginCheck(ctx context.Context, req *apb.LoginRequest) (*apb.LoginReply, error) {
	// allow everyone
	rep := &apb.LoginReply{
		Result: true,
	}
	return rep, nil
}

func (s *server) TelemetrySubscribe(req *tpb.SubscriptionRequest, stream tpb.OpenConfigTelemetry_TelemetrySubscribeServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		log.Println("Client metadata:")
		log.Println(md)
	}

	// send metadata to client
	header := metadata.Pairs("jtisim", "yes")
	stream.SendHeader(header)

	plist := req.GetPathList()
	ch := make(chan *tpb.OpenConfigData)
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

func (s *server) CancelTelemetrySubscription(ctx context.Context, req *tpb.CancelSubscriptionRequest) (*tpb.CancelSubscriptionReply, error) {
	return nil, nil
}

func (s *server) GetTelemetrySubscriptions(ctx context.Context, req *tpb.GetSubscriptionsRequest) (*tpb.GetSubscriptionsReply, error) {
	return nil, nil
}

func (s *server) GetTelemetryOperationalState(ctx context.Context, req *tpb.GetOperationalStateRequest) (*tpb.GetOperationalStateReply, error) {
	return nil, nil
}

func (s *server) GetDataEncodings(ctx context.Context, req *tpb.DataEncodingRequest) (*tpb.DataEncodingReply, error) {
	return nil, nil
}

func (s *server) Capabilities(context.Context, *gnmipb.CapabilityRequest) (*gnmipb.CapabilityResponse, error) {
	return nil, nil
}
func (s *server) Get(context.Context, *gnmipb.GetRequest) (*gnmipb.GetResponse, error) {
	return nil, nil
}
func (s *server) Set(context.Context, *gnmipb.SetRequest) (*gnmipb.SetResponse, error) {
	return nil, nil
}
func (s *server) Subscribe(stream gnmipb.GNMI_SubscribeServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		log.Println("Client metadata:")
		log.Println(md)
	}

	// send metadata to client
	header := metadata.Pairs("jtisim", "yes")
	stream.SendHeader(header)

	req, err := stream.Recv()
	if err != nil {
		log.Fatalf("Recv failed: %v", err)
	}

	subReq := req.GetSubscribe()
	if subReq == nil {
		log.Fatalf("Invalid subscribe request, received %v", req.GetRequest())
	}

	if subReq.GetEncoding() != gnmipb.Encoding_PROTO {
		log.Fatalf("Only PROTO encoding supported, received %v", subReq.GetEncoding())
	}

	if subReq.GetMode() != gnmipb.SubscriptionList_STREAM {
		log.Fatalf("Only STREAM mode supported, received %v", subReq.GetMode())
	}

	if subReq.GetUseAliases() {
		log.Fatalf("Aliases not supported, received %v", subReq.GetUseAliases())
	}

	stream.Send(&gnmipb.SubscribeResponse{Response: &gnmipb.SubscribeResponse_SyncResponse{SyncResponse: true}})
	ch := make(chan *gnmipb.SubscribeResponse)
	for _, sub := range subReq.GetSubscription() {
		sub.GetSampleInterval()
		gnmiPath := sub.GetPath()
		pname, _, _ := gnmiParsePath("", gnmiPath.GetElem(), nil, nil)
		switch {
		case strings.HasPrefix(pname, "/interfaces"):
			go s.gnmiStreamInterfaces(ch, pname, sub)
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

func dialOutToGnmiCollector(s *server, host string, port int32) {
	var opts []grpc.DialOption
	if s.jtisim.skipVerify {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	} else {
		certificate, _ := tls.LoadX509KeyPair(s.jtisim.deviceCert, s.jtisim.deviceKey)
		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(s.jtisim.CACert)
		if err != nil {
			log.Fatalf("[%s] failed to read ca cert: %s", host, err)
		}

		if ok := certPool.AppendCertsFromPEM(bs); !ok {
			log.Fatalf("[%s] failed to append certs", host)
		}

		transportCreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			ServerName:   s.jtisim.SSLServerName,
			RootCAs:      certPool,
		})
		opts = append(opts, grpc.WithTransportCredentials(transportCreds))
	}
	hostname := host + ":" + strconv.Itoa(int(port))

	var conn *grpc.ClientConn
	err := errors.New(fmt.Sprintf("[%s] Could not dial: %v", host, nil))
	for err != nil {
		conn, err = grpc.Dial(hostname, opts...)
		if nil == err {
			break
		}

		log.Printf(fmt.Sprintf("[%s] could not dial: %v", host, err))
		time.Sleep(10 * time.Second)
	}

	var ctx context.Context
	if s.jtisim.skipVerify {
		md := metadata.New(map[string]string{"server": s.jtisim.serverName})
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	} else {
		ctx = context.Background()
	}
	dialOutSubHandle, err := dialoutpb.NewSubscriberClient(conn).DialOutSubscriber(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("gNMI host: %v, subscribe handle creation failed, err: %v", hostname, err))
		return
	}

	err = dialOutSubHandle.Send(&gnmi.SubscribeResponse{})
	if err != nil {
		log.Fatalf(fmt.Sprintf("gNMI host: %v, send request failed: %v", hostname, err))
		return
	}

	var subReq *gnmi.SubscribeRequest
	for {
		log.Printf("gNMI host: %v, Waiting for susbcription list", hostname)
		subReq, err = dialOutSubHandle.Recv()
		if err == io.EOF {
			log.Fatalf(fmt.Sprintf("gNMI host: %v, received eof", hostname))
		}

		if err != nil {
			log.Fatalf(fmt.Sprintf("gNMI host: %v, received error: %v", hostname, err))
		}

		break
	}

	subscriptionList := subReq.GetSubscribe()
	if subscriptionList == nil {
		log.Fatalf("gNMI host: %v, Invalid subscribe request, received %v", hostname, subReq.GetRequest())
	}

	for {
		ch := make(chan *gnmipb.SubscribeResponse)

		for _, sub := range subscriptionList.GetSubscription() {
			sub.GetSampleInterval()
			gnmiPath := sub.GetPath()
			pname, _, _ := gnmiParsePath("", gnmiPath.GetElem(), nil, nil)
			switch {
			case strings.HasPrefix(pname, "/interfaces"):
				go s.gnmiStreamInterfaces(ch, pname, sub)
			default:
				log.Printf("gNMI host: %v, Sensor (%s) is not yet supported", hostname, pname)
			}
		}

		for {
			select {
			case data := <-ch:
				log.Println("%s", data)
				if err := dialOutSubHandle.Send(data); err != nil {
					log.Fatalf("[gNMI host: %v, Error sending response: %v", hostname, err)
				}
			}
		}
	}
}
