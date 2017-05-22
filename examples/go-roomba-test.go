package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xa4a/go-roomba"
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
	//r.Start()
	r.Safe()
	v, e := r.Sensors(31)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(v)
	//r.LEDs(false, false, false, false, 0, 0)
	r.Drive(100, 1)
	t := time.Tick(1000 * time.Millisecond)
	<-t
	r.Stop()
}
