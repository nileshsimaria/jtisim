package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	pb "github.com/nileshsimaria/jtimon/telemetry"
)

// IDesc Interface description structrue
type IDesc struct {
	Desc Description `json:"desc"`
	IFD  IFDCounters `json:"ifd-counters"`
	IFL  IFLCounters `json:"ifl-counters"`
}

// Description of interfaces
type Description struct {
	Media   string `json:"media"`
	FPC     int    `json:"fpc"`
	PIC     int    `json:"pic"`
	PORT    int    `json:"port"`
	Logical int    `json:"logical"`
}

// IFDCounters of interfaces
type IFDCounters struct {
	INPkts      uint64 `json:"in-pkts"`
	INOctets    uint64 `json:"in-octets"`
	AdminStatus bool   `json:"admin-status"`
	OperStatus  bool   `json:"oper-status"`
}

// IFLCounters of interfaces
type IFLCounters struct {
	INUnicastPkts   uint64 `json:"in-unicast-pkts"`
	INMulticastPkts uint64 `json:"in-multicast-pkts"`
}

func parseInterfacesJSON() *IDesc {
	file, err := ioutil.ReadFile("interfaces.json")
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	var iDesc IDesc
	if err := json.Unmarshal(file, &iDesc); err != nil {
		panic(err)
	}
	return &iDesc
}

type interfaces struct {
	ifds []*ifd
}
type ifd struct {
	name        string
	inPkts      uint64
	inOctets    uint64
	adminStatus string
	operStatus  string
	ifls        []*ifl
}

type ifl struct {
	index   int
	inUPkts uint64
	inMPkts uint64
}

func generateIList(idesc *IDesc) *interfaces {
	fpc := idesc.Desc.FPC
	pic := idesc.Desc.PIC
	port := idesc.Desc.PORT
	media := idesc.Desc.Media
	logical := idesc.Desc.Logical

	interfaces := &interfaces{
		ifds: make([]*ifd, fpc*pic*port),
	}

	cnt := 0
	for i := 0; i < fpc; i++ {
		for j := 0; j < pic; j++ {
			for k := 0; k < port; k++ {
				name := fmt.Sprintf("%s=%d/%d/%d", media, i, j, k)
				ifd := &ifd{
					name: name,
				}
				ifd.ifls = make([]*ifl, logical)

				for index := 0; index < logical; index++ {
					ifl := ifl{
						index: index,
					}
					ifd.ifls[index] = &ifl
				}

				interfaces.ifds[cnt] = ifd
				cnt++

			}
		}
	}
	return interfaces
}

func streamInterfaces(ch chan *pb.OpenConfigData, path *pb.Path) {
	seq := uint64(0)
	pname := path.GetPath()
	freq := path.GetSampleFrequency()
	fmt.Println(pname, freq)
	iDesc := parseInterfacesJSON()
	interfaces := generateIList(iDesc)

	for {
		ifds := interfaces.ifds
		for _, ifd := range ifds {
			prefixV := fmt.Sprintf("/interfaces/interface[name='%s']", ifd.name)
			inp := ifd.inPkts + 100 // TODO : introduce step concept
			ifd.inPkts = inp
			ino := ifd.inOctets + 1000 // TODO : introduce step concept
			ifd.inOctets = ino
			ops := "UP"
			ads := "DOWN"

			kv := []*pb.KeyValue{
				{Key: "__prefix__", Value: &pb.KeyValue_StrValue{StrValue: prefixV}},
				{Key: "name", Value: &pb.KeyValue_StrValue{StrValue: ifd.name}},
				{Key: "state/oper-status", Value: &pb.KeyValue_StrValue{StrValue: ops}},
				{Key: "state/admin-status", Value: &pb.KeyValue_StrValue{StrValue: ads}},
				{Key: "state/counters/in-pkts", Value: &pb.KeyValue_UintValue{UintValue: inp}},
				{Key: "state/counters/in-octets", Value: &pb.KeyValue_UintValue{UintValue: ino}},
			}

			d := &pb.OpenConfigData{
				SystemId:       "jvsim",
				ComponentId:    1,
				Timestamp:      uint64(MakeMSTimestamp()),
				SequenceNumber: seq,
				Kv:             kv,
			}
			seq++
			ch <- d

			for _, ifl := range ifd.ifls {
				prefixVifl := fmt.Sprintf("/interfaces/interface[name='%s']/subinterfaces/subinterface[index='%d']", ifd.name, ifl.index)
				inup := ifl.inUPkts + 100
				ifl.inUPkts = inup
				inmp := ifl.inMPkts + 1000
				ifl.inMPkts = inmp

				kvifl := []*pb.KeyValue{
					{Key: "__prefix__", Value: &pb.KeyValue_StrValue{StrValue: prefixVifl}},
					{Key: "subinterfaces/subinterface/index", Value: &pb.KeyValue_UintValue{UintValue: uint64(ifl.index)}},
					{Key: "state/counters/in-unicast-pkts", Value: &pb.KeyValue_UintValue{UintValue: inup}},
					{Key: "state/counters/in-multicast-pkts", Value: &pb.KeyValue_UintValue{UintValue: inmp}},
				}
				d := &pb.OpenConfigData{
					SystemId:       "jvsim",
					ComponentId:    1,
					Timestamp:      uint64(MakeMSTimestamp()),
					SequenceNumber: seq,
					Kv:             kvifl,
				}
				seq++
				ch <- d

			}

		} //finish one wrap

		time.Sleep(time.Duration(freq) * time.Millisecond)
	}
}
