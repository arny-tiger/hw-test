package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	leftPtr := 0
	rightPtr := 1
	var builder strings.Builder

	runeArr := []rune(str)
	for leftPtr < len(runeArr) {
		firstChar := runeArr[leftPtr]
		var secondChar rune
		if rightPtr < len(runeArr) {
			secondChar = runeArr[rightPtr]
		}
		if unicode.IsDigit(firstChar) {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(secondChar) {
			repeatsNum, _ := strconv.Atoi(string(secondChar))
			builder.WriteString(strings.Repeat(string(firstChar), repeatsNum))
			rightPtr += 2
			leftPtr += 2
		} else {
			builder.WriteString(string(firstChar))
			rightPtr++
			leftPtr++
		}
	}
	return builder.String(), nil
}
