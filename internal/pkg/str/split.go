package str

import "strings"

func SplitMax(text, sep string, max int) []string {
	var data []string
	for {
		if len(text) > max {
			i := strings.LastIndex(text[:max], sep)
			if i <= 0 {
				i = max
			}
			data = append(data, text[:i])
			text = text[i+1:]
		} else if len(text) > 0 {
			data = append(data, text)
			break
		} else {
			break
		}
	}
	return data
}
