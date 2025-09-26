package commons

import (
	"fmt"
	"math"
	"strconv"
)

func IntToString(integer int64) string {
	return strconv.FormatInt(integer, 10)
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