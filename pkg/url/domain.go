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

func MustDomain(url string) string {
	d, err := Domain(url)
	if err != nil {
		panic(err)
	}
	return d
}
