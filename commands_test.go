package main

import (
	"testing"

	"github.com/uaraven/gotp"
)

type mockKeys struct {
	Keys
	keys []OTPKey
}

func (mk *mockKeys) ListOTPs() ([]OTPKey, error) {
	return make([]OTPKey, 0), nil
}

func (mk *mockKeys) AddKey(id string, name string, otpUri string) error {
	otp, err := gotp.OTPFromUri(otpUri)
	if err != nil {
		return err
	}
	otpkey := OTPKey{
		Id:         id,
		OTPKeyData: *otp,
	}
	mk.keys = append(mk.keys, otpkey)
	return nil
}

func NewMockKeys() *mockKeys {
	return &mockKeys{}
}

func TestAdd(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth://totp/vendor:label?secret=NNSXS&issuer=vendor",
		Name: "",
	}
	keys := NewMockKeys()
	err := Add(&cmd, keys)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(keys.keys) != 1 {
		t.Errorf("expected 1 key in the keyring")
	}
	if keys.keys[0].Label != "vendor:label" {
		t.Errorf("expected label to be 'label', but got %s instead", keys.keys[0].Label)
	}
}
