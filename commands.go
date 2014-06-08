// Package roomba iRobot roomba Open Interface.
//
// The Roomba OI has four operating modes: Off, Passive, Safe, and Full.
// When roomba starts the OI is in “off” mode. When it is off, the OI listens
// for an OI Start command. Once it receives the Start command, you can enter
// into any one of the four operating modes by sending a mode command to the OI.
//
// Passive mode: entered upon sending one of the cleaning commands. You can
// only read sensor data in the passive mode and can't change the actuators
// state.
//
// Safe mode: gives full control of Roomba, except for safety restrictions:
//   * Cliff detection when moving forward.
//	 * Detection of wheel drop.
// 	 * Charger plugged in and powered.
// When any of the events ocurs, Roomba switches to passive mode.
//
// Full mode: gives full control over Romoba, disabling the safety
// restrictions.

package roomba

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/xa4a/go-roomba/constants"
)

func to_byte(b bool) byte {
	var res byte
	switch b {
	case false:
		res = 0
	case true:
		res = 1
	}
	return res
}

var OpCodes = constants.OpCodes

// MakeRoomba initializes a new Roomba structure and sets up a serial port.
// By default, Roomba communicates at 115200 baud.
func MakeRoomba(port_name string) (*Roomba, error) {
	roomba := &Roomba{PortName: port_name, StreamPaused: make(chan bool, 1)}
	baud := uint(115200)
	err := roomba.Open(baud)
	return roomba, err
}

// Start command starts the OI. You must always send the Start command before
// sending any other commands to the OI.
// Note: Use the Start command (128) to change the mode to Passive.
func (this *Roomba) Start() error {
	return this.WriteByte(OpCodes["Start"])
}

// TODO: Baud command.

// Passive switches Roomba to passive mode by sending the Start command.
func (this *Roomba) Passive() error {
	return this.Start()
}

// This command puts the OI into Safe mode, enabling user control of Roomba.
// It turns off all LEDs.
func (this *Roomba) Safe() error {
	return this.WriteByte(OpCodes["Safe"])
}

// Full command gives you complete control over Roomba by putting the OI into
// Full mode, and turning off the cliff, wheel-drop and internal charger safety
// features.
func (this *Roomba) Full() error {
	return this.WriteByte(OpCodes["Full"])
}

// Control command's effect and usage are identical to the Safe command.
func (this *Roomba) Control() error {
	this.Passive()
	return this.WriteByte(130) // ?
}

// Clean command starts the default cleaning mode.
func (this *Roomba) Clean() error {
	return this.WriteByte(OpCodes["Clean"])
}

// TODO: Max command.

// Spot command starts the Spot cleaning mode.
func (this *Roomba) Spot() error {
	return this.WriteByte(OpCodes["Spot"])
}

// SeekDock command sends Roomba to the dock.
func (this *Roomba) SeekDock() error {
	return this.WriteByte(OpCodes["SeekDock"])
}

// TODO: Schedule, Set Day/Time.

// Power command powers down Roomba.
func (this *Roomba) Power() error {
	return this.WriteByte(OpCodes["Power"])
}

// Drive command controls Roomba’s drive wheels. It takes two 16-bit signed
// values. The first one specifies the average velocity of the drive wheels in
// millimeters per second (mm/s).  The next one specifies the radius in
// millimeters at which Roomba will turn. The longer radii make Roomba drive
// straighter, while the shorter radii make Roomba turn more. The radius is
// measured from the center of the turning circle to the center of Roomba. A
// Drive command with a positive velocity and a positive radius makes Roomba
// drive forward while turning toward the left. A negative radius makes Roomba
// turn toward the right. Special cases for the radius make Roomba turn in place
// or drive straight. A negative velocity makes Roomba drive backward. Velocity
// is in range (-500 – 500 mm/s), radius (-2000 – 2000 mm). Special cases:
// straight = 32768 or 32767 = hex 8000 or 7FFF, turn in place clockwise = -1,
// turn in place counter-clockwise = 1
func (this *Roomba) Drive(velocity, radius int16) error {
	if !(-500 <= velocity && velocity <= 500) {
		return fmt.Errorf("invalid velocity: %d", velocity)
	}
	if !(-2000 <= radius && radius <= 2000) {
		fmt.Errorf("invalid readius: %d", radius)
	}
	return this.Write(OpCodes["Drive"], Pack([]interface{}{velocity, radius}))
}

