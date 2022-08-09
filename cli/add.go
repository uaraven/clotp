package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/uaraven/clotp/keyrings"
	"github.com/uaraven/gotp"
)

type AddCmd struct {
	Uri    string `arg:"positional,required"`
	Name   string `arg:"--name" help:"Optional name of the code to refer to it later"`
	IsCode bool   `arg:"--code" help:"Pass just secret code instead of full URI"`
}

func parseURI(uri string) (*gotp.OTPKeyData, error) {
	_, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	otp, err := gotp.OTPFromUri(uri)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func Add(cmd *AddCmd, keys keyrings.Keys) (string, error) {
	var otp *gotp.OTPKeyData
	var err error
	var uri string
	if cmd.IsCode {
		if cmd.Name == "" {
			return "", fmt.Errorf("name must be specified when create the key from the code")
		}
		key, err := gotp.DecodeKey(cmd.Uri)
		if err == nil {
			totp := gotp.NewDefaultTOTP(key)
			otp = &gotp.OTPKeyData{
				OTP:     totp,
				Account: cmd.Name,
				Issuer:  cmd.Name,
			}
			uri = totp.ProvisioningUri(cmd.Name, cmd.Name)
		}
	} else {
		otp, err = parseURI(cmd.Uri)
		uri = cmd.Uri
	}
	if err != nil {
		return "", err
	}

	var name string

	if cmd.Name != "" {
		name = cmd.Name
	} else {
		name = otp.GetLabelRepr()
	}
	var output string
	if otp.Issuer == "" {
		output = fmt.Sprintf("Added %s with name %s", otp.Account, name)
	} else {
		output = fmt.Sprintf("Added %s(%s) with name %s", otp.Account, otp.Issuer, name)
	}

	return output, keys.AddKey(name, uri)
}

func AddMigration(migrationUri string, keys keyrings.Keys) (string, error) {
	otps, err := otpFromMigrationUri(migrationUri)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for _, otp := range otps {
		err = keys.AddKey(otp.GetLabelRepr(), otp.OTP.ProvisioningUri(otp.Account, otp.Issuer))
		if err != nil {
			return "", err
		}
		result.WriteString(fmt.Sprintf("Added OTP %s\n", otp.GetLabelRepr()))
	}
	return result.String(), nil
}
