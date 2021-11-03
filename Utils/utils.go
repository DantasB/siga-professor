package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var layout = "02/01/2006"

func ToUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func ParseDate(date string) (time.Time, error) {
	splittedDate := strings.Split(date, " ")
	if len(splittedDate) != 4 {
		return time.Time{}, errors.New("could not split the string")
	}

	parsedTime, err := time.Parse(layout, splittedDate[2])
	return parsedTime, err
}

func RemoveSeparators(data string) string {
	return strings.Replace(strings.Replace(strings.Replace(data, "\t", " ", -1), "\r", " ", -1), "\n", " ", -1)
}

func ReplaceMultipleSpacesByPipe(data string) string {
	regex := regexp.MustCompile(` {6,}`)
	return regex.ReplaceAllString(data, "|")
}

func FillNilDataWithLastLineData(disciplines [][]string) ([][]string, error) {
	lastLine := []string{}
	for _, line := range disciplines {

		if line[0] != "\u00a0" {
			lastLine = line
		}

		for columnIndex, column := range line {
			if column == "\u00a0" {
				line[columnIndex] = lastLine[columnIndex]
			}
		}
	}

	return disciplines, nil
}