// Stop commands is equivalent to Drive(0, 0).
func (this *Roomba) Stop() error {
	return this.Drive(0, 0)
}

// DirectDrive command lets you control the forward and backward motion of
// Roomba’s drive wheels independently. It takes two 16-bit signed values.
// The first specifies the velocity of the right wheel in millimeters per second
// (mm/s), The next one specifies the velocity of the left wheel A positive
// velocity makes that wheel drive forward, while a negative velocity makes it
// drive backward. Right wheel velocity (-500 – 500 mm/s). Left wheel velocity
// (-500 – 500 mm/s).
func (this *Roomba) DirectDrive(right, left int16) error {
	if !(-500 <= right && right <= 500) ||
		!(-500 <= left && left <= 500) {
		return fmt.Errorf("invalid velocity. one of %d or %d", right, left)
	}
	return this.Write(OpCodes["DirectDrive"], Pack([]interface{}{right, left}))
}

// TODO: Drive PWM, Motors, PWM Motors commands.

// LEDs command controls the LEDs common to all models of Roomba 500. The
// Clean/Power LED is specified by two data bytes: one for the color and the
// other for the intensity. Color: 0 = green, 255 = red. Intermediate values are
// intermediate colors (orange, yellow, etc). Intensitiy: 0 = off, 255 = full
// intensity. Intermediate values are intermediate intensities.
func (this *Roomba) LEDs(check_robot, dock, spot, debris bool, power_color, power_intensity byte) error {
	var led_bits byte

	for _, bit := range []bool{check_robot, dock, spot, debris} {
		led_bits <<= 1
		led_bits |= to_byte(bit)
	}
	return this.Write(OpCodes["LEDs"], Pack([]interface{}{
		led_bits, power_color, power_intensity}))
}

// TODO: Scheduling LEDs, Digit LEDs ASCII, Buttons, Song, Play.

// Sensors command requests the OI to send a packet of sensor data bytes. There
// are 58 different sensor data packets. Each provides a value of a specific
// sensor or group of sensors.
func (this *Roomba) Sensors(packet_id byte) ([]byte, error) {
	bytes_to_read, ok := constants.SENSOR_PACKET_LENGTH[packet_id]
	if !ok {
		return []byte{}, fmt.Errorf("unknown packet id requested: %d", packet_id)
	}

	this.Write(OpCodes["Sensors"], []byte{packet_id})
	var err error
	var n int
	result := make([]byte, bytes_to_read)
	for byte(n) < bytes_to_read {
		result_view := result[n:]
		bytes_to_read -= byte(n)
		n, err = this.Read(result_view)
		if err != nil {
			log.Printf("error %v", err)
			return result, fmt.Errorf("failed reading sensors data for packet id %d: %s", packet_id, err)
		}
	}
	return result, nil
}

// QueryList command lets you ask for a list of sensor packets. The result is
// returned once, as in the Sensors command. The robot returns the packets in
/// the order you specify.
func (this *Roomba) QueryList(packet_ids []byte) ([][]byte, error) {
	for _, packet_id := range packet_ids {
		_, ok := constants.SENSOR_PACKET_LENGTH[packet_id]
		if !ok {
			return [][]byte{}, fmt.Errorf("unknown packet id requested: %d", packet_id)
		}
	}

	b := new(bytes.Buffer)
	b.WriteByte(byte(len(packet_ids)))
	b.Write(packet_ids)
	this.Write(OpCodes["QueryList"], b.Bytes())

	var err error
	var n int
	result := make([][]byte, len(packet_ids))
	for i, packet_id := range packet_ids {
		bytes_to_read := constants.SENSOR_PACKET_LENGTH[packet_id]
		result[i] = make([]byte, bytes_to_read)
		err, n = nil, 0
		for byte(n) < bytes_to_read {
			result_view := result[i][n:]
			bytes_to_read -= byte(n)
			n, err = this.Read(result_view)
			if err != nil {
				return result, fmt.Errorf("failed reading sensors data for packet id %d: %s", packet_id, err)
			}
		}
	}
	return result, nil
}

