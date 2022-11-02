package str

import (
	"github.com/microcosm-cc/bluemonday"
	"html"
	"strings"
	"sync"
	"unicode/utf8"
)

const newLine = rune('\n')

var (
	p  *bluemonday.Policy
	mu sync.Mutex
)

func init() {
	mu.Lock()
	defer mu.Unlock()

	p = bluemonday.StripTagsPolicy()
}

func StripNewLine(s string, maxNewLine int) string {
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	j := 0
	for _, c := range s {
		if c == newLine {
			j++
			if j <= maxNewLine {
				builder.WriteRune(c)
			}
			continue
		}
		builder.WriteRune(c)
		j = 0
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
