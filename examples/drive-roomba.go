package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
)

const (
	defaultPort = "/dev/tty.usbserial-DA01NYRU"
)

var (
	portName = flag.String("port", defaultPort, "roomba's serial port name")
)

func main() {
	flag.Parse()
	r, err := roomba.MakeRoomba(*portName, "LCD-D22")
	if err != nil {
		log.Fatalf("Making roomba failedi %v", err)
	}
	r.Start(true)
	r.Safe()

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
		case 'n':
			r.MainBrush(false, true)
		case 'y':
			r.MainBrush(true, true)
		}

		time.Sleep(500 * time.Millisecond)
		r.Drive(0, 0)
		//r.Power()
	}
}
