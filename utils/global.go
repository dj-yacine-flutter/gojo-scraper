package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func CleanStringArray(inputStrings []string) []string {
	seenStrings := make(map[string]bool)
	outputStrings := []string{}

	for _, inputString := range inputStrings {
		inputString = strings.ToLower(inputString)

		if seenStrings[inputString] {
			continue
		}

		seenStrings[inputString] = true

		if strings.Contains(inputString, ",") {
			splitStrings := strings.Split(inputString, ",")
			outputStrings = append(outputStrings, splitStrings...)
		} else {
			outputStrings = append(outputStrings, inputString)
		}
	}

	return outputStrings
}

func CleanOverview(input string) string {
	pattern := regexp.MustCompile(`\n\n\[\]|\[[^\]]*]`)

	cleanedText := pattern.ReplaceAllString(input, "")

	return cleanedText
}

func CleanTag(input string) string {
	input = strings.ToLower(input)
	if strings.Contains(input, "maintenance") {
		return ""
	}
	if strings.Contains(input, "to episode") {
		return ""
	}
	if strings.Contains(input, "moved to") {
		return ""
	}
	if strings.Contains(input, "tag") {
		return ""
	}
	if strings.Contains(input, "element") {
		return ""
	}
	if strings.Contains(input, "setting") {
		return ""
	}
	if strings.Contains(input, "themes") {
		return ""
	}

	return input
}

func CleanRuntime(input string) string {
	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, "", "")
	input = strings.ReplaceAll(input, "hours", "h ")
	input = strings.ReplaceAll(input, "hour", "h ")
	input = strings.ReplaceAll(input, "hr", "h ")
	input = strings.ReplaceAll(input, "minutes", "m ")
	input = strings.ReplaceAll(input, "minute", "m ")
	input = strings.ReplaceAll(input, "seconds", "s ")
	input = strings.ReplaceAll(input, "second", "s ")
	return input
}

func CleanRating(input string) (string, error) {
	bracketPattern := regexp.MustCompile(`\([^)]*\)`)

	cleanedText := bracketPattern.ReplaceAllString(input, "")

	cleanedText = strings.TrimSpace(cleanedText)
	cleanedText = strings.ReplaceAll(cleanedText, " ", "")
	return cleanedText, nil
}

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
