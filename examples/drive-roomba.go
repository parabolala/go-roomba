package main

import (
	"bufio"
	"flag"
	"log"
	"os"
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
	r.Start()
	r.Safe()
	//r.LEDs(false, false, false, false, 0, 0)
	//	r.Drive(100, 1)
	//	t := time.Tick(1000 * time.Millisecond)
	//	<-t
	//	r.Stop()
	in := bufio.NewReader(os.Stdin)

	for {

		b, err := in.ReadByte()
		if err != nil {
			log.Printf("Failed to read stdin: %v", err)
		}

		// Validate size of input

		switch b {
		case 'f':
			r.Drive(50, -1)
		case 'a':
			r.Drive(50, 1)
		case 's':
			r.Drive(100, 32767)
		case 'd':
			r.Drive(-100, 32767)
		}

		t := time.Tick(1000 * time.Millisecond)
		<-t
		r.Drive(0, 0)
		//r.Power()

	}

}
