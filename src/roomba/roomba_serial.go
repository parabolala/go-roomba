// Low-level roomba interaction entities.

package roomba

import (
        "fmt"
        "errors"
        "log"
        "bytes"
        "github.com/tarm/goserial"
        "encoding/binary"
        "io"
    )
const (
    // Getting started commands
    OPCODE_START byte = 128
    OPCODE_BAUD byte = 129

    // Mode commands
    OPCODE_SAFE byte = 131
    OPCODE_FULL byte = 132

    // Cleaning commands
    OPCODE_CLEAN byte = 135
    OPCODE_MAX byte = 136
    OPCODE_SPOT byte = 134

    OPCODE_SEEK_DOCK byte = 143
    OPCODE_SCHEDULE byte = 167
    OPCODE_SET_DAY_TIME byte = 168
    OPCODE_POWER byte = 133

    // Actuator commands 
    OPCODE_DRIVE byte = 137
    OPCODE_DIRECT_DRIVE byte = 145
    OPCODE_DRIVE_PWM byte = 146
    OPCODE_MOTORS byte = 138
    OPCODE_PWM_MOTORS byte = 144
    OPCODE_LEDS byte = 139
    //OPCODE_SCHEDULING_LEDS byte = 162
    //OPCODE_DIGITAL_LEDS_RAW byte = 163
    //OPCODE_DIGITAL_LEDS_ASCII byte = 164
    //OPCODE_BUTTONS byte = 165
    OPCODE_SONG byte = 140
    OPCODE_PLAY byte = 141

    // Input commands
    OPCODE_SENSORS byte = 142
    OPCODE_QUERY_LIST byte = 149
    OPCODE_STREAM byte = 148
    OPCODE_PAUSE_STREAM byte = 150
)

const (
    SENSOR_BUMP_WHEELS_DROPS = 7
    SENSOR_WALL = 8
    SENSOR_CLIFF_LEFT = 9
    SENSOR_CLIFF_FRONT_LEFT = 10
    SENSOR_CLIFF_FRONT_RIGHT = 11
    SENSOR_CLIFF_RIGHT = 12
    SENSOR_VIRTUAL_WALL = 13
    SENSOR_WHEEL_OVERCURRENT = 14
    SENSOR_DIRT_DETECT = 15
    //unused = 16
    SENSOR_IR_OMNI = 17
    SENSOR_IR_LEFT = 52
    SENSOR_IR_RIGHT = 53
    SENSOR_BUTTONS = 18
    SENSOR_DISTANCE = 19
    SENSOR_ANGLE = 20
    SENSOR_CHARGING = 21
    SENSOR_VOLTAGE = 22
    SENSOR_CURRENT = 23
    SENSOR_TEMPERATURE = 24
    SENSOR_BATTERY_CHARGE = 25
    SENSOR_BATTERY_CAPACITY = 26
    SENSOR_WALL_SIGNAL = 27
    SENSOR_CLIFF_LEFT_SIGNAL = 28
    SENSOR_CLIFF_FRONT_LEFT_SIGNAL = 29
    SENSOR_CLIFF_FRONT_RIGHT_SIGNAL = 30
    SENSOR_CLIFF_RIGHT_SIGNAL = 31
    //unused = 32-33
    SENSOR_CHARGING_SOURCE = 34
    SENSOR_OI_MODE = 35
    SENSOR_SONG_NUMBER = 36
    SENSOR_SONG_PLAYING = 37
    SENSOR_NUM_STREAM_PACKETS = 38
    SENSOR_REQUESTED_VELOCITY = 39
    SENSOR_REQUESTED_RADIUS = 40
    //....
)

const WHEEL_SEPARATION = 298  // mm

type Roomba struct {
    PortName string
    S io.ReadWriteCloser
}

func pack(data []interface{}) []byte {
    buf := new(bytes.Buffer)
    for _, v := range data {
        err := binary.Write(buf, binary.BigEndian, v)
        if err != nil {
            log.Fatal("binary.Write failed:", err)
        }
    }
    return buf.Bytes()
}

func (this *Roomba) Open(baud uint) error {
    if baud != 115200 && baud != 19200 {
        return errors.New(fmt.Sprintf("Invalid baud rate: %u. Must be one of 115200, 19200", baud))
    }

    c := &serial.Config{Name: this.PortName, Baud: int(baud)}
    port, err := serial.OpenPort(c)

    if err != nil {
        log.Printf("Failed to open serial port: %s", this.PortName)
        return err
    }
    this.S = port
    log.Printf("Opened serial port: %s", this.PortName)
    return nil
}

func (this *Roomba) Write(p []byte) (n int, err error) {
    return this.S.Write(p)
}

func (this *Roomba) Read(p []byte) (n int, err error) {
    return this.S.Read(p)
}

func (this *Roomba) ChangeBaudRate(baud uint) error {
    this.S.Close()
    return this.Open(baud)
}
