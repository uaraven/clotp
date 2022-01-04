package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/uaraven/clotp/migration"
	"github.com/uaraven/gotp"
)

func Add(cmd *AddCmd, keys Keys) error {
	id := uuid.New().String()
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

	return keys.AddKey(id, name, cmd.Uri)
}

func List(cmd *ListCmd, keys Keys) error {
	otps, err := keys.ListOTPs()
	if err != nil {
		return err
	}
	if len(otps) == 0 {
		fmt.Println("No stored codes")

	} else {
		for _, otp := range otps {
			fmt.Printf("[%s] %s - %s\n", otp.Id, otp.Label, otp.OTP.ProvisioningUri(otp.Label, otp.Issuer))
		}
	}
	return nil
}

func Remove(cmd *RemoveCmd, keys Keys) error {
	if cmd.Id != "" {
		fmt.Printf("Removing code with id '%s'\n", cmd.Id)
		return keys.RemoveById(cmd.Id)
	} else if cmd.Name != "" {
		fmt.Printf("Removing code with name '%s'\n", cmd.Id)
		return keys.RemoveByName(cmd.Name)
	} else {
		return fmt.Errorf("neither id nor name specified")
	}
}

func Decode(cmd *DecodeCmd) error {
	payload, err := migration.UnmarshalURL(cmd.Uri)
	if err != nil {
		return err
	}
	for _, otpParam := range payload.OtpParameters {
		otp, err := otpFromParameters(otpParam)
		if err != nil {
			return err
		}
		fmt.Println(otp.ProvisioningUri(otpParam.Name, otpParam.Issuer))
	}
	return nil
}
