//Tests for roomba package functions
package roomba

import (
	"bytes"
	"testing"
)

type CloseableBuffer struct {
	bytes.Buffer
}

func (this *CloseableBuffer) Close() error {
	return nil
}

func MakeTestRoomba() *Roomba {
	buffer := &CloseableBuffer{}
	r := &Roomba{S: buffer}
	return r
}

func VerifyOutput(r *Roomba, expected []byte, t *testing.T) {
	if buffer, ok := r.S.(*CloseableBuffer); ok {
		actual := buffer.Bytes()

		if len(actual) != len(expected) {
			t.Errorf("actual written length (%u) doesn't match expected (%u).",
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
	VerifyOutput(r, expected, t)
}

func TestLEDs(t *testing.T) {
	expected := []byte{139, 4, 0, 128}
	r := MakeTestRoomba()
	r.LEDs(false, true, false, false, 0, 128)
	VerifyOutput(r, expected, t)
}
