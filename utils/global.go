package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ExtractYear(input string) (int, error) {
	re := regexp.MustCompile(`(\d{4})`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return 0, fmt.Errorf("no four digits found")
	}
	year, err := strconv.ParseInt(matches[1], 0, 0)
	if err != nil {
		return 0, fmt.Errorf("cannot make four digits string a int")
	}
	return int(year), nil
}
