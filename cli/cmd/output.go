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

type OutputCommand struct {
	fs     *flag.FlagSet
	output int
	action bool
}

func NewOutputCommand() *OutputCommand {
	c := &OutputCommand{
		fs: flag.NewFlagSet("output", flag.ExitOnError),
	}

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: automationhat %s <1|2|3> <on|off>\n", c.fs.Name())
	}

	return c
}

func (c *OutputCommand) Name() string {
	return c.fs.Name()
}

func (c *OutputCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	if c.fs.NArg() < 2 {
		return errors.New("Missing required arguments")
	}

	if c.fs.NArg() > 2 {
		return errors.New("Too many arguments")
	}

	output := c.fs.Arg(0)
	action := c.fs.Arg(1)

	availableOutputs := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	outputNum, ok := availableOutputs[output]
	if !ok {
		return errors.New("unrecognized output, must be one of: [1, 2, 3]")
	}

	c.output = outputNum

	switch action {
	case "on":
		c.action = true
	case "off":
		c.action = false
	default:
		return errors.New("unrecognized action, must be one of: [on, off]")
	}

	return nil
}

func (c *OutputCommand) Execute() error {
	actionName := "off"
	if c.action == true {
		actionName = "on"
	}
	fmt.Printf("Turning output %d %s\n", c.output, actionName)

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	hat, err := automationhat.NewAutomationHat(&automationhat.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer hat.Halt()

	var output gpio.PinOut
	switch c.output {
	case 1:
		output = hat.GetOutput1()
	case 2:
		output = hat.GetOutput2()
	case 3:
		output = hat.GetOutput3()
	}

	value := gpio.Low
	if c.action {
		value = !value
	}

	if err := output.Out(value); err != nil {
		log.Fatal(err)
	}

	return nil
}
