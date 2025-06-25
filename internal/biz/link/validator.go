package link

import (
	"errors"
	"hash/crc32"
	"strconv"
)

func (link *Link) ValidateSimple(code []byte) error {
	if len(code) > 20 || len(code) < 6 {
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
