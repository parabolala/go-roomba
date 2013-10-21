package roomba

import (
    "testing"
    "bytes"
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

		if len(actual) != len(expected) {
			t.Errorf("actual written length (%d) doesn't match expected (%d).",
				len(actual), len(expected))
			//return
		}
		for i, b := range expected {
			if b != actual[i] {
				t.Errorf("Expected output: % d, actual output: % d. Byte %d doesn't match",
					expected, actual, i)
			}
		}
	}
}
