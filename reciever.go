package gps_shifter

import (
	"github.com/adrianmo/go-nmea"
	"go.bug.st/serial"
	"log"
)

// Move to config file
const portName = "/dev/ttyUSB0"
const maxBufferSize = 512

func RecieveNMEAData() string {
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
	}

	start := nmea.GLL{}

	buf := make([]byte, 100)
	phrase := make([]byte, 0, 512)
	sum := 0
	for {
		if sum > maxBufferSize {
			stop, err := GetLocation(string(phrase))
			if err != nil {
				// Change error behavior
				log.Println(err)
				sum = 0
				phrase = make([]byte, 0, 512)
				continue
			}

			delta := GetDelta(start, stop)

			log.Println(delta.String())

			start = stop

			sum = 0
			phrase = make([]byte, 0, 512)
		}
		n, err := port.Read(buf)
		if err != nil {
			// Change error behavior
			log.Fatal(err)
		}

		//log.Println(n)

		if n < len(buf) {
			buf = buf[:n]
		}

		sum += n

		phrase = append(phrase, buf...)

	}
}
