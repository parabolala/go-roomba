// Package roomba defines and implements the interface for interacting with
// iRobot Roomba Open Interface robots.
package roomba

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"gobot.io/x/gobot/drivers/gpio"

	rb "github.com/deepakkamesh/go-roomba/constants"
)

type Roomba struct {
	PortName     string
	S            io.ReadWriter
	StreamPaused chan bool
	brcPin       *gpio.DirectPinDriver
	keepAlive    bool
}

// MakeRoomba initializes a new Roomba structure and sets up a serial port.
// By default, Roomba communicates at 115200 baud. Providing a brc port will
// send periodic pulses to keep roomba alive in passive mode.
func MakeRoomba(port_name string, brc *gpio.DirectPinDriver) (*Roomba, error) {

	roomba := &Roomba{
		PortName:     port_name,
		StreamPaused: make(chan bool, 1),
		keepAlive:    false,
	}

	baud := uint(115200)
	if err := roomba.Open(baud); err != nil {
		return nil, err
	}

	if brc == nil {
		return roomba, nil
	}

	// Setup BRC routine to keepalive.
	brc.DigitalWrite(1)
	roomba.brcPin = brc
	go roomba.PulseBRC()

	return roomba, nil
}

func (this *Roomba) PulseBRC() {
	for {
		time.Sleep(30 * time.Second)
		if !this.keepAlive {
			continue
		}

		// Pulse iRobot BRC port low for 1 second every 30 seconds to keep alive.
		this.brcPin.DigitalWrite(0)
		time.Sleep(1 * time.Second)
		this.brcPin.DigitalWrite(1)
		// When roomba is docked, it seems to go into sleep where pulsing BRC does
		// not wake it up. SeekDock button seems to keep it up.
		// TODO: only press this button when already docked and charging.
		this.ButtonPush(0x10)
	}
}

// Start command starts the OI. You must always send the Start command before
// sending any other commands to the OI.
// Note: Use the Start command (128) to change the mode to Passive.
func (this *Roomba) Start(keepAlive bool) error {
	this.keepAlive = keepAlive
	return this.WriteByte(rb.START)
}

// Passive switches Roomba to passive mode by sending the Start command.
func (this *Roomba) Passive() error {
	return this.Start(this.keepAlive)
}

func (this *Roomba) Reset() error {
	return this.WriteByte(rb.RESET)
}

// Safe sends the OI into Safe mode, enabling user control of Roomba.
// It turns off all LEDs.
func (this *Roomba) Safe() error {
	return this.WriteByte(rb.SAFE)
}

// Full command gives you complete control over Roomba by putting the OI into
// Full mode, and turning off the cliff, wheel-drop and internal charger safety
// features.
func (this *Roomba) Full() error {
	return this.WriteByte(rb.FULL)
}

// Control command's effect and usage are identical to the Safe command.
func (this *Roomba) Control() error {
	this.Passive()
	return this.WriteByte(rb.CONTROL)
}

// Clean command starts the default cleaning mode.
func (this *Roomba) Clean() error {
	return this.WriteByte(rb.CLEAN)
}

// Spot command starts the Spot cleaning mode.
func (this *Roomba) Spot() error {
	return this.WriteByte(rb.SPOT)
}

// SeekDock command sends Roomba to the dock.
func (this *Roomba) SeekDock() error {
	return this.WriteByte(rb.SEEK_DOCK)
}

// Power command powers down Roomba.
func (this *Roomba) Power() error {
	return this.WriteByte(rb.POWER)
}

// This command stops the OI. All streams will stop and the robot will no longer respond to commands.
// Use this command when you are finished working with the robot.
func (this *Roomba) Stop() error {
	return this.WriteByte(rb.STOP)
}

/*************************** ACTUATOR FUNCTIONS ********************************/

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
// straight = 32768 or 32767 = ex 8000 or 7FFF, turn in place clockwise = -1,
// turn in place counter-clockwise = 1
func (this *Roomba) Drive(velocity, radius int16) error {
	if !(-500 <= velocity && velocity <= 500) {
		return fmt.Errorf("invalid velocity: %d", velocity)
	}
	if !(-2000 <= radius && radius <= 2000) {
		fmt.Errorf("invalid readius: %d", radius)
	}
	return this.Write(rb.DRIVE, Pack([]interface{}{velocity, radius}))
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
	return this.Write(rb.DRIVE_DIRECT, Pack([]interface{}{right, left}))
}

