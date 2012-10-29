// iRobot roomba open interface

package roomba

func MakeRoomba(port_name string) (*Roomba, error) {
    roomba := &Roomba{PortName: port_name}
    baud := uint(115200)
    err := roomba.Open(baud)
    return roomba, err
}

func (this *Roomba) Start() error {
    _, err := this.S.Write(pack([]interface{}{OPCODE_START}))
    return err
}

func (this *Roomba) Passive() error {
    return this.Start()
}

func (this *Roomba) Control() error {
    this.Passive()
    _, err := this.Write(pack([]interface{}{byte(130)})) // ?
    return err
}

func (this *Roomba) Drive(velocity, radius int16) error {
    _, err := this.Write(pack([]interface{}{OPCODE_DRIVE, velocity, radius}))
    return err
}

func (this *Roomba) Stop() error {
    return this.Drive(0, 0)
}

