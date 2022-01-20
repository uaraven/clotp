package cli

import (
	"testing"

	"github.com/uaraven/gotp"
)

func TestSetCounter(t *testing.T) {
	key := []byte("key")
	totp := gotp.NewDefaultHOTP(key, 1)
	keys := NewMockKeys()
	keys.AddKey("1", totp.ProvisioningUri("account", "issuer"))

	cmd := &SetCounterCmd{
		Name:    "1",
		Counter: 10,
	}

	_, err := SetCounter(cmd, keys)

	if err != nil {
		t.Error(err)
	}

	otp, err := keys.GetByName("1")
	if err != nil {
		t.Error(err)
	}

	if otp.Key.Counter != 10 {
		t.Errorf("failed to set counter, expected 10, got: %d", otp.Key.Counter)
	}
}

func TestSetCounterForTOTP(t *testing.T) {
	key := []byte("key")
	hotp := gotp.NewDefaultTOTP(key)
	keys := NewMockKeys()
	keys.AddKey("1", hotp.ProvisioningUri("account", "issuer"))

	cmd := &SetCounterCmd{
		Name:    "1",
		Counter: 10,
	}

	_, err := SetCounter(cmd, keys)

	if err == nil {
		t.Errorf("SetCounter should fail for TOTP code")
	}
}
