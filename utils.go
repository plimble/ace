package ace

import (
	"reflect"
	"unsafe"
)

func bytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func stringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func concat(s ...string) string {
	size := 0
	for i := 0; i < len(s); i++ {
		size += len(s[i])
	}

	buf := make([]byte, 0, size)

	for i := 0; i < len(s); i++ {
		buf = append(buf, stringToBytes(s[i])...)
	}

	return bytesToString(buf)
}
