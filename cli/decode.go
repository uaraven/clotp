package cli

import (
	"fmt"
	"strings"

	"github.com/uaraven/clotp/migration"
	"github.com/uaraven/gotp"
)

type DecodeCmd struct {
	Uri   string `arg:"positional,required" help:"Google Authenticator Export URI"`
	Parse bool   `arg:"-p,--parse" help:"Parse URIs and print each part separately"`
}

func parseOTP(account string, issuer string, otp gotp.OTP) (string, error) {
	switch otpg := otp.(type) {
	case *gotp.TOTP:
		hash, _ := gotp.HashAlgorithmName(otpg.Hash)
		return fmt.Sprintf("AccountName=%s, Issuer=%s, Secret=%s, Hash=%s, Digits=%d, TimeStep=%d", account, issuer,
			gotp.EncodeKey(otpg.Secret), hash, otpg.Digits, otpg.TimeStep), nil
	case *gotp.HOTP:
		hash, _ := gotp.HashAlgorithmName(otpg.Hash)
		return fmt.Sprintf("AccountName=%s, Issuer=%s, Secret=%s, Hash=%s, Digits=%d, Counter=%d", account, issuer,
			gotp.EncodeKey(otpg.Secret), hash, otpg.Digits, otpg.GetCounter()), nil
	default:
		return "", fmt.Errorf("unknown OTP type")
	}
}

func Decode(cmd *DecodeCmd) (string, error) {
	payload, err := migration.UnmarshalURL(cmd.Uri)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for _, otpParam := range payload.OtpParameters {
		otp, err := otpFromParameters(otpParam)
		if err != nil {
			return "", err
		}
		if cmd.Parse {
			line, err := parseOTP(otpParam.Name, otpParam.Issuer, otp)
			if err != nil {
				return "", err
			}
			result.WriteString(line)
		}
		result.WriteString(fmt.Sprintf("%s\n", otp.ProvisioningUri(otpParam.Name, otpParam.Issuer)))
	}
	return result.String(), nil
}
