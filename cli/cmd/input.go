package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"

	"github.com/asssaf/automationhat-go/automationhat"
)

type InputCommand struct {
	fs    *flag.FlagSet
	input int
}

func NewInputCommand() *InputCommand {
	c := &InputCommand{
		fs: flag.NewFlagSet("input", flag.ExitOnError),
	}

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: automationhat %s <1|2|3>\n", c.fs.Name())
	}

	return c
}

func (c *InputCommand) Name() string {
	return c.fs.Name()
}

func (c *InputCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	if c.fs.NArg() < 1 {
		return errors.New("Missing required arguments")
	}

	if c.fs.NArg() > 1 {
		return errors.New("Too many arguments")
	}

	input := c.fs.Arg(0)

	availableInputs := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	inputNum, ok := availableInputs[input]
	if !ok {
		return errors.New("unrecognized input, must be one of: [1, 2, 3]")
	}

	c.input = inputNum

	return nil
}

func (c *InputCommand) Execute() error {
	fmt.Printf("Reading input %d\n", c.input)

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	hat, err := automationhat.NewAutomationHat(&automationhat.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer hat.Halt()

	var input gpio.PinIn
	switch c.input {
	case 1:
		input = hat.GetInput1()
	case 2:
		input = hat.GetInput2()
	case 3:
		input = hat.GetInput3()
	}

	err = input.In(gpio.PullNoChange, gpio.NoEdge)
	if err != nil {
		log.Fatal(err)
	}

	value := input.Read()
	fmt.Printf("Value read: %s\n", value)

	return nil
}
