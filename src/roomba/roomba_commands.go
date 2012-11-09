// iRobot roomba open interface
package roomba

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

func MakeRoomba(port_name string) (*Roomba, error) {
	roomba := &Roomba{PortName: port_name, StreamPaused: make(chan bool, 1)}
	baud := uint(115200)
	err := roomba.Open(baud)
	return roomba, err
}

func (this *Roomba) Start() error {
	return this.WriteByte(OpCodes["Start"])
}

func (this *Roomba) Passive() error {
	return this.Start()
}

func (this *Roomba) Safe() error {
	return this.WriteByte(OpCodes["Safe"])
}

// Note: Use the Start command (128) to change the mode to Passive.
func (this *Roomba) Full() error {
	return this.WriteByte(OpCodes["Full"])
}

func (this *Roomba) Control() error {
	this.Passive()
	return this.WriteByte(130) // ?
}

func (this *Roomba) Clean() error {
	return this.WriteByte(OpCodes["Clean"])
}

func (this *Roomba) Spot() error {
	return this.WriteByte(OpCodes["Spot"])
}

func (this *Roomba) SeekDock() error {
	return this.WriteByte(OpCodes["SeekDock"])
}

func (this *Roomba) Power() error {
	return this.WriteByte(OpCodes["Power"])
}

func (this *Roomba) Drive(velocity, radius int16) error {
	if !(-500 <= velocity && velocity <= 500) {
		return fmt.Errorf("invalid velocity: %d", velocity)
	}
	if !(-2000 <= radius && radius <= 2000) {
		fmt.Errorf("invalid readius: %d", radius)
	}
	return this.Write(OpCodes["Drive"], pack([]interface{}{velocity, radius}))
}

func (this *Roomba) Stop() error {
	return this.Drive(0, 0)
}

func (this *Roomba) DirectDrive(right, left int16) error {
	if !(-500 <= right && right <= 500) ||
		!(-500 <= left && left <= 500) {
		return fmt.Errorf("invalid velocity. one of %d or %d", right, left)
	}
	return this.Write(OpCodes["DirectDrive"], pack([]interface{}{right, left}))
}

func (this *Roomba) LEDs(check_robot, dock, spot, debris bool, power_color, power_intensity byte) error {
	var led_bits byte

	for _, bit := range []bool{check_robot, dock, spot, debris} {
		led_bits <<= 1
		led_bits |= to_byte(bit)
	}
	return this.Write(OpCodes["LEDs"], pack([]interface{}{
		led_bits, power_color, power_intensity}))
}

func (this *Roomba) Sensors(packet_id byte) ([]byte, error) {
	this.Write(OpCodes["Sensor"], []byte{packet_id})
	bytes_to_read := SENSOR_PACKET_LENGTH[packet_id]
	result := make([]byte, bytes_to_read)
	n, err := this.Read(result)
	if byte(n) != bytes_to_read || err != nil {
		return result, fmt.Errorf("failed reading sensors data for packet id %d", packet_id)
	}
	return result, nil
}

func (this *Roomba) QueryList(packet_ids []byte) ([][]byte, error) {
	b := new(bytes.Buffer)
	b.WriteByte(byte(len(packet_ids)))
	b.Write(packet_ids)
	this.Write(OpCodes["QueryList"], b.Bytes())

	result := make([][]byte, len(packet_ids))
	for i, packet_id := range packet_ids {
		bytes_to_read := SENSOR_PACKET_LENGTH[packet_id]
		result[i] = make([]byte, bytes_to_read)
		n, err := this.Read(result[i])
		if byte(n) != bytes_to_read || err != nil {
			return result, fmt.Errorf("failed reading sensors data for packet id %d", packet_id)
		}
	}
	return result, nil
}

func (this *Roomba) PauseStream() {
	this.StreamPaused <- true
}

func (this *Roomba) ReadStream(packet_ids []byte, out chan<- [][]byte) {

	var data_length byte = 0
	for _, packet_id := range packet_ids {
		data_length += SENSOR_PACKET_LENGTH[packet_id]
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
			bytes_read := 0
			for bytes_read < len(buf) {
				n, err := this.S.Read(buf[bytes_read:])
				if n != 0 {
					bytes_read += n
				}
				if err != nil {
					if err == io.EOF {
						log.Printf("eof in stream: %s", err)
						return
					}
					goto Loop
				}
			}
			buf_r := bytes.NewReader(buf)
			if b, err := buf_r.ReadByte(); err != nil || b != 19 {
				log.Fatalf("stream data doesn't start with header 19")
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
				bytes_to_read := int(SENSOR_PACKET_LENGTH[packet_id])
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
