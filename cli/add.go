package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/uaraven/clotp/keyrings"
	"github.com/uaraven/gotp"
)

type AddCmd struct {
	Uri  string `arg:"positional,required"`
	Name string `arg:"--name" help:"Optional name of the code to refer to it later"`
}

func Add(cmd *AddCmd, keys keyrings.Keys) (string, error) {
	u, err := url.Parse(cmd.Uri)
	if err != nil {
		return "", err
	}
	if u.Scheme == migrationScheme {
		return AddMigration(cmd.Uri, keys)
	}
	otp, err := gotp.OTPFromUri(cmd.Uri)
	if err != nil {
		return "", err
	}

	var name string

	if cmd.Name != "" {
		name = cmd.Name
	} else {
		name = otp.GetLabelRepr()
	}

	return fmt.Sprintf("Added %s(%s) as %s", otp.Account, otp.Issuer, name), keys.AddKey(name, cmd.Uri)
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
