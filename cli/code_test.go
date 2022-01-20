package cli

import (
	"testing"

	"github.com/uaraven/gotp"
)

func TestCodeTOTP(t *testing.T) {
	key := []byte("key")
	totp := gotp.NewDefaultTOTP(key)
	keys := NewMockKeys()
	keys.AddKey("1", totp.ProvisioningUri("account", "issuer"))

	cmd := &CodeCmd{
		Name: "1",
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
	keys.AddKey("1", hotp.ProvisioningUri("account", "issuer"))

	cmd := &CodeCmd{
		Name:    "1",
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
