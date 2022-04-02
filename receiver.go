package gps_shifter

import (
	"github.com/adrianmo/go-nmea"
	"go.bug.st/serial"
	"log"
)

// Move to config file
const portName = "/dev/ttyUSB0"
const maxBufferSize = 512

func ReceiveNMEAData() string {
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
	}

	start := nmea.GLL{}

	buf := make([]byte, 128)
	byteChannel := make(chan byte, 1000)
	gllPhrase := make([]byte, 0, maxBufferSize)

	go receiveGLLData(port, byteChannel, buf)

	for {
		phrase := getGLLData(gllPhrase, byteChannel)
		log.Println(phrase)
		loc, err := GetLocation(phrase)
		if err != nil {
			log.Println(err)
			continue
		}

		delta := GetDelta(start, loc)

		log.Println(delta.String())

		start = loc

	}
}

func receiveGLLData(port serial.Port, outCh chan byte, buffer []byte) {
	for {
		n, err := port.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		if n < len(buffer) {
			buffer = buffer[:n]
		}

		for _, b := range buffer {
			outCh <- b
		}
	}
}

func getGLLData(phraseBuffer []byte, bytesFromReceiver chan byte) string {
	gllStarted := false
	for b := range bytesFromReceiver {
		if !gllStarted {
			switch string(phraseBuffer) {
			case "":
				if b == '$' {
					phraseBuffer = append(phraseBuffer, b)
				}
			case "$":
				if b == 'G' {
					phraseBuffer = append(phraseBuffer, b)
				}
			case "$G":
				if b == 'P' {
					phraseBuffer = append(phraseBuffer, b)
				}
			case "$GP":
				if b == 'G' {
					phraseBuffer = append(phraseBuffer, b)
				}
			case "$GPG":
				if b == 'L' {
					phraseBuffer = append(phraseBuffer, b)
				}
			case "$GPGL":
				if b == 'L' {
					phraseBuffer = append(phraseBuffer, b)
				}
			case "$GPGLL":
				gllStarted = true
				phraseBuffer = append(phraseBuffer, b)
			}
			continue
		}
		if b == '\n' {
			phrase := string(phraseBuffer)
			phraseBuffer = make([]byte, 0, 512)
			gllStarted = false
			return phrase
		}
		phraseBuffer = append(phraseBuffer, b)
	}
	return ""
}
