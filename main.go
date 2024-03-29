package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/uaraven/clotp/cli"
	"github.com/uaraven/clotp/keyrings"
)

type Options struct {
	List       *cli.ListCmd       `arg:"subcommand:list" help:"List stored OTPs"`
	Add        *cli.AddCmd        `arg:"subcommand:add" help:"Add new OTP code"`
	Remove     *cli.RemoveCmd     `arg:"subcommand:remove" help:"Remove existing OTP code"`
	Code       *cli.CodeCmd       `arg:"subcommand:code" help:"Generate OTP code"`
	Copy       *cli.CopyCmd       `arg:"subcommand:copy" help:"Copy OTP code to the clipboard"`
	Decode     *cli.DecodeCmd     `arg:"subcommand:decode" help:"Decode Google Authenticator migration URI"`
	SetCounter *cli.SetCounterCmd `arg:"subcommand:set-counter" help:"Set HOTP counter"`
	View       *cli.ViewCmd       `arg:"subcommand:view" help:"View OTP code details"`
	Scan       *cli.ScanCmd       `arg:"subcommand:scan" help:"Scan and decode QR code"`
	Gen        *cli.GenCmd        `arg:"subcommand:gen" help:"Genereate TOTP code on the fly"`
}

func (Options) Version() string {
	return "CLotp 0.1.2"
}

var options Options

// todo:
// - add parameter to generate TOTP for a given timestamp

func main() {
	arg.MustParse(&options)

	keys, err := keyrings.NewKeys()
	if err != nil {
		fmt.Printf("Failed to create keyring, %v\n", err)
		return
	}
	var output string
	err = nil
	if options.Add != nil {
		output, err = cli.Add(options.Add, keys)
	} else if options.List != nil {
		output, err = cli.List(options.List, keys)
	} else if options.Remove != nil {
		output, err = cli.Remove(options.Remove, keys)
	} else if options.Code != nil {
		output, err = cli.Code(options.Code, keys)
	} else if options.Copy != nil {
		output, err = cli.Copy(options.Copy, keys)
	} else if options.Decode != nil {
		output, err = cli.Decode(options.Decode)
	} else if options.SetCounter != nil {
		output, err = cli.SetCounter(options.SetCounter, keys)
	} else if options.View != nil {
		output, err = cli.View(options.View, keys)
	} else if options.Scan != nil {
		output, err = cli.ScanQrCode(options.Scan)
	} else if options.Gen != nil {
		output, err = cli.Gen(options.Gen)
	} else {
		fmt.Println("Must provide a command. Run with --help to see command line options")
		return
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println(output)
	}
}
