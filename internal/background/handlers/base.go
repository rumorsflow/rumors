package handlers

import (
	"strconv"
	"strings"
)

func Args(str string) []string {
	return strings.Split(str, " ")
}

func Id(str string) (id int64, rest []string) {
	a := Args(str)
	if len(a) > 1 {
		if v, err := strconv.ParseInt(a[1], 10, 64); err == nil {
			id = v
		}
		rest = a[2:]
	}
	return
}

func Index(str string) uint64 {
	if v, err := strconv.ParseUint(str, 10, 64); err == nil {
		return v
	}
	return 0
}

func Size(str string) uint32 {
	if v, err := strconv.ParseUint(str, 10, 32); err == nil && v > 0 {
		if v > 50 {
			return 50
		}
		return uint32(v)
	}
	return 10
}

func Pagination(str string) (i uint64, s uint32, filters []string) {
	a := Args(str)
	if strings.ToLower(a[0]) == "list" {
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
