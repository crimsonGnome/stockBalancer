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

func ReverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
