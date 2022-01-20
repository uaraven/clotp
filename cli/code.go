package cli

import (
	"fmt"

	"github.com/uaraven/clotp/keyrings"
	"github.com/uaraven/gotp"
)

type CodeCmd struct {
	Name    string `arg:"positional"`
	Counter int64  `arg:"--counter" help:"Override counter for the HOTP code"`
	Copy    bool   `arg:"--copy" help:"Copy generated code to clipboard"`
}

func Code(cmd *CodeCmd, keys keyrings.Keys) (string, error) {
	if cmd.Name == "" {
		return "", fmt.Errorf("OTP code name must be specified")
	}
	otpKey, err := keys.GetByName(cmd.Name)
	if err != nil {
		return "", err
	}
	var output string
	otp, err := otpKey.AsOTP()
	if err != nil {
		return "", err
	}
	switch otpKey.Key.Type {
	case keyrings.TOTP:
		output = otp.(*gotp.TOTP).Now()
	case keyrings.HOTP:
		// generating HOTP means we need to update the counter
		hotp := otp.(*gotp.HOTP)
		var counter int64
		if cmd.Counter > 0 {
			counter = cmd.Counter
		} else {
			counter = otpKey.Key.Counter
		}
		if counter == 0 {
			counter = 1
		}
		output = hotp.GenerateOTP(counter)
		newCounter := hotp.GetCounter()
		otpKey.Key.Counter = newCounter
		keys.UpdateKey(cmd.Name, otpKey)
	}
	return output, err
}
