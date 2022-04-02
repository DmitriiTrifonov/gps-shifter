package gps_shifter

import (
	"github.com/adrianmo/go-nmea"
	"go.bug.st/serial"
	"log"
)

type Receiver struct {
	portName   string
	bufferSize int
	baudRate   int
	port       serial.Port
}

func NewReceiver(portName string, bufferSize, baudRate int) (*Receiver, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	return &Receiver{
		portName:   portName,
		bufferSize: bufferSize,
		baudRate:   baudRate,
		port:       port,
	}, nil
}

func (r *Receiver) ReceiveNMEAData(out chan nmea.GLL) {
	buf := make([]byte, 128)
	byteChannel := make(chan byte, 1000)
	gllPhrase := make([]byte, 0, r.bufferSize)

	go receiveGLLData(r.port, byteChannel, buf)

	for {
		phrase := getGLLData(gllPhrase, byteChannel)
		loc, err := GetLocation(phrase)
		if err != nil {
			log.Println(err)
			continue
		}

		out <- loc
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
