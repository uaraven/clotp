package cli

import (
	"fmt"

	"github.com/uaraven/clotp/keyrings"
)

type SetCounterCmd struct {
	Name    string `arg:"positional,required" help:"HOTP code name"`
	Counter int64  `arg:"positional,required" help:"New counter value"`
}

func SetCounter(cmd *SetCounterCmd, keys keyrings.Keys) (string, error) {
	if cmd.Name == "" {
		return "", fmt.Errorf("OTP code name must be specified")
	}
	key, err := keys.GetKeyByName(cmd.Name)
	if err != nil {
		return "", err
	}
	if key.Type == keyrings.TOTP {
		return "", fmt.Errorf("cannot set counter for TOTP code")
	}
	otp, err := keys.GetByName(cmd.Name)
	if err != nil {
		return "", err
	}
	var output string
	otp.Key.Counter = cmd.Counter
	err = keys.UpdateKey(cmd.Name, otp)
	output = fmt.Sprintf("Set counter to %d", cmd.Counter)
	return output, err
}
