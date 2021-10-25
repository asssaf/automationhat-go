package cmd

import (
	"errors"
	"flag"
	"fmt"
)

type execFunc func([]string) error

type Command interface {
	Init([]string) error
	Execute() error
	Name() string
}

var (
	outputCmd = flag.NewFlagSet("output", flag.ExitOnError)
)

var rootCmds = map[string]execFunc{}

func Execute() error {
	commands := []Command{
		NewOutputCommand(),
		NewLedCommand(),
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: automationhat <command> ...\n")
		fmt.Fprintf(flag.CommandLine.Output(), "The commands are:\n")
		for _, c := range commands {
			fmt.Fprintf(flag.CommandLine.Output(), "\t%s\n", c.Name())
		}
		flag.PrintDefaults()
	}

	flag.Parse()

	command := flag.Arg(0)
	if command == "" {
		return errors.New("Missing command")
	}

	args := flag.Args()
	for _, c := range commands {
		if command == c.Name() {
			if err := c.Init(args[1:]); err != nil {
				return err
			}
			return c.Execute()
		}
	}

	return errors.New(fmt.Sprintf("unknown command: %s", command))
}
