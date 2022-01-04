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

var options struct {
	List   *ListCmd   `arg:"subcommand:list"`
	Add    *AddCmd    `arg:"subcommand:add"`
	Remove *RemoveCmd `arg:"subcommand:remove"`
}

func main() {
	arg.MustParse(&options)

	keys, err := NewKeys()
	if err != nil {
		fmt.Printf("Failed to create keyring, %v\n", err)
		return
	}
	err = nil
	if options.Add != nil {
		err = Add(options.Add, keys)
	} else if options.List != nil {
		err = List(options.List, keys)
	} else if options.Remove != nil {
		err = Remove(options.Remove, keys)
	} else {
		fmt.Println("Must provide a command. Run with --help to see command line options")
		return
	}
	if err != nil {
		fmt.Println(err)
	}
}
