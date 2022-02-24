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
		return fmt.Sprintf("AccountName='%s', Issuer='%s', Secret='%s', Hash=%s, Digits=%d, TimeStep=%d", account, issuer,
			gotp.EncodeKey(otpg.Secret), hash, otpg.Digits, otpg.TimeStep), nil
	case *gotp.HOTP:
		hash, _ := gotp.HashAlgorithmName(otpg.Hash)
		return fmt.Sprintf("AccountName='%s', Issuer='%s', Secret='%s', Hash=%s, Digits=%d, Counter=%d", account, issuer,
			gotp.EncodeKey(otpg.Secret), hash, otpg.Digits, otpg.GetCounter()), nil
	default:
		return "", fmt.Errorf("unknown OTP type")
	}
}

func decodeMigrationUri(migrationUri string, parse bool) ([]string, error) {
	payload, err := migration.UnmarshalURL(migrationUri)
	if err != nil {
		return nil, err
	}
	var result = make([]string, 0)
	for _, otpParam := range payload.OtpParameters {
		otp, err := otpFromParameters(otpParam)
		if err != nil {
			return nil, err
		}
		if parse {
			line, err := parseOTP(otpParam.Name, otpParam.Issuer, otp)
			if err != nil {
				return nil, err
			}
			result = append(result, fmt.Sprintf("%s\n", line))
		} else {
			result = append(result, fmt.Sprintf("%s\n", otp.ProvisioningUri(otpParam.Name, otpParam.Issuer)))
		}
	}
	return result, nil
}

func Decode(cmd *DecodeCmd) (string, error) {
	results, err := decodeMigrationUri(cmd.Uri, cmd.Parse)
	if err != nil {
		return "", err
	}
	return strings.Join(results, "\n"), nil
}
