package gps_shifter

import (
	"fmt"
	nmea "github.com/adrianmo/go-nmea"
	"strings"
)

func GetLocation(sentence string) (nmea.GLL, error) {
	sentences := strings.Split(sentence, "$")
	var gllStr string
	for _, s := range sentences {
		if strings.HasPrefix(s, "GPGLL") {
			gllStr = "$" + s
			break
		}
	}
	parsed, err := nmea.Parse(gllStr)
	if err != nil {
		return nmea.GLL{}, fmt.Errorf("cannot parse sentence: %s", sentence)
	}
	loc := parsed.(nmea.GLL)
	return loc, nil
}

func GetDelta(start, stop nmea.GLL) Vector2D {
	return NewVector2D(stop.Latitude, stop.Longitude,
		start.Latitude, start.Longitude)
}
