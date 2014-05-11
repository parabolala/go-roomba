package sim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"roomba/constants"
)

type RoombaSimulator struct {
	rw           io.ReadWriter
	writeQ       chan []byte
	WrittenBytes bytes.Buffer
	ReadBytes    bytes.Buffer
}

var MockSensorValues = map[byte][]byte{
	constants.SENSOR_BUMP_WHEELS_DROPS: []byte{3},
	constants.SENSOR_VIRTUAL_WALL:      []byte{5},
	constants.SENSOR_CLIFF_RIGHT:       []byte{42},
	constants.SENSOR_DISTANCE:          []byte{10, 20},
	constants.SENSOR_WALL:              []byte{35},
}

func (mock *RoombaSimulator) Serve() {
	// Write bytes from channel asynchronously.
	go func() {
		for {
			bs := <-mock.writeQ
			mock.WrittenBytes.Write(bs)
			mock.rw.Write(bs)
		}
	}()

	for {
		mock.ExecuteCMD()
	}
}

func (mock *RoombaSimulator) ExecuteCMD() error {
	cmdBuf := mock.read(1)
	if len(cmdBuf) != 1 {
		return fmt.Errorf("failed reading opcode")
	}
	switch cmdBuf[0] {
	case constants.OpCodes["Sensors"]:
		packetId := mock.read(1)[0]
		value, ok := MockSensorValues[packetId]
		if !ok {
			log.Printf("no mock value for sensor packet id %d", packetId)
		}
		log.Printf("sensor %d value: %v", packetId, value)
		mock.write(value)
	case constants.OpCodes["QueryList"]:
		nPackets := mock.read(1)[0]
		for i := 0; i < int(nPackets); i++ {
			packetId := mock.read(1)[0]
			value, ok := MockSensorValues[packetId]
			if !ok {
				log.Printf("no mock value for sensor packet id %d", packetId)
			}
			log.Printf("sensor %d value: %v", packetId, value)
			mock.write(value)
		}
	case constants.OpCodes["Stream"]:
		nBytes := mock.read(1)[0]
		if nBytes != 0 {
			output := []byte{19, 5, 29, 2, 25, 13, 0, 182}
			mock.write(output)
		}
	case constants.OpCodes["Start"]:
		log.Printf("switched to passive mode")
	case constants.OpCodes["Safe"]:
		log.Printf("switched to safe mode")
	case constants.OpCodes["DirectDrive"]:
		data := mock.read(4)
		var rigthVelocity, leftVelocity int16
		binary.Read(bytes.NewReader(data[:2]), binary.BigEndian, &rigthVelocity)
		binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &leftVelocity)
		log.Printf("DirectDrive: %d, %d (%v)", rigthVelocity, leftVelocity, data)
	default:
		log.Printf("unknown opcode: %d", cmdBuf[0])
	}

	return nil
}

// Reads given number of bytes from the Reader mock.r.
func (mock *RoombaSimulator) read(n int) []byte {
	buf := make([]byte, n)
	nRead, err := mock.rw.Read(buf)
	if n != nRead {
		if err != nil {
			log.Printf("error reading in RoombaSimulator: %v", err)
		}
		//log.Printf("read %d bytes when expected %d", nRead, n)
		return []byte{}
	}
	log.Printf("roomba reads: %v", buf)
	mock.ReadBytes.Write(buf)
	return buf
}

// Writes bytes to the Writer w asynchronously.
func (mock *RoombaSimulator) write(b []byte) {
	log.Printf("roomba says: %v", b)
	mock.writeQ <- b
}

// Helper for merging reader and writer into a ReadWriter.
type readWriter struct {
	io.Reader
	io.Writer
}

func MakeRoombaSim() (*RoombaSimulator, *readWriter) {
	// Input: driver writes, simulator reads.
	inp_r, inp_w := io.Pipe()

	// Ouput: simulator writes, driver reads.
	out_r, out_w := io.Pipe()

	readBytes := &bytes.Buffer{}
	writtenBytes := &bytes.Buffer{}

	mock := &RoombaSimulator{
		rw: &readWriter{
			// Log all read bytes to ReadBytes.
			io.TeeReader(inp_r, readBytes),
			// Log all written bytes to writtenBytes.
			io.MultiWriter(out_w, writtenBytes),
		},
		writeQ:    make(chan []byte, 15),
		ReadBytes: *readBytes,
	}
	go mock.Serve()

	rw := &readWriter{out_r, inp_w}

	return mock, rw
}
