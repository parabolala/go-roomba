package main

import (
	"log"
	"roomba"
	"time"
)

func main() {
	r, err := roomba.MakeRoomba("/dev/cu.usbserial-FTTL3AW0")
    roomba.Start()
	if err != nil {
		log.Fatal("Making roomba failed")
	}
	r.Drive(40, 200)
	t := time.Tick(700 * time.Millisecond)
	<-t
	r.Stop()
}
