// Low-level roomba interaction entities.
package roomba

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tarm/goserial"
	"io"
	"log"
)

var OpCodes = map[string]byte{
	// Getting started commands
	"Start": 128,
	"Baud":  129,

	// Mode commands
	"Safe": 131,
	"Full": 132,

	// Cleaning commands
	"Clean": 135,
	"Max":   136,
	"Spot":  134,

	"Seek_dock":  143,
	"Schedule":   167,
	"SetDayTime": 168,
	"Power":      133,

	// Actuator commands 
	"Drive":       137,
	"DirectDrive": 145,
	"DrivePwm":    146,
	"Motors":      138,
	"PwmMotors":   144,
	"LEDs":        139,
	//SchedulingLeds: 162
	//DigitalLedsRaw: 163
	//DigitalLedsASCII: 164
	//Buttons: 165
	"Song": 140,
	"Play": 141,

	// Input commands
	"Sensors":     142,
	"QueryList":   149,
	"Stream":      148,
	"PauseStream": 150,
}

const (
	SENSOR_BUMP_WHEELS_DROPS = 7
	SENSOR_WALL              = 8
	SENSOR_CLIFF_LEFT        = 9
	SENSOR_CLIFF_FRONT_LEFT  = 10
	SENSOR_CLIFF_FRONT_RIGHT = 11
	SENSOR_CLIFF_RIGHT       = 12
	SENSOR_VIRTUAL_WALL      = 13
	SENSOR_WHEEL_OVERCURRENT = 14
	SENSOR_DIRT_DETECT       = 15
	//unused = 16
	SENSOR_IR_OMNI                  = 17
	SENSOR_IR_LEFT                  = 52
	SENSOR_IR_RIGHT                 = 53
	SENSOR_BUTTONS                  = 18
	SENSOR_DISTANCE                 = 19
	SENSOR_ANGLE                    = 20
	SENSOR_CHARGING                 = 21
	SENSOR_VOLTAGE                  = 22
	SENSOR_CURRENT                  = 23
	SENSOR_TEMPERATURE              = 24
	SENSOR_BATTERY_CHARGE           = 25
	SENSOR_BATTERY_CAPACITY         = 26
	SENSOR_WALL_SIGNAL              = 27
	SENSOR_CLIFF_LEFT_SIGNAL        = 28
	SENSOR_CLIFF_FRONT_LEFT_SIGNAL  = 29
	SENSOR_CLIFF_FRONT_RIGHT_SIGNAL = 30
	SENSOR_CLIFF_RIGHT_SIGNAL       = 31
	//unused = 32-33
	SENSOR_CHARGING_SOURCE    = 34
	SENSOR_OI_MODE            = 35
	SENSOR_SONG_NUMBER        = 36
	SENSOR_SONG_PLAYING       = 37
	SENSOR_NUM_STREAM_PACKETS = 38
	SENSOR_REQUESTED_VELOCITY = 39
	SENSOR_REQUESTED_RADIUS   = 40
	//....
	SENSOR_ALL = 100
)

var SENSOR_PACKET_LENGTH = map[byte]byte{
	SENSOR_BUMP_WHEELS_DROPS: 1,
	SENSOR_WALL:              1,
	SENSOR_CLIFF_LEFT:        1,
	SENSOR_CLIFF_FRONT_LEFT:  1,
	SENSOR_CLIFF_FRONT_RIGHT: 1,
	SENSOR_CLIFF_RIGHT:       1,
	SENSOR_VIRTUAL_WALL:      1,
	SENSOR_WHEEL_OVERCURRENT: 1,
	SENSOR_DIRT_DETECT:       1,
	//unused
	16:                              3,
	SENSOR_IR_OMNI:                  1,
	SENSOR_IR_LEFT:                  1,
	SENSOR_IR_RIGHT:                 1,
	SENSOR_BUTTONS:                  1,
	SENSOR_DISTANCE:                 2,
	SENSOR_ANGLE:                    2,
	SENSOR_CHARGING:                 1,
	SENSOR_VOLTAGE:                  2,
	SENSOR_CURRENT:                  2,
	SENSOR_TEMPERATURE:              1,
	SENSOR_BATTERY_CHARGE:           2,
	SENSOR_BATTERY_CAPACITY:         2,
	SENSOR_WALL_SIGNAL:              2,
	SENSOR_CLIFF_LEFT_SIGNAL:        2,
	SENSOR_CLIFF_FRONT_LEFT_SIGNAL:  2,
	SENSOR_CLIFF_FRONT_RIGHT_SIGNAL: 2,
	SENSOR_CLIFF_RIGHT_SIGNAL:       2,
	//unused
	32: 3,
	33: 3,
	SENSOR_CHARGING_SOURCE:    1,
	SENSOR_OI_MODE:            1,
	SENSOR_SONG_NUMBER:        1,
	SENSOR_SONG_PLAYING:       1,
	SENSOR_NUM_STREAM_PACKETS: 1,
	SENSOR_REQUESTED_VELOCITY: 2,
	SENSOR_REQUESTED_RADIUS:   2,
	//....
	// Group packets.
	0:          26,
	1:          10,
	2:          6,
	3:          10,
	4:          14,
	5:          12,
	6:          52,
	SENSOR_ALL: 100,
	101:        28,
	106:        12,
	107:        9,
}

const WHEEL_SEPARATION = 298 // mm

type Roomba struct {
	PortName string
	S        io.ReadWriteCloser
}

func pack(data []interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			log.Fatal("failed packing bytes:", err)
		}
	}
	return buf.Bytes()
}

func (this *Roomba) Open(baud uint) error {
	if baud != 115200 && baud != 19200 {
		return errors.New(fmt.Sprintf("invalid baud rate: %u. Must be one of 115200, 19200", baud))
	}

	c := &serial.Config{Name: this.PortName, Baud: int(baud)}
	port, err := serial.OpenPort(c)

	if err != nil {
		log.Printf("failed to open serial port: %s", this.PortName)
		return err
	}
	this.S = port
	log.Printf("opened serial port: %s", this.PortName)
	return nil
}

func (this *Roomba) Write(opcode byte, p []byte) error {
	n, err := this.S.Write([]byte{opcode})
	if n != 1 || err != nil {
		return fmt.Errorf("failed writing opcode %d to serial interface",
			opcode)
	}
	n, err = this.S.Write(p)
	if n != len(p) || err != nil {
		return fmt.Errorf("failed writing command to serial interface: % d", p)
	}
	return nil
}

func (this *Roomba) WriteByte(opcode byte) error {
	return this.Write(opcode, []byte{})
}

func (this *Roomba) Read(p []byte) (n int, err error) {
	return this.S.Read(p)
}
