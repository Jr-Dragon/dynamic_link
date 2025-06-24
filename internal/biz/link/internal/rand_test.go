package internal

import (
	"hash/crc32"
	"strconv"
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	testcases := []struct {
		name   string
		length int
	}{
		{"Normal Length", 6},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := GenerateRandomString(tc.length, crc32.IEEE)
			if err != nil {
				t.Errorf("GenerateRandomString failed: %v", err)
			}

			splitted := strings.Split(gen, "-")
			if len(splitted) != 2 {
				t.Errorf("splitted length should be 2, got %d", len(splitted))
			}

			k, c := splitted[0], splitted[1]
			if len(k) != tc.length {
				t.Errorf("splitted length should be %d, got %d", tc.length, len(splitted[0]))
			}

			cint, err := strconv.ParseUint(c, 36, 32)
			if err != nil {
				t.Errorf("failed to parse checksum: %v", err)
			}

			if uint32(cint) != crc32.Checksum([]byte(k), crc32.IEEETable) {
				t.Errorf("checksum mismatch: expected %d, got %d", crc32.Checksum([]byte(k), crc32.IEEETable), cint)
			}
		})
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateRandomString(6, crc32.IEEE)
		if err != nil {
			b.Errorf("GenerateRandomString failed: %v", err)
		}
	}
}
