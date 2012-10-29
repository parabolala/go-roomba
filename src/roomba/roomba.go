// iRobot roomba open interface

package roomba

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

func (this *Roomba) Control() error {
	this.Passive()
	return this.Write0(130) // ?
}

func (this *Roomba) Drive(velocity, radius int16) error {
	_, err := this.Write(OpCodes["Drive"], pack([]interface{}{velocity, radius}))
	return err
}

func (this *Roomba) Stop() error {
	return this.Drive(0, 0)
}

func (this *Roomba) Clean() error {
	return this.Write0(OpCodes["Clean"])
}

func (this *Roomba) Spot() error {
	return this.Write0(OpCodes["Spot"])
}
