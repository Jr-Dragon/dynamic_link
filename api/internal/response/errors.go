package response

import (
	"errors"
	"fmt"
)

var ErrInvalidCode = errors.New("invalid code")

func InvalidCode(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidCode, err)
}
