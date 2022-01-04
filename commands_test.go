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
	return mk.keys, nil
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

func TestAddTOTP(t *testing.T) {
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

func TestAddHOTP(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth://hotp/vendor:label?secret=NNSXS&issuer=vendor&counter=1",
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

func TestAddMigration(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth-migration://offline?data=ChwKBEtleTESBVRlc3QxGgdJc3N1ZXIxIAEoATACChwKBEtleTISBVRlc3QyGgdJc3N1ZXIyIAEoATAC",
		Name: "",
	}
	keys := NewMockKeys()
	err := Add(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
	if len(keys.keys) != 2 {
		t.Errorf("expected 2 key in the keyring")
	}
	if keys.keys[0].Label != "Issuer1:Test1" {
		t.Errorf("expected label to be 'Issuer1:Test1', but got %s instead", keys.keys[0].Label)
	}
	if keys.keys[1].Label != "Issuer2:Test2" {
		t.Errorf("expected label to be 'Issuer2:Test2', but got %s instead", keys.keys[0].Label)
	}
}

func TestList(t *testing.T) {
	otp1, _ := gotp.OTPFromUri("otpauth://hotp/label1?secret=NNSXS&issuer=vendor1&counter=1")
	otp2, _ := gotp.OTPFromUri("otpauth://totp/label2?secret=NNSXS&issuer=vendor2")
	keys := NewMockKeys()
	keys.keys = []OTPKey{
		{
			Id:         "1",
			OTPKeyData: *otp1,
		},
		{
			Id:         "2",
			OTPKeyData: *otp2,
		},
	}

	cmd := ListCmd{}
	err := List(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
}
