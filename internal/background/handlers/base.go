package handlers

import (
	"github.com/spf13/cast"
	"strings"
)

func Args(str string) []string {
	return strings.Split(str, " ")
}

func Id(str string) (id int64, rest []string) {
	a := Args(str)
	if len(a) > 1 {
		id = cast.ToInt64(a[1])
		rest = a[2:]
	}
	return
}

func Index(str string) uint64 {
	return cast.ToUint64(str)
}

func Size(str string) uint32 {
	if v := cast.ToUint32(str); v > 0 {
		if v > 50 {
			return 50
		}
		return v
	}
	return 10
}

func Pagination(str string) (i uint64, s uint32, filters []string) {
	a := Args(str)
	if strings.ToLower(a[0]) == "list" || a[0] == "l" || a[0] == "L" {
		a = append([]string{}, a[1:]...)
	}
	t := len(a)

	if t > 0 {
		i = Index(a[0])
	}
	if t > 1 {
		s = Size(a[1])
	} else {
		s = 10
	}
	if t > 2 {
		filters = a[2:]
	}
	return
}
