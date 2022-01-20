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
	keys.keys[*keyrings.NewKey(keyrings.HOTP, "vendor1", "label1")] = *otp1
	keys.keys[*keyrings.NewKey(keyrings.TOTP, "vendor2", "label2")] = *otp2

	cmd := ListCmd{}
	out, err := List(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "HOTP") || !strings.Contains(out, "vendor1") {
		t.Errorf("List should contain HOTP and vendor1, actual:\n%s", out)
	}
	if !strings.Contains(out, "TOTP") || !strings.Contains(out, "vendor2") {
		t.Errorf("List should contain TOTP and vendor2, actual:\n%s", out)
	}
}
