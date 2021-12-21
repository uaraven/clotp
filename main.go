package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

type ListCmd struct {
}

type AddCmd struct {
	Url string `arg:"positional,required"`
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

	if options.List != nil {
		keys, err := keys.ListKeys()
		if err != nil {
			fmt.Println(err)
			return
		}
		for key := range keys {
			fmt.Println(key)
		}
	}
}
