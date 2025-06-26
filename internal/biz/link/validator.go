package link

import (
	"errors"
	"hash/crc32"
	"strconv"
)

func (link *Link) ValidateSimple(code []byte) error {
	// The code is formatted in "`$key`-`$checksum`"
	// - len(key): 6
	// - minimum checksum: "0"
	// - maximum checksum: "1z141z3" (4294967295)
	// 8 <= len(code) <= 14
	if len(code) > 14 || len(code) < 8 {
		return errors.New("code length")
	}

	if code[6] != byte('-') {
		return errors.New("code format")
	}

	k, crc := code[:6], code[7:]
	cint, err := strconv.ParseUint(string(crc), 36, 32)
	if err != nil {
		return err
	}
	if crc32.Checksum(k, crc32.MakeTable(link.cfg.App.Key)) != uint32(cint) {
		return errors.New("code checksum")
	}

	return nil
}
