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
)

func main() {
	flag.Parse()
	r, err := roomba.MakeRoomba(*portName)
	if err != nil {
		log.Fatal("Making roomba failed")
	}

	if e := r.Power(); e != nil {
		fmt.Println(e)
	}

}
