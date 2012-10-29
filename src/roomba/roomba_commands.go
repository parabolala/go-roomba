// iRobot roomba open interface

package roomba

import (
	"log"
)

func MakeRoomba(port_name string) (*Roomba, error) {
	roomba := &Roomba{PortName: port_name}
	baud := uint(115200)
	err := roomba.Open(baud)
	return roomba, err
}

func (this *Roomba) Start() error {
	return this.Write0(OpCodes["Start"])
}

func (this *Roomba) Passive() error {
	return this.Start()
}

func (this *Roomba) Safe() error {
	return this.Write0(OpCodes["Safe"])
}

// Note: Use the Start command (128) to change the mode to Passive.
func (this *Roomba) Full() error {
	return this.Write0(OpCodes["Full"])
}

func (this *Roomba) Control() error {
	this.Passive()
	return this.Write0(130) // ?
}

func (this *Roomba) Clean() error {
	return this.Write0(OpCodes["Clean"])
}

func (this *Roomba) Spot() error {
	return this.Write0(OpCodes["Spot"])
}

func (this *Roomba) SeekDock() error {
	return this.Write0(OpCodes["SeekDock"])
}

func (this *Roomba) Power() error {
	return this.Write0(OpCodes["Power"])
}

func (this *Roomba) Drive(velocity, radius int16) error {
	if !(-500 <= velocity && velocity <= 500) {
		log.Fatalf("Invalid velocity: %d", velocity)
	}
	if !(-2000 <= radius && radius <= 2000) {
		log.Fatalf("Invalid readius: %d", radius)
	}
	_, err := this.Write(OpCodes["Drive"], pack([]interface{}{velocity, radius}))
	return err
}

func (this *Roomba) Stop() error {
	return this.Drive(0, 0)
}

func (this *Roomba) DirectDrive(right, left int16) error {
	if !(-500 <= right && right <= 500) ||
		!(-500 <= left && left <= 500) {
		log.Fatalf("Invalid velocity")
	}
	_, err := this.Write(OpCodes["DirectDrive"], pack([]interface{}{right, left}))
	return err
}

// Drive PWM
// Motors
// PWM Motors
