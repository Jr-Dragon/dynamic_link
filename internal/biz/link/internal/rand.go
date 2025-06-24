package internal

import (
	"hash/crc32"
	"strconv"

	"github.com/gookit/goutil/strutil"
)

func checksum(s []byte, k uint32) string {
	tbl := crc32.MakeTable(k)
	return strconv.FormatUint(uint64(crc32.Checksum(s, tbl)), 36)
}

func GenerateRandomString(n int, k uint32) (string, error) {
	b, err := strutil.RandomBytes(n)
	if err != nil {
		return "", err
	}

	return string(b) + "-" + checksum(b, k), nil
}
