package main

import (
	"fmt"

	"github.com/uaraven/gotp"
)

func Add(cmd *AddCmd, keys *Keys) error {
	otp, err := gotp.OTPFromUri(cmd.Uri)
	if err != nil {
		return err
	}

	var name string

	if cmd.Name != "" {
		name = cmd.Name
	} else {
		name = otp.Label
	}

	return keys.AddKey(name, cmd.Uri)
}

func List(cmd *ListCmd, keys *Keys) error {
	otps, err := keys.ListOTPs()
	if err != nil {
		return err
	}
	for _, otp := range otps {
		fmt.Printf("%s - %s\n", otp.Label, otp.OTP.ProvisioningUri(otp.Label, otp.Issuer))
	}
	return nil
}
