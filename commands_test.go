//Tests for roomba package functions
package roomba_test

import (
	"testing"

	"github.com/xa4a/go-roomba/constants"
	rt "github.com/xa4a/go-roomba/testing"
)

func TestDrive(t *testing.T) {
	expected := []byte{137, 255, 56, 1, 244}
	r := rt.MakeTestRoomba()
	defer rt.ClearTestRoomba()
	r.Drive(-200, 500)
	rt.VerifyWritten(r, expected, t)
}

func TestLEDs(t *testing.T) {
	expected := []byte{139, 4, 0, 128}
	r := rt.MakeTestRoomba()
	defer rt.ClearTestRoomba()
	r.LEDs(false, true, false, false, 0, 128)
	rt.VerifyWritten(r, expected, t)
}

func TestQueryLists(t *testing.T) {
	output := []byte{3, 5}
	r := rt.MakeTestRoomba()
	defer rt.ClearTestRoomba()

	expected_input := []byte{149, 2, 7, 13}
	res, err := r.QueryList([]byte{
		constants.SENSOR_BUMP_WHEELS_DROPS,
		constants.SENSOR_VIRTUAL_WALL})
	if err != nil {
		t.Fatalf("error querying sensors: %s", err)
	}
	rt.VerifyWritten(r, expected_input, t)
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
	expected_data := [][]byte{{2, 25}, {5}}
	r := rt.MakeTestRoomba()
	defer rt.ClearTestRoomba()

	expected_input := []byte{148, 2, 29, 13}
	out, err := r.Stream([]byte{
		constants.SENSOR_CLIFF_FRONT_LEFT_SIGNAL,
		constants.SENSOR_VIRTUAL_WALL})
	if err != nil {
		t.Fatal("error querying senors")
	}
	response := <-out
	rt.VerifyWritten(r, expected_input, t)
	for i, packet_data := range response {
		for j, packet_byte := range packet_data {
			if expected_data[i][j] != packet_byte {
				t.Errorf("output byte doesn't match (byte %d in packet %d) expected %v != actual %v", j, i, expected_data[i], packet_data)
			}
		}
	}
}

func TestPauseStream(t *testing.T) {
	r := rt.MakeTestRoomba()
	defer rt.ClearTestRoomba()
	r.PauseStream()
	out, _ := r.Stream([]byte{})
	_, ok := <-out
	if ok {
		t.Fatalf("non-empty channel return by empty stream")
	}
	expected_input := []byte{148, 0, 150, 0}
	rt.VerifyWritten(r, expected_input, t)
}
