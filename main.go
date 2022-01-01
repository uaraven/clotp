package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

type ListCmd struct {
}

type AddCmd struct {
	Uri  string `arg:"positional,required"`
	Name string `arg:"--name"`
}

var options struct {
	List *ListCmd `arg:"subcommand:list"`
	Add  *AddCmd  `arg:"subcommand:add"`
}

func main() {
	arg.MustParse(&options)

	if options.List == nil && options.Add == nil {
		fmt.Println("Must provide a command. Run with --help to see command line options")
		return
	}

	keys, err := NewKeys()
	if err != nil {
		fmt.Printf("Failed to create keyring, %v", err)
		return
	}

	if options.Add != nil {
		err = Add(*options.Add, keys)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if options.List != nil {
		otps, err := keys.ListOTPs()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, otp := range otps {
			fmt.Printf("%s - %s\n", otp.Label, otp.OTP.ProvisioningUri(otp.Label, otp.Issuer))
		}
	}
}
