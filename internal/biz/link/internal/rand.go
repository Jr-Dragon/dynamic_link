package internal

import (
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randomBytes(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func checksum(s []byte, k uint32) string {
	tbl := crc32.MakeTable(k)
	return strconv.FormatUint(uint64(crc32.Checksum(s, tbl)), 36)
}

func GenerateRandomString(n int, k uint32) string {
	b := randomBytes(n)

	var sb strings.Builder
	sb.Write(b)
	sb.WriteRune('-')
	sb.WriteString(checksum(b, k))

	return sb.String()
}
