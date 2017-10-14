package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	roomba "github.com/xa4a/go-roomba"
	c "github.com/xa4a/go-roomba/constants"
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

	// Check OI mode and switch to SAFE mode.
	var (
		mode []byte
		e    error
	)
	/*
		mode, e = r.Sensors(c.SENSOR_OI_MODE)
		if e != nil {
			fmt.Printf("Unable to change modes %v\n", e)
		}
		fmt.Printf("current OI Mode %d\n", mode[0])

		switch mode[0] {
		case 0: // If off turn on.
			r.Start()
		case 1, 3: // if Passive put in safe
			r.Safe()
		}
	*/
	r.Start()
	//r.Passive()
	//r.Safe()
	mode, e = r.Sensors(c.SENSOR_OI_MODE)
	if e != nil {
		fmt.Printf("Unable to query modes %v\n", e)
	}
	fmt.Printf("current OI Mode %d\n", mode[0])
	/*
	   mode, e = r.Sensors(c.SENSOR_OI_MODE)
	   	if e != nil || mode[0] != 2 {
	   		fmt.Printf("IO Mode should be SAFE expect:2, got: %d\n", mode[0])
	   	}
	*/
	// Update this list to demo different features.
	demoList := []byte{
		//c.DIGIT_LEDS_ASCII,
		//c.SONG,
		//	c.PLAY,
		c.SENSORS,
	}

	for _, demo := range demoList {
		switch demo {
		case c.DIGIT_LEDS_ASCII:
			if e := r.LEDDisplay([]byte{'b', 'c', '4', '2'}); e != nil {
				fmt.Println(e)
			}

		case c.SONG:
			if e := r.Song(1, 38, 32, 56, 32, 38, 32); e != nil {
				fmt.Println(e)
			}

		case c.PLAY:
			if e := r.Play(1); e != nil {
				fmt.Println(e)
			}

		case c.SENSORS:
			var (
				d []byte
				e error
				a int16
			)
			d, e = r.Sensors(c.SENSOR_VOLTAGE)
			if e != nil {
				fmt.Println(e)
				continue
			}
			a = int16(d[0])
			a = a<<8 | int16(d[1])
			fmt.Printf("Voltage %v\n", a)

			d, e = r.Sensors(c.SENSOR_CURRENT)
			if e != nil {
				fmt.Println(e)
				continue
			}
			a = int16(d[0])
			a = a<<8 | int16(d[1])
			fmt.Printf("Got Sensor Data %v\n", d)
			fmt.Printf("Current %v\n", a)

			d, e = r.Sensors(c.SENSOR_CHARGING)
			if e != nil {
				fmt.Println(e)
				continue
			}
			fmt.Printf("Charging State: %d\n", d[0])

			d, e = r.Sensors(c.SENSOR_BATTERY_CHARGE)
			if e != nil {
				fmt.Println(e)
				continue
			}
			a = int16(d[0])
			a = a<<8 | int16(d[1])
			fmt.Printf("Charg Left: %vmAh\n", a)

			d, e = r.Sensors(c.SENSOR_BATTERY_CAPACITY)
			if e != nil {
				fmt.Println(e)
				continue
			}
			a = int16(d[0])
			a = a<<8 | int16(d[1])
			fmt.Printf("Batt cap: %vmAh\n", a)

		}
		time.Sleep(200 * time.Millisecond)
	}
}
