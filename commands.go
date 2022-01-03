package main

import (
	"fmt"

	"github.com/google/uuid"
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
	for _, otp := range otps {
		fmt.Printf("[%s: %s] %s\n", otp.Id, otp.Label, otp.OTP.ProvisioningUri(otp.Label, otp.Issuer))
	}
	return nil
}

func Remove(cmd *RemoveCmd, keys Keys) error {
	if cmd.Id != "" {
		keys.RemoveById(cmd.Id)
		return nil
	} else if cmd.Name != "" {
		keys.RemoveByName(cmd.Name)
		return nil
	} else {
		return fmt.Errorf("neither id nor name specified")
	}
}
