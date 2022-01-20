package cli

import (
	"fmt"

	"github.com/uaraven/clotp/keyrings"
)

type SetCounterCmd struct {
	Id      string `arg:"--id" help:"HOTP code identifier"`
	Name    string `arg:"--name" help:"HOTP code name"`
	Counter int64  `arg:"positional,required" help:"New counter value"`
}

func SetCounter(cmd *SetCounterCmd, keys keyrings.Keys) (string, error) {
	if cmd.Name == "" {
		return "", fmt.Errorf("OTP code name must be specified")
	}
	otp, err := keys.GetByName(cmd.Name)
	var output string
	switch otp.Key.Type {
	case keyrings.TOTP:
		return "", fmt.Errorf("cannot set counter for TOTP code")
	case keyrings.HOTP:
		otp.Key.Counter = cmd.Counter
		err = keys.UpdateKey(cmd.Name, otp)
		output = fmt.Sprintf("Set counter to %d", cmd.Counter)
	}
	return output, err
}
