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
	domain := strings.ToLower(u.Hostname())
	if data := strings.Split(domain, "."); len(data) > 2 {
		domain = strings.Join(data[len(data)-2:], ".")
	}
	return domain, nil
}

func SafeDomain(url string) string {
	d, _ := Domain(url)
	return d
}
