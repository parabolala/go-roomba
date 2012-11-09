//Tests for roomba package functions
package roomba

import (
	"bytes"
	"testing"
)

type CloseableRWBuffer struct {
	r, w bytes.Buffer
}

func (this *CloseableRWBuffer) Read(p []byte) (int, error) {
	return this.r.Read(p)
}

func (this *CloseableRWBuffer) Write(p []byte) (int, error) {
	return this.w.Write(p)
}

func (this *CloseableRWBuffer) Close() error {
	return nil
}

func MakeTestRoomba() *Roomba {
	buffer := &CloseableRWBuffer{}
	r := &Roomba{S: buffer, StreamPaused: make(chan bool, 1)}
	return r
}

func VerifyWritten(r *Roomba, expected []byte, t *testing.T) {
	if buffer, ok := r.S.(*CloseableRWBuffer); ok {
		actual := buffer.w.Bytes()

		t.Logf("%d ", actual)
		if len(actual) != len(expected) {
			t.Errorf("actual written length (%d) doesn't match expected (%d).",
				len(actual), len(expected))
			return
		}
		for i, b := range expected {
			if b != actual[i] {
				t.Errorf("Expected output: % d, actual output: % d. Byte %d doesn't match",
					expected, actual, i)
			}
		}
	}
}

func TestDrive(t *testing.T) {
	expected := []byte{137, 255, 56, 1, 244}
	r := MakeTestRoomba()
	r.Drive(-200, 500)
	VerifyWritten(r, expected, t)
}

func TestLEDs(t *testing.T) {
	expected := []byte{139, 4, 0, 128}
	r := MakeTestRoomba()
	r.LEDs(false, true, false, false, 0, 128)
	VerifyWritten(r, expected, t)
}

func TestQueryLists(t *testing.T) {
	output := []byte{3, 5}
	r := MakeTestRoomba()
	r.S.(*CloseableRWBuffer).r.Write(output)

	expected_input := []byte{149, 2, 7, 13}
	res, err := r.QueryList([]byte{SENSOR_BUMP_WHEELS_DROPS,
		SENSOR_VIRTUAL_WALL})
	if err != nil {
		t.Fatal("error querying senors")
	}
	VerifyWritten(r, expected_input, t)
	for i, b := range res {
		if len(b) != 1 {
			t.Errorf("query_list returned wrong packet len for packet_id %d",
				expected_input[i+2])
		}
		if b[0] != output[i] {
			t.Errorf("query_list result %d doesn't match expected value", i)
		}
	}
}

func TestStream(t *testing.T) {
	expected_data := [][]byte{{2, 25}, {0}}
	output := []byte{19, 5, 29, 2, 25, 13, 0, 182}
	r := MakeTestRoomba()
	r.S.(*CloseableRWBuffer).r.Write(output)

	expected_input := []byte{148, 2, 29, 13}
	out, err := r.Stream([]byte{SENSOR_CLIFF_FRONT_LEFT_SIGNAL,
		SENSOR_VIRTUAL_WALL})
	if err != nil {
		t.Fatal("error querying senors")
	}
	response := <-out
	VerifyWritten(r, expected_input, t)
	for i, packet_data := range response {
		for j, packet_byte := range packet_data {
			if expected_data[i][j] != packet_byte {
				t.Error("output byte doesn't match")
			}
		}
	}
}

func TestPauseStream(t *testing.T) {
	output := []byte{19, 1, 235}
	r := MakeTestRoomba()
	r.S.(*CloseableRWBuffer).r.Write(output)
	r.PauseStream()
	out, _ := r.Stream([]byte{})
	_, ok := <-out
	if ok {
		t.Fatalf("non-empty channel return by empty stream")
	}
	expected_input := []byte{148, 0, 150, 0}
	VerifyWritten(r, expected_input, t)
}
