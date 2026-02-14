package commons

import (
	"fmt"
	"log"
	"math"
	"strings"
	"strconv"

	"github.com/fatih/color"
)

const (
	percent = 100.0
)

func IntToString(integer int) string {
	return strconv.FormatInt(int64(integer), 10)
}

func Int64ToString(integer int64) string {
	return strconv.FormatInt(integer, 10)
}

func FloatToString(number float64) string {
	return strconv.FormatFloat(number, 'f', 2, 64)
}

func ParseInt(input string) (int, error) {
	output, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(output), nil
}

func ParseInt64(input string) (int64, error) {
	return strconv.ParseInt(input, 10, 64)
}

func MustParseInt(input string) int {
	output, err := ParseInt(input)
	if err != nil {
		log.Fatalf("Failed to parse integer: %s", input)
	}
	return output
}

func MustParseInt64(input string) int64 {
	output, err := ParseInt64(input)
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

func FormatPercentage(ratio float64, precision int) string {
	output := fmt.Sprintf("%+." + IntToString(precision) + "f%%", percent * ratio)
	if ratio > 0 {
		output = Green(output)
	} else if ratio < 0 {
		output = Red(output)
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
	fractionalString := fmt.Sprintf("%.2f", amount)
	fractionalPart := fractionalString[len(fractionalString) - 2:]
	output = fmt.Sprintf("%s.%s", output, fractionalPart)
	if amount < 0.0 {
		output = fmt.Sprintf("-%s", output)
	}
	return output
}

func FormatProfit(amount float64) string {
	formatted := FormatMoney(amount)
	if amount > 0 {
		formatted = Green(fmt.Sprintf("+%s", formatted))
	} else {
		formatted = Red(formatted)
	}
	return formatted
}

func Trim(input string) string {
	return strings.Trim(input, " \r\t\n")
}

func Green(input string) string {
	green := color.New(color.FgGreen).SprintFunc()
	return green(input)
}

func Red(input string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(input)
}