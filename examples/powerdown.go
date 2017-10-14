package main

import (
	"flag"
	"fmt"
	"log"

	roomba "github.com/xa4a/go-roomba"
)

const (
	defaultPort = "/dev/ttyS0"
)

var (
	portName = flag.String("port", defaultPort, "roomba's serial port name")
	action   = flag.String("action", "", "safe,passive,stop,full,reset,start")
	brc      = flag.String("brc", "LCD-D23", "safe,passive,stop,full,reset,start")
)

func main() {
	flag.Parse()
	r, err := roomba.MakeRoomba(*portName, *brc)
	if err != nil {
		log.Fatal("Making roomba failed")
	}

	switch *action {
	case "start":
		r.Start(true)
	case "safe":
		r.Safe()
	case "stop":
		r.Stop()
	case "reset":
		r.Reset()
	case "full":
		r.Full()
	case "power":
		r.Power()
	case "pulse":
		r.PulseBRC()
	case "dock":
		r.ButtonPush(4)
	default:
		fmt.Println("no action specified")
	}

}
