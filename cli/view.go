package cli

import (
	"fmt"
	"strings"

	"github.com/uaraven/clotp/keyrings"
	"github.com/uaraven/gotp"
)

type ViewCmd struct {
	Name string `arg:"positional,required" help:"Name of the OTP code to view"`
}

func View(cmd *ViewCmd, keys keyrings.Keys) (string, error) {
	if cmd.Name == "" {
		return "", fmt.Errorf("name must be provided")
	}
	code, err := keys.GetByName(cmd.Name)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("        Name: %s\n", code.Key.Name))
	sb.WriteString(fmt.Sprintf("Account Name: %s\n", code.Key.Account))
	sb.WriteString(fmt.Sprintf("      Issuer: %s\n", code.Key.Issuer))
	sb.WriteString(fmt.Sprintf("    OTP Type: %s\n", code.Key.GetTypeString()))
	hashName, err := gotp.HashAlgorithmName(code.OTP.Hash)
	if err != nil {
		return "", err
	}
	sb.WriteString(fmt.Sprintf("   Hash Type: %s\n", hashName))
	sb.WriteString(fmt.Sprintf(" Code Digits: %d\n", code.OTP.Digits))
	switch code.Key.Type {
	case keyrings.TOTP:
		sb.WriteString(fmt.Sprintf(" Time offset: %d\n", code.OTP.StartTime))
		sb.WriteString(fmt.Sprintf("   Time step: %d\n", code.OTP.TimeStep))
	case keyrings.HOTP:
		sb.WriteString(fmt.Sprintf("     Counter: %d\n", code.Key.Counter))
	default:
		return "", fmt.Errorf("unsupported OTP type: %d", code.Key.Type)
	}
	sb.WriteString(fmt.Sprintf("      Secret: %s\n", gotp.EncodeKey(code.OTP.Secret)))
	otpCode, err := code.AsOTP()
	if err != nil {
		return "", err
	}
	sb.WriteString(fmt.Sprintf("    Auth URI: %s\n", otpCode.ProvisioningUri(code.Key.Account, code.Key.Issuer)))
	return sb.String(), nil
}
