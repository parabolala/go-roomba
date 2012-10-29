package main

import (
	"log"
	"time"
	"roomba"
)

func main() {
	r, err := roomba.MakeRoomba("/dev/cu.usbserial-FTTL3AW0")
	if err != nil {
		log.Fatal("Making roomba failed")
	}
	r.Drive(40, 200)
	t := time.Tick(700 * time.Millisecond)
	<-t
	r.Stop()
}