// PauseStream command lets you stop steam without clearing the list of
// requested packets.
func (this *Roomba) PauseStream() {
	this.StreamPaused <- true
}

func (this *Roomba) ReadStream(packet_ids []byte, out chan<- [][]byte) {
	var data_length byte
	for _, packet_id := range packet_ids {
		packet_length, ok := constants.SENSOR_PACKET_LENGTH[packet_id]
		if !ok {
			log.Printf("unknown packet id requested: %d", packet_id)
			return
		}
		data_length += packet_length
	}

	// Input buffer. 3 is for 19, N-bytes and checksum.
	buf := make([]byte, data_length+byte(len(packet_ids))+3)

	for {
	Loop:
		select {
		case <-this.StreamPaused:
			// Pause stream.
			this.Write(OpCodes["ResumeStream"], []byte{0})
			close(out)
			return
		default:
			// Read single stream frame.
			bytes_read := 0
			for bytes_read < len(buf) {
				n, err := this.S.Read(buf[bytes_read:])
				if n != 0 {
					bytes_read += n
				}
				if err != nil {
					if err == io.EOF {
						return
					}
					goto Loop
				}
			}
			// Process frame.
			buf_r := bytes.NewReader(buf)
			if b, err := buf_r.ReadByte(); err != nil || b != 19 {
				log.Fatalf("stream data doesn't start with header 19")
				return
			}
			if b, err := buf_r.ReadByte(); err != nil || b != byte(len(buf)-3) {
				log.Fatalf("invalid N-bytes: %d, expected %d.", buf[1],
					len(buf)-3)
			}

			result := make([][]byte, len(packet_ids))

			i := 0
			// Used for verifying checksum.
			sum := byte(len(buf) - 3) // N-bytes
			packet_id, err := buf_r.ReadByte()
			for ; err == nil; packet_id, err = buf_r.ReadByte() {
				sum += packet_id
				bytes_to_read := int(constants.SENSOR_PACKET_LENGTH[packet_id])
				bytes_read := 0
				result[i] = make([]byte, bytes_to_read)

				for bytes_to_read > 0 {
					n, err := buf_r.Read(result[i][bytes_read:])
					bytes_read += n
					bytes_to_read -= n
					if err != nil {
						log.Fatalf("error reading packet data")
					}
				}
				for _, b := range result[i] {
					sum += b
				}
				i += 1
				if buf_r.Len() == 1 {
					break
				}
			}

			expected_checksum, err := buf_r.ReadByte()
			if err != nil {
				log.Fatalf("missing checksum")
			}
			sum += expected_checksum
			if sum != 0 {
				log.Fatalf("computed checksum didn't match: %d", sum)
			}
			out <- result
		}
	}
}

// Stream command starts a stream of data packets. The list of packets
// requested is sent every 15 ms, which is the rate Roomba uses to update data.
// This method of requesting sensor data is best if you are controlling Roomba
// over a wireless network (which has poor real-time characteristics) with
// software running on a desktop computer.
func (this *Roomba) Stream(packet_ids []byte) (<-chan [][]byte, error) {
	b := new(bytes.Buffer)
	b.WriteByte(byte(len(packet_ids)))
	b.Write(packet_ids)
	err := this.Write(OpCodes["Stream"], b.Bytes())
	if err != nil {
		return nil, err
	}

	out := make(chan [][]byte)
	go this.ReadStream(packet_ids, out)
	return out, nil
}
