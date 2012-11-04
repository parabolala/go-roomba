// iRobot roomba open interface
package roomba

import (
	"fmt"
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

// Drive PWM
// Motors
// PWM Motors

func (this *Roomba) LEDs(check_robot, dock, spot, debris bool, power_color, power_intensity byte) error {
	var led_bits byte

	for _, bit := range []bool{check_robot, dock, spot, debris} {
		led_bits <<= 1
		led_bits |= to_byte(bit)
	}
	return this.Write(OpCodes["LEDs"], pack([]interface{}{
		led_bits, power_color, power_intensity}))
}
