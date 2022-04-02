package main

import (
	shifter "gps-shifter"
	"log"
)

// Move to config file
const (
	portName      = "/dev/ttyUSB0"
	maxBufferSize = 512
	baudRate      = 9600
)

var startPoint = shifter.NewVector2D(56.947340, 24.123012, 0, 0)

func main() {
	receiver, err := shifter.NewReceiver(portName, maxBufferSize, baudRate)
	if err != nil {
		log.Fatal(err)
	}

	sh := shifter.NewShifter(receiver, startPoint)

	shiftedCh := make(chan shifter.Vector2D, 1000)

	go sh.Shift(shiftedCh, startPoint)

	for {
		select {
		case vec := <-shiftedCh:
			log.Println(vec.String())
		default:
			continue
		}
	}
}
