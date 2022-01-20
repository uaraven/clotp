package main

import (
	"fmt"
	"strings"
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

func (mk *mockKeys) GetById(id string) (*OTPKey, error) {
	for _, item := range mk.keys {
		if item.Id == id {
			return &OTPKey{item.OTPKeyData, id}, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func NewMockKeys() *mockKeys {
	return &mockKeys{}
}

func TestCodeTOTP(t *testing.T) {
	key := []byte("key")
	totp := gotp.NewDefaultTOTP(key)
	keys := NewMockKeys()
	keys.AddKey("1", "1", totp.ProvisioningUri("account", "issuer"))

	cmd := &CodeCmd{
		Id: "1",
	}

	expected := totp.Now()
	actual, err := Code(cmd, keys)

	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("Expected code to be '%s', but got '%s'", expected, actual)
	}
}

func TestCodeHOTP(t *testing.T) {
	key := []byte("key")
	hotp := gotp.NewDefaultHOTP(key, 1)
	keys := NewMockKeys()
	keys.AddKey("1", "1", hotp.ProvisioningUri("account", "issuer"))

	cmd := &CodeCmd{
		Id:      "1",
		Counter: -1,
	}

	expected := "023900"
	actual, err := Code(cmd, keys)

	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("Expected code to be '%s', but got '%s'", expected, actual)
	}

	expected2 := "296062"
	actual2, err := Code(cmd, keys)

	if err != nil {
		t.Error(err)
	}

	if actual2 != expected2 {
		t.Errorf("Expected code to be '%s', but got '%s' for the second invocation", expected2, actual2)
	}
}

func TestSetCounter(t *testing.T) {
	key := []byte("key")
	totp := gotp.NewDefaultHOTP(key, 1)
	keys := NewMockKeys()
	keys.AddKey("1", "1", totp.ProvisioningUri("account", "issuer"))

	cmd := &SetCounterCmd{
		Id:      "1",
		Name:    "",
		Counter: 10,
	}

	_, err := SetCounter(cmd, keys)

	if err != nil {
		t.Error(err)
	}

	otp, err := keys.GetById("1")
	if err != nil {
		t.Error(err)
	}

	hotpc := otp.OTP.(*gotp.HOTP)

	if hotpc.GetCounter() != 10 {
		t.Errorf("failed to set counter, expected 10, got: %d", hotpc.GetCounter())
	}
}

func TestSetCounterForTOTP(t *testing.T) {
	key := []byte("key")
	hotp := gotp.NewDefaultTOTP(key)
	keys := NewMockKeys()
	keys.AddKey("1", "1", hotp.ProvisioningUri("account", "issuer"))

	cmd := &SetCounterCmd{
		Id:      "1",
		Name:    "",
		Counter: 10,
	}

	_, err := SetCounter(cmd, keys)

	if err == nil {
		t.Errorf("SetCounter should fail for TOTP code")
	}
}
