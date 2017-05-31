package main

import (
	"flag"
	"fmt"
	"log"

	roomba "github.com/xa4a/go-roomba"
)

const (
	defaultPort = "/dev/tty.usbserial-DA01NYRU"
)

var (
	portName = flag.String("port", defaultPort, "roomba's serial port name")
	action   = flag.String("action", "", "safe,passive,stop,full,reset,start")
)

func main() {
	flag.Parse()
	r, err := roomba.MakeRoomba(*portName)
	if err != nil {
		log.Fatal("Making roomba failed")
	}

	switch *action {
	case "start":
		r.Start()
	case "safe":
		r.Safe()
	case "stop":
		r.Stop()
	case "reset":
		r.Reset()
	case "full":
		r.Full()
	default:
		fmt.Println("no action specified")
	}

}
