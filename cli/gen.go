package cli

import (
	"fmt"
	"github.com/uaraven/gotp"
)

type GenCmd struct {
	Key    string `arg:"--key" help:"Base32-encoded secret key" required:"true"`
	Hash   string `arg:"--hash" help:"Hash algorithm SHA-1, SHA-256 or SHA-512" choice:"SHA-1" choice:"SHA-256" choice:"SHA-512" default:"SHA-1"`
	Digits int    `arg:"--digits" help:"Number of digits in the code, 6-10" default:"6"`
	Period uint   `arg:"--period" help:"TOTP code rotation period, default is 30 seconds" default:"30"`
}

func Gen(cmd *GenCmd) (string, error) {
	secret, err := gotp.DecodeKey(cmd.Key)
	if err != nil {
		return "", fmt.Errorf("invalid secret key: %v", err)
	}
	hash, err := nameToHash(cmd.Hash)
	if err != nil {
		return "", fmt.Errorf("invalid hash algorithm: %v", err)
	}
	if cmd.Digits < 6 || cmd.Digits > 10 {
		return "", fmt.Errorf("invalid number of digits, must be in 6-10 range")
	}
	otp := gotp.NewTOTPHash(secret, cmd.Digits, int(cmd.Period), 0, hash)
	return otp.Now(), nil
}
