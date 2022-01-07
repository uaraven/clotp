package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/uaraven/clotp/migration"
	"github.com/uaraven/gotp"
)

func Add(cmd *AddCmd, keys Keys) (string, error) {
	id := uuid.New().String()
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
		name = otp.Label
	}

	return fmt.Sprintf("Added with id=%s", id), keys.AddKey(id, name, cmd.Uri)
}

func AddMigration(migrationUri string, keys Keys) (string, error) {
	otps, err := otpFromMigrationUri(migrationUri)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for _, otp := range otps {
		id := uuid.New().String()
		err = keys.AddKey(id, otp.Label, otp.OTP.ProvisioningUri(otp.Label, otp.Issuer))
		if err != nil {
			return "", err
		}
		result.WriteString(fmt.Sprintf("Added OTP %s with id %s\n", otp.Label, id))
	}
	return result.String(), nil
}

func List(cmd *ListCmd, keys Keys) (string, error) {
	otps, err := keys.ListOTPs()
	if err != nil {
		return "", err
	}
	if len(otps) == 0 {
		return "No stored OTPs", nil
	} else {
		var result strings.Builder
		for _, otp := range otps {
			result.WriteString(fmt.Sprintf("[%s] %s - %s\n", otp.Id, otp.Label, otp.OTP.ProvisioningUri(otp.Label, otp.Issuer)))
		}
		return result.String(), nil
	}
}

func Remove(cmd *RemoveCmd, keys Keys) (string, error) {
	var err error
	var output string
	if cmd.Id != "" && cmd.Name != "" {
		return "", fmt.Errorf("one of --id or --name must be specified")
	}
	if cmd.Id != "" {
		output = fmt.Sprintf("Removed OTP with id '%s'\n", cmd.Id)
		err = keys.RemoveById(cmd.Id)
	} else if cmd.Name != "" {
		output = fmt.Sprintf("Removing code with name '%s'\n", cmd.Id)
		keys.RemoveByName(cmd.Name)
	} else {
		return "", fmt.Errorf("neither id nor name specified")
	}
	return output, err
}

func Code(cmd *CodeCmd, keys Keys) (string, error) {
	if cmd.Id != "" && cmd.Name != "" {
		return "", fmt.Errorf("one of --id or --name must be specified")
	}
	var otp *OTPKey
	var err error
	if cmd.Id != "" {
		otp, err = keys.GetById(cmd.Id)
	} else if cmd.Name != "" {
		otp, err = keys.GetByName(cmd.Name)
	} else {
		return "", fmt.Errorf("neither id nor name specified")
	}
	var output string
	switch otpGen := otp.OTP.(type) {
	case *gotp.TOTP:
		output = otpGen.Now()
	case *gotp.HOTP:
		// generating HOTP means we need to update the counter
		if cmd.Counter >= 0 {
			otpGen.SetCounter(cmd.Counter)
		}
		output = otpGen.CurrentOTP()
		uri := otpGen.ProvisioningUri(otp.Label, otp.Issuer)
		err = keys.AddKey(otp.Id, otp.Label, uri)
	}
	return output, err
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
		result.WriteString(fmt.Sprintf("%s\n", otp.ProvisioningUri(otpParam.Name, otpParam.Issuer)))
	}
	return result.String(), nil
}
