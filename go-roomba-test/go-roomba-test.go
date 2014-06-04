package main

import (
	"flag"
	"log"
	"time"

	"github.com/xa4a/go-roomba"
)

const (
	defaultPort = "/dev/cu.usbserial-FTTL3AW0"
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
	r.Safe()
	r.Drive(40, 200)
	t := time.Tick(1000 * time.Millisecond)
	<-t
	r.Stop()
}
