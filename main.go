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

type RemoveCmd struct {
	Id   string `arg:"--id"`
	Name string `arg:"--name"`
}

type CodeCmd struct {
	Id      string `arg:"--id"`
	Name    string `arg:"positional"`
	Counter int64  `arg:"--counter" default:"-1"`
}

type DecodeCmd struct {
	Uri string `arg:"positional,required"`
}

var options struct {
	List   *ListCmd   `arg:"subcommand:list" help:"List stored OTPs"`
	Add    *AddCmd    `arg:"subcommand:add"`
	Remove *RemoveCmd `arg:"subcommand:remove"`
	Code   *CodeCmd   `arg:"subcommand:code"`
	Decode *DecodeCmd `arg:"subcommand:decode"`
}

// todo:
// - add command to set counter for HOTP
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
