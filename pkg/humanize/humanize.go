package humanize

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ToBytes(humanSize string) int {
	var multiplier float64 = 1
	if strings.HasSuffix(humanSize, "K") {
		humanSize = strings.TrimSuffix(humanSize, "K")
		multiplier = 1024
	} else if strings.HasSuffix(humanSize, "M") {
		humanSize = strings.TrimSuffix(humanSize, "M")
		multiplier = 1024 * 1024
	} else if strings.HasSuffix(humanSize, "G") {
		humanSize = strings.TrimSuffix(humanSize, "G")
		multiplier = 1024 * 1024 * 1024
	}
	value, err := strconv.ParseFloat(humanSize, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(multiplier * value)

}

func FromBytes(byteSize int) string {
	if byteSize > 1024*1024*1024 {
		return fmt.Sprintf("%.2fG", float64(byteSize)/(1024*1024*1024))
	}
	if byteSize > 1024*1024 {
		return fmt.Sprintf("%.2fM", float64(byteSize)/(1024*1024))
	}
	if byteSize > 1024 {
		return fmt.Sprintf("%.2fK", float64(byteSize)/1024)
	}
	return strconv.Itoa(byteSize)
}
