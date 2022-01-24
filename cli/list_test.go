package cli

import (
	"strings"
	"testing"

	"github.com/uaraven/clotp/keyrings"
)

func TestList(t *testing.T) {
	otpUrl1 := "otpauth://hotp/vendor1:label1?counter=1&issuer=vendor1&secret=NNSXS"
	otpUrl2 := "otpauth://totp/vendor2:label2?issuer=vendor2&secret=NNSXS"
	otp1, _ := keyrings.ParamsFromUri(otpUrl1)
	otp2, _ := keyrings.ParamsFromUri(otpUrl2)
	keys := NewMockKeys()
	keys.keys["label1"] = keyrings.KeyringItem{
		Key: keyrings.NewKey(keyrings.HOTP, "vendor1", "label1"),
		OTP: otp1,
	}
	keys.keys["customLabel"] = keyrings.KeyringItem{
		Key: keyrings.NewKey(keyrings.TOTP, "vendor2", "customLabel"),
		OTP: otp2,
	}
	cmd := ListCmd{}
	out, err := List(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "HOTP") || !strings.Contains(out, "label1 vendor1") {
		t.Errorf("List should contain HOTP and label1 vendor1, actual:\n%s", out)
	}
	if !strings.Contains(out, "TOTP") || !strings.Contains(out, "customLabel vendor2") {
		t.Errorf("List should contain TOTP and customLabel vendor2, actual:\n%s", out)
	}
}
