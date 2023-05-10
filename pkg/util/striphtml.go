package util

import (
	"github.com/microcosm-cc/bluemonday"
	"html"
	"strings"
	"sync"
	"unicode/utf8"
)

const newLine = rune('\n')
const space = rune(' ')

var (
	p  *bluemonday.Policy
	mu sync.Mutex
)

func init() {
	mu.Lock()
	defer mu.Unlock()

	p = bluemonday.StrictPolicy()
}

func StripNewLine(s string, maxNewLine int) string {
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	sNewLine := 0
	sSpace := 0
	for _, c := range s {
		if c == newLine {
			sNewLine++
			if sNewLine <= maxNewLine {
				builder.WriteRune(c)
			}
			continue
		}
		if c == space {
			sSpace++
			if sSpace == 1 {
				builder.WriteRune(c)
			}
			continue
		}
		builder.WriteRune(c)
		sNewLine = 0
		sSpace = 0
	}

	return strings.TrimSpace(builder.String())
}

func StripHTMLTags(s string) string {
	mu.Lock()
	defer mu.Unlock()

	s = html.UnescapeString(s)
	s = p.Sanitize(s)
	s = html.UnescapeString(s)
	s = StripNewLine(s, 2)

	return strings.TrimSpace(s)
}
