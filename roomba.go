// Package roomba defines and implements the interface for interacting with
// iRobot Roomba Open Interface robots.
package roomba

import (
	"io"
)

type Roomba struct {
	PortName     string
	S            io.ReadWriter
	StreamPaused chan bool
}
