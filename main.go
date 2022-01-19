package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

type ListCmd struct {
	Parse bool `arg:"-p,--parse" help:"Parse URIs and print each part separately"`
}

type AddCmd struct {
	Uri  string `arg:"positional,required"`
	Name string `arg:"--name" help:"Optional name of the code to refer to it later"`
}

type RemoveCmd struct {
	Id   string `arg:"--id" help:"ID of the code to remove"`
	Name string `arg:"--name" help:"Name of the code to remove. Either name or ID must be provided"`
}

type CodeCmd struct {
	Id      string `arg:"--id" help:"Look up key by its ID, instead of name"`
	Name    string `arg:"positional"`
	Counter int64  `arg:"--counter" help:"Override counter for the HOTP code"`
	Copy    bool   `arg:"--copy" help:"Copy generated code to clipboard"`
}

type DecodeCmd struct {
	Uri   string `arg:"positional,required" help:"Google Authenticator Export URI"`
	Parse bool   `arg:"-p,--parse" help:"Parse URIs and print each part separately"`
}

type SetCounterCmd struct {
	Id      string `arg:"--id" help:"HOTP code identifier"`
	Name    string `arg:"--name" help:"HOTP code name"`
	Counter int64  `arg:"positional,required" help:"New counter value"`
}

var options struct {
	List       *ListCmd       `arg:"subcommand:list" help:"List stored OTPs"`
	Add        *AddCmd        `arg:"subcommand:add"`
	Remove     *RemoveCmd     `arg:"subcommand:remove"`
	Code       *CodeCmd       `arg:"subcommand:code"`
	Decode     *DecodeCmd     `arg:"subcommand:decode"`
	SetCounter *SetCounterCmd `arg:"subcommand:set-counter" help:"Set HOTP counter"`
}

// todo:
// - add parameter to generate HOTP for a given counter or TOTP for a given timestamp

func main() {
	arg.MustParse(&options)

	keys, err := NewKeys()
	if err != nil {
		fmt.Printf("Failed to create keyring, %v\n", err)
		return
	}
	var output string
	err = nil
	if options.Add != nil {
		output, err = Add(options.Add, keys)
	} else if options.List != nil {
		output, err = List(options.List, keys)
	} else if options.Remove != nil {
		output, err = Remove(options.Remove, keys)
	} else if options.Code != nil {
		output, err = Code(options.Code, keys)
	} else if options.Decode != nil {
		output, err = Decode(options.Decode)
	} else if options.SetCounter != nil {
		output, err = SetCounter(options.SetCounter, keys)
	} else {
		fmt.Println("Must provide a command. Run with --help to see command line options")
		return
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}
}
