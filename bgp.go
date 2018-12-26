package jtisim

import (
	"log"
	"time"

	pb "github.com/nileshsimaria/jtisim/telemetry"
)

func (s *server) streamBGP(ch chan *pb.OpenConfigData, path *pb.Path) {
	pname := path.GetPath()
	freq := path.GetSampleFrequency()
	log.Println(pname, freq)

	seq := uint64(0)
	for {
		kv := []*pb.KeyValue{
			{Key: "__prefix__", Value: &pb.KeyValue_StrValue{StrValue: "/bgp"}},
			{Key: "state/foo", Value: &pb.KeyValue_UintValue{UintValue: 1111}},
		}

		d := &pb.OpenConfigData{
			SystemId:       "jtisim",
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
