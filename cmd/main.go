package main

import (
	"log"

	js "github.com/nileshsimaria/jtisim"

	flag "github.com/spf13/pflag"
)

var (
	host        = flag.String("host", "127.0.0.1", "host name or ip")
	port        = flag.Int32("port", 50051, "grpc server port")
	random      = flag.Bool("random", false, "Use random number to generate counter values")
	desc        = flag.String("description-dir", "../desc", "description directory")
	versionOnly = flag.Bool("version", false, "Print version and build-time of the binary and exit")

	jtisimVersion = "version-not-available"
	buildTime     = "build-time-not-available"
)

func main() {
	flag.Parse()

	log.Printf("Version: %s BuildTime %s\n", jtisimVersion, buildTime)
	if *versionOnly {
		return
	}

	jtisim := js.NewJTISim(*host, *port, *random, *desc)
	if err := jtisim.Start(); err != nil {
		log.Printf("can not start jti simulator: %v", err)
	}
}
