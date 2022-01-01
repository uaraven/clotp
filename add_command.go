package main

import "github.com/uaraven/gotp"

func Add(cmd AddCmd, keys *Keys) error {
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
