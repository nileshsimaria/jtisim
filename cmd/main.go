package main

import (
	"log"

	js "github.com/nileshsimaria/jtisim"

	flag "github.com/spf13/pflag"
)

var (
	host   = flag.String("host", "127.0.0.1", "host name or ip")
	port   = flag.Int32("port", 50051, "grpc server port")
	random = flag.Bool("random", false, "Use random number to generate counter values")
	desc   = flag.String("description-dir", "../desc", "description directory")
)

func main() {
	flag.Parse()
	jtisim := js.NewJTISim(*host, *port, *random, *desc)
	if err := jtisim.Start(); err != nil {
		log.Printf("can not start jti simulator: %v", err)
	}
}
