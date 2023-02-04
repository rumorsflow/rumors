package strutil

import (
	"strings"
	"unicode/utf8"
)

func MaxLen(str string, max int) string {
	if utf8.RuneCountInString(str) > max {
		i := 0
		for j := range str {
			if i == max {
				return strings.TrimRight(str[:j], "\n")
			}
			i++
		}
	}
	return str
}
