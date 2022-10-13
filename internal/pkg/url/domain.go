package url

import (
	gourl "net/url"
	"strings"
)

func Domain(url string) (string, error) {
	u, err := gourl.Parse(url)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ToLower(u.Hostname()), "www.", ""), nil
}

func SafeDomain(url string) string {
	d, _ := Domain(url)
	return d
}
