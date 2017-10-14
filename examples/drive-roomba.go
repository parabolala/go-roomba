package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
)

const (
	defaultPort = "/dev/tty.usbserial-DA01NYRU"
)

var (
	portName = flag.String("port", defaultPort, "roomba's serial port name")
)

func main() {
	flag.Parse()

	// Need GPIO control for BRC.
	pi := raspi.NewAdaptor()
	if err := pi.Connect(); err != nil {
		panic(err)
	}
	d := gpio.NewDirectPinDriver(pi, "7")
	if err := d.Start(); err != nil {
		panic(err)
	}

	// Enable pin
	en := gpio.NewDirectPinDriver(pi, "11")
	if err := en.Start(); err != nil {
		panic(err)
	}
	en.DigitalWrite(0)

	// Initialize Roomba.
	r, err := roomba.MakeRoomba(*portName, d)
	if err != nil {
		log.Fatalf("Making roomba failedi %v", err)
	}
	r.Start(false)
	r.Safe()

	in := bufio.NewReader(os.Stdin)

	for {

		b, err := in.ReadByte()
		if err != nil {
			log.Printf("Failed to read stdin: %v", err)
		}

		// Validate size of input

		switch b {
		case 'q':
			if err := r.LEDDisplay([]byte{'a', 'b', '1', '2'}); err != nil {
				panic(err)
			}
		case 'f':
			if err := r.Drive(50, -1); err != nil {
				panic(err)
			}
		case 'a':
			if err := r.Drive(50, 1); err != nil {
				panic(err)
			}
		case 's':
			if err := r.Drive(100, 32767); err != nil {
				panic(err)
			}
		case 'd':
			if err := r.Drive(-100, 32767); err != nil {
				panic(err)
			}
		case 'n':
			if err := r.MainBrush(false, true); err != nil {
				panic(err)
			}
		case 'y':
			if err := r.MainBrush(true, true); err != nil {
				panic(err)
			}

		case 'e':
			if err := en.DigitalWrite(1); err != nil {
				panic(err)
			}

		case 'r':
			if err := en.DigitalWrite(0); err != nil {
				panic(err)
			}

		case 'g':
			d, err := r.Sensors(constants.SENSOR_CURRENT)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Sensor Query Current: %v\n", int16(d[0])<<8|int16(d[1]))
			d, err = r.Sensors(constants.SENSOR_VOLTAGE)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Sensor Query voltage: %v\n", int16(d[0])<<8|int16(d[1]))
			d, err = r.Sensors(constants.SENSOR_OI_MODE)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Sensor Mode: %v\n", d[0])

		}

		time.Sleep(500 * time.Millisecond)
		r.Drive(0, 0)
		//r.Power()
	}
}
