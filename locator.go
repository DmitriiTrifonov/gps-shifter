package gps_shifter

import (
	"fmt"
	nmea "github.com/adrianmo/go-nmea"
)

func GetLocation(sentence string) (nmea.GLL, error) {
	parsed, err := nmea.Parse(sentence)
	if err != nil {
		return nmea.GLL{}, fmt.Errorf("cannot parse sentence: %s", sentence)
	}
	loc := parsed.(nmea.GLL)
	return loc, nil
}

func GetDelta(start, stop nmea.GLL) Vector2D {
	if start.Latitude == 0 || start.Longitude == 0 {
		return Vector2D{}
	}
	return NewVector2D(stop.Latitude, stop.Longitude,
		start.Latitude, start.Longitude)
}
