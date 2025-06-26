package link

import (
	"errors"
	"testing"
)

func TestLink_ValidateSimple(t *testing.T) {
	testcases := []struct {
		name    string
		code    []byte
		wantErr error
	}{
		{
			name: "success",
			code: []byte(code),
		},
		{
			name:    "failed: code length too short",
			code:    []byte("foo"),
			wantErr: errors.New("code length"),
		},
		{
			name:    "failed: code length too long",
			code:    []byte("this_is_a_too_long_code"),
			wantErr: errors.New("code length"),
		},
		{
			name:    "failed: code format: without '-'",
			code:    []byte("invalid-format"),
			wantErr: errors.New("code format"),
		},
		{
			name:    "failed: code checksum mismatch",
			code:    []byte("hSxIIn-abcdef"),
			wantErr: errors.New("code checksum"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			link := NewLink(cfg, nil)
			if err := link.ValidateSimple(tc.code); err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("ValidateSimple() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
