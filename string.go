package commons

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

func IntToString(integer int64) string {
	return strconv.FormatInt(integer, 10)
}

func FloatToString(number float64) string {
	return strconv.FormatFloat(number, 'f', 2, 64)
}

func ParseInt(input string) (int64, error) {
	return strconv.ParseInt(input, 10, 64)
}

func MustParseInt(input string) int64 {
	output, err := ParseInt(input)
	if err != nil {
		log.Fatalf("Failed to parse integer: %s", input)
	}
	return output
}

func ParseFloat(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func MustParseFloat(input string) float64 {
	output, err := ParseFloat(input)
	if err != nil {
		log.Fatalf("Failed to convert string to float: %s", input)
	}
	return output
}

func FormatMoney(amount float64) string {
	amountString := fmt.Sprintf("%d", int64(math.Abs(amount)))
	output := "$"
	for i, character := range amountString {
		if i > 0 && (len(amountString) - i) % 3 == 0 {
			output += ","
		}
		output += string(character)
	}
	if amount < 0.0 {
		output = fmt.Sprintf("-%s", output)
	}
	return output
}