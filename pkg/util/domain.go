package util

import (
	"net/url"
	"strings"
)

func Domain(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	domain := strings.ToLower(u.Hostname())
	if data := strings.Split(domain, "."); len(data) > 2 {
		domain = strings.Join(data[len(data)-2:], ".")
	}
	return domain, nil
}

func SafeDomain(link string) string {
	d, _ := Domain(link)
	return d
}
