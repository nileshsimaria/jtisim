package main

import (
	js "github.com/nileshsimaria/jtisim"
	flag "github.com/spf13/pflag"
	"log"
)

var (
	host          = flag.String("host", "127.0.0.1", "host name or ip")
	port          = flag.Int32("port", 50051, "grpc server port")
	random        = flag.Bool("random", false, "Use random number to generate counter values")
	desc          = flag.String("description-dir", "../desc", "description directory")
	versionOnly   = flag.Bool("version", false, "Print version and build-time of the binary and exit")
	dialOut       = flag.Bool("dial-out", false, "Dialout")
	skipVerify    = flag.Bool("skip-verify", false, "Skip verify")
	serverName    = flag.String("server", "0.0.0.0", "Server name for inserting into metadata for skip-verify")
	CACert        = flag.String("ca-cert", "./certs/self_signed/ca-cert.pem", "Path of CA cert")
	deviceCert    = flag.String("cert", "./certs/self_signed/client-cert.pem", "Path of server cert")
	deviceKey     = flag.String("pem", "./certs/self_signed/client-key.pem", "Path of server key")
	SSLServerName = flag.String("ssl-server-name", "jcloud_demo.juniper.net", "SSL Server name as per cert")
	staticDialout = flag.Bool("static-dial-out", false, "Static Dial Out")

	jtisimVersion = "version-not-available"
	buildTime     = "build-time-not-available"
)

func main() {
	flag.Parse()

	log.Printf("Version: %s BuildTime %s\n", jtisimVersion, buildTime)
	if *versionOnly {
		return
	}
	if *staticDialout {
		*dialOut = true
	}

	jtisim := js.NewJTISim(*host, *port, *random, *desc, *dialOut, *staticDialout, *skipVerify, *deviceCert, *deviceKey, *CACert, *SSLServerName, *serverName)
	if err := jtisim.Start(); err != nil {
		log.Printf("can not start jti simulator: %v", err)
	}
}
