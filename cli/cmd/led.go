package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"periph.io/x/host/v3"

	"github.com/asssaf/automationhat-go/automationhat"
)

type LedCommand struct {
	fs    *flag.FlagSet
	led   int
	value int
}

func NewLedCommand() *LedCommand {
	c := &LedCommand{
		fs: flag.NewFlagSet("led", flag.ExitOnError),
	}

	c.fs.IntVar(&c.led, "led", 0, "Led to set (0-17)")
	c.fs.IntVar(&c.value, "brightness", 0, "Brightness to set (0-255)")

	return c
}

func (c *LedCommand) Name() string {
	return c.fs.Name()
}

func (c *LedCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if c.led < 0 || c.led > 17 {
		return errors.New(fmt.Sprintf("Led number out of range: %d", c.led))
	}

	if c.value < 0 || c.value > 255 {
		return errors.New(fmt.Sprintf("Led brightness out of range: %d", c.led))
	}
	return nil
}

func (c *LedCommand) Execute() error {
	fmt.Printf("Turning led %d to %d\n", c.led, c.value)

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	hat, err := automationhat.NewAutomationHat()
	if err != nil {
		log.Fatal(err)
	}

	//defer hat.Halt()
	// calling halt will reset the chip

	leds := hat.GetLeds()

	if on, brightness, err := leds.GetState(1); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("state %t %d\n", on, brightness)
	}

	if err := leds.SwitchAll(true); err != nil {
		log.Fatal(err)
	}

	if err := leds.Brightness(c.led, byte(c.value)); err != nil {
		log.Fatal(err)
	}

	if err := leds.Switch(c.led, c.value != 0); err != nil {
		log.Fatal(err)
	}

	if err := leds.WakeUp(); err != nil {
		log.Fatal(err)
	}

	if on, brightness, err := leds.GetState(1); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("state %t %d\n", on, brightness)
	}

	return nil
}
