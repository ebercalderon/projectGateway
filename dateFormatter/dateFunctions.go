package dateFormatter

import (
	"os"
	"time"
)

func GetStartOfDay(epoch int64) int64 {
	tz := os.Getenv("TIMEZONE")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("Error al cargar el location: " + err.Error())
	}
	fecha := time.UnixMilli(epoch)
	year, month, day := fecha.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc).UnixMilli()
}

func GetEndOfDay(epoch int64) int64 {
	tz := os.Getenv("TIMEZONE")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("Error al cargar el location: " + err.Error())
	}

	fecha := time.UnixMilli(epoch)
	year, month, day := fecha.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, loc).UnixMilli()
}
