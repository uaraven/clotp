package cli

import (
	"strings"
	"testing"

	"github.com/uaraven/clotp/keyrings"
)

func TestView(t *testing.T) {
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

	cmd := &ViewCmd{
		Name: "label1",
	}

	output, err := View(cmd, keys)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(output, "Secret: NNSXS") {
		t.Errorf("Expected output to contain 'Secret: NNSXS', but it was: %s", output)
	}
}
