package cli

import (
	"testing"

	"github.com/uaraven/clotp/keyrings"
)

func TestRemove(t *testing.T) {
	otpUrl1 := "otpauth://hotp/vendor1:label1?counter=1&issuer=vendor1&secret=NNSXS"
	otpUrl2 := "otpauth://totp/vendor2:label2?issuer=vendor2&secret=NNSXS"
	otp1, _ := keyrings.ParamsFromUri(otpUrl1)
	otp2, _ := keyrings.ParamsFromUri(otpUrl2)
	keys := NewMockKeys()
	keys.keys["label1"] = keyrings.KeyringItem{
		Key: keyrings.NewKey(keyrings.HOTP, "vendor1", "label1"),
		OTP: otp1,
	}
	keys.keys["label2"] = keyrings.KeyringItem{
		Key: keyrings.NewKey(keyrings.HOTP, "vendor2", "label2"),
		OTP: otp2,
	}

	cmd := RemoveCmd{
		Name: "label1",
	}

	_, err := Remove(&cmd, keys)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(keys.keys) != 1 {
		t.Errorf("expected 1 key in the keyring")
	}
}

func TestRemoveInvalidName(t *testing.T) {
	otpUrl1 := "otpauth://hotp/vendor1:label1?counter=1&issuer=vendor1&secret=NNSXS"
	otpUrl2 := "otpauth://totp/vendor2:label2?issuer=vendor2&secret=NNSXS"
	otp1, _ := keyrings.ParamsFromUri(otpUrl1)
	otp2, _ := keyrings.ParamsFromUri(otpUrl2)
	keys := NewMockKeys()
	keys.keys["label1"] = keyrings.KeyringItem{
		Key: keyrings.NewKey(keyrings.HOTP, "vendor1", "label1"),
		OTP: otp1,
	}
	keys.keys["label2"] = keyrings.KeyringItem{
		Key: keyrings.NewKey(keyrings.HOTP, "vendor2", "label2"),
		OTP: otp2,
	}

	cmd := RemoveCmd{
		Name: "label3",
	}

	_, err := Remove(&cmd, keys)
	if err == nil {
		t.Error("expected remove to fail")
	}
	if len(keys.keys) != 2 {
		t.Errorf("expected 2 keys in the keyring")
	}
}
