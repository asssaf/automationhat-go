package automationhat

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/sn3218"
	"periph.io/x/host/v3/rpi"
)

// Dev represents an Automation HAT
type Dev struct {
	output1 gpio.PinOut
	output2 gpio.PinOut
	output3 gpio.PinOut
	leds    *sn3218.Dev
}

// NewAutomationHat returns a automationhat driver.
func NewAutomationHat() (*Dev, error) {
	i2cPort, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		return nil, err
	}

	leds, err := sn3218.New(i2cPort)
	if err != nil {
		return nil, err
	}

	dev := &Dev{
		output1: rpi.P1_29, // GPIO 5
		output2: rpi.P1_32, // GPIO 12
		output3: rpi.P1_31, // GPIO 6
		leds:    leds,
	}

	return dev, nil
}

// GetOutput1 returns gpio.PinOut corresponding to output 1
func (d *Dev) GetOutput1() gpio.PinOut {
	return d.output1
}

// GetOutput1 returns gpio.PinOut corresponding to output 2
func (d *Dev) GetOutput2() gpio.PinOut {
	return d.output2
}

// GetOutput1 returns gpio.PinOut corresponding to output 3
func (d *Dev) GetOutput3() gpio.PinOut {
	return d.output3
}

func (d *Dev) GetLeds() *sn3218.Dev {
	return d.leds
}

// Halt all internal devices.
func (d *Dev) Halt() error {
	if err := d.output1.Halt(); err != nil {
		return err
	}

	if err := d.output2.Halt(); err != nil {
		return err
	}

	if err := d.output3.Halt(); err != nil {
		return err
	}

	if err := d.leds.Halt(); err != nil {
		return err
	}

	return nil
}
