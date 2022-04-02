package gps_shifter

import (
	"github.com/adrianmo/go-nmea"
)

type Shifter struct {
	receiver   *Receiver
	startPoint Vector2D
}

func NewShifter(receiver *Receiver, startPoint Vector2D) *Shifter {
	return &Shifter{receiver: receiver, startPoint: startPoint}
}

func (s *Shifter) Shift(out chan Vector2D, startPoint Vector2D) {
	gllChan := make(chan nmea.GLL, 1000)
	prevRealGLL := nmea.GLL{}

	go s.receiver.ReceiveNMEAData(gllChan)

	for {
		select {
		case gll := <-gllChan:
			// TODO: Add delta len check
			delta := GetDelta(prevRealGLL, gll)
			prevRealGLL = gll
			newPoint := Vector2D{
				startPoint.x + delta.x,
				startPoint.y + delta.y,
			}
			out <- newPoint
			startPoint = newPoint
		default:
			continue
		}
	}
}
