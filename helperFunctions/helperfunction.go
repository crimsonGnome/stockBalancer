package helperFunction

import (
	"log"
	"strconv"
)

func FloatConverter(number string) float64 {
	floatNumber, err := strconv.ParseFloat(number, 64)
	if err != nil {
		log.Fatalf("Failed to parse Amount, %v", err)
	}
	return floatNumber
}