// MainBrush controls the main brush motor.
func (this *Roomba) MainBrush(on bool, defaultDir bool) error {
	var cmd byte
	if !defaultDir {
		cmd = cmd | 16
	}

	if on {
		cmd = cmd | rb.MAIN_BRUSH
		return this.Write(rb.MOTORS, Pack([]interface{}{cmd}))
	}
	cmd = cmd
	return this.Write(rb.MOTORS, Pack([]interface{}{cmd}))
}

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
	return this.Write(rb.LEDS, Pack([]interface{}{
		led_bits, power_color, power_intensity}))
}

func (this *Roomba) ButtonPush(button byte) error {
	return this.Write(rb.BUTTONS, Pack([]interface{}{button}))
}

// LEDDisplay displays data in the 7 segment display.
func (this *Roomba) LEDDisplay(data []byte) error {
	if l := len(data); l != 4 {
		return fmt.Errorf("not enough digits. Expected 4 got %d ", l)
	}
	return this.Write(rb.DIGIT_LEDS_ASCII, Pack([]interface{}{data}))
}

func (this *Roomba) Song(num byte, notes ...byte) error {
	if num > 4 {
		return fmt.Errorf("only 4 songs are allowed")
	}
	l := len(notes)
	if l%2 != 0 || l/2 > 16 {
		return fmt.Errorf("there should be even number of notes and max of 16. ie. freq, dur")
	}

	return this.Write(rb.SONG, Pack([]interface{}{num, byte(l / 2), notes}))
}

// Play lets you select a song num to play from the songs added to Roomba using the Song command.
// You must add one or more songs to Roomba using the Song function in order for the Play command to
// work.
func (this *Roomba) Play(num byte) error {
	if num > 4 {
		return fmt.Errorf("only 4 songs are allowed")
	}

	return this.Write(rb.PLAY, Pack([]interface{}{num}))
}

/*************************** SENSOR FUNCTIONS ********************************/

// Sensors command requests the OI to send a packet of sensor data bytes. There
// are 58 different sensor data packets. Each provides a value of a specific
// sensor or group of sensors.
func (this *Roomba) Sensors(packet_id byte) ([]byte, error) {
	bytes_to_read, ok := rb.SENSOR_PACKET_LENGTH[packet_id]
	if !ok {
		return []byte{}, fmt.Errorf("unknown packet id requested: %d", packet_id)
	}

	this.Write(rb.SENSORS, []byte{packet_id})
	var err error
	var n int
	result := make([]byte, bytes_to_read)
	for byte(n) < bytes_to_read {
		result_view := result[n:]
		bytes_to_read -= byte(n)
		n, err = this.Read(result_view)
		if err != nil {
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
		_, ok := rb.SENSOR_PACKET_LENGTH[packet_id]
		if !ok {
			return [][]byte{}, fmt.Errorf("unknown packet id requested: %d", packet_id)
		}
	}

	b := new(bytes.Buffer)
	b.WriteByte(byte(len(packet_ids)))
	b.Write(packet_ids)
	this.Write(rb.QUERY_LIST, b.Bytes())

	var err error
	var n int
	result := make([][]byte, len(packet_ids))
	for i, packet_id := range packet_ids {
		bytes_to_read := rb.SENSOR_PACKET_LENGTH[packet_id]
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
		packet_length, ok := rb.SENSOR_PACKET_LENGTH[packet_id]
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
			this.Write(rb.PAUSE_RESUME_STREAM, []byte{0})
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
				bytes_to_read := int(rb.SENSOR_PACKET_LENGTH[packet_id])
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
	err := this.Write(rb.STREAM, b.Bytes())
	if err != nil {
		return nil, err
	}

	out := make(chan [][]byte)
	go this.ReadStream(packet_ids, out)
	return out, nil
}

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
