package subsystem

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

// SizeToBytes convert string（100m、1G、512）to bytes
func SizeToBytes(memStr string) (int64, error) {
	if strings.TrimSpace(memStr) == constant.LiteralMax {
		return math.MaxInt64, nil
	}

	memStr = strings.TrimSpace(memStr)
	if memStr == "" {
		return 0, errors.New("value can't be empty")
	}

	var numStr string
	var unit string
	for i, c := range memStr {
		if !unicode.IsDigit(c) && c != '.' {
			numStr = memStr[:i]
			unit = strings.ToUpper(memStr[i:])
			break
		}
	}
	if numStr == "" {
		numStr = memStr
		unit = "B"
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, fmt.Errorf("memory can't be negative : [%s]", memStr)
	}

	var multiplier float64
	switch unit {
	case "B":
		multiplier = 1
	case "K", "KB":
		multiplier = 1024
	case "M", "MB":
		multiplier = 1024 * 1024
	case "G", "GB":
		multiplier = 1024 * 1024 * 1024
	case "T", "TB":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unsupported unit for [%s]（has to be one of : B/KB/MB/GB/TB）", unit)
	}

	bytes := int64(num * multiplier)
	return bytes, nil
}

func ParseMemoryToBits(memStr string) (int64, error) {
	bytes, err := SizeToBytes(memStr)
	if err != nil {
		return 0, err
	}
	return bytes * 8, nil
}
