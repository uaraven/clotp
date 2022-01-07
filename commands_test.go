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

func TestAddTOTP(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth://totp/vendor:label?secret=NNSXS&issuer=vendor",
		Name: "",
	}
	keys := NewMockKeys()
	_, err := Add(&cmd, keys)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(keys.keys) != 1 {
		t.Errorf("expected 1 key in the keyring")
	}
	if keys.keys[0].Label != "label" {
		t.Errorf("expected label to be 'label', but got %s instead", keys.keys[0].Label)
	}
}

func TestAddHOTP(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth://hotp/vendor:label?secret=NNSXS&issuer=vendor&counter=1",
		Name: "",
	}
	keys := NewMockKeys()
	_, err := Add(&cmd, keys)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(keys.keys) != 1 {
		t.Errorf("expected 1 key in the keyring")
	}
	if keys.keys[0].Label != "label" {
		t.Errorf("expected label to be 'label', but got %s instead", keys.keys[0].Label)
	}
}

func TestAddMigration(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth-migration://offline?data=ChwKBEtleTESBVRlc3QxGgdJc3N1ZXIxIAEoATACChwKBEtleTISBVRlc3QyGgdJc3N1ZXIyIAEoATAC",
		Name: "",
	}
	keys := NewMockKeys()
	_, err := Add(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
	if len(keys.keys) != 2 {
		t.Errorf("expected 2 keys in the keyring")
	}
	if keys.keys[0].Label != "Test1" {
		t.Errorf("expected label to be 'Test1', but got %s instead", keys.keys[0].Label)
	}
	if keys.keys[1].Label != "Test2" {
		t.Errorf("expected label to be 'Test2', but got %s instead", keys.keys[0].Label)
	}
}

func TestList(t *testing.T) {
	otpUrl1 := "otpauth://hotp/vendor1:label1?counter=1&issuer=vendor1&secret=NNSXS"
	otpUrl2 := "otpauth://totp/vendor2:label2?issuer=vendor2&secret=NNSXS"
	otp1, _ := gotp.OTPFromUri(otpUrl1)
	otp2, _ := gotp.OTPFromUri(otpUrl2)
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
	out, err := List(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, otpUrl1) {
		t.Errorf("List should contain %s, actual:\n%s", otpUrl1, out)
	}
	if !strings.Contains(out, otpUrl2) {
		t.Errorf("List should contain %s, actual:\n%s", otpUrl2, out)
	}
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
