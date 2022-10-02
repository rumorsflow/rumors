package handlers

import (
	"encoding/binary"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/spf13/cast"
	"strings"
)

func Args(str string) []string {
	return strings.Split(str, " ")
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
	i = Index(slice.Safe(a, 0))
	s = Size(slice.Safe(a, 1))

	if len(a) > 2 {
		filters = a[2:]
	}
	return
}

func Int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func BytesToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
