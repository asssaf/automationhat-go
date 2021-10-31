package automationhat

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/sn3218"
	"periph.io/x/host/v3/rpi"
)

var (
	LedADC1     = 0
	LedADC2     = 1
	LedADC3     = 2
	LedOutput1  = 3
	LedOutput2  = 4
	LedOutput3  = 5
	LedRelay1NO = 6
	LedRelay1NC = 7
	LedRelay2NO = 8
	LedRelay2NC = 9
	LedRelay3NO = 10
	LedRelay3NC = 11
	LedInput3   = 12
	LedInput2   = 13
	LedInput1   = 14
	LedWarn     = 15 // red
	LedComm     = 16 // blue
	LedPower    = 17 // green
)

type Opts struct {
	AutoLeds bool
}

var DefaultOpts = Opts{
	AutoLeds: true,
}

// Dev represents an Automation HAT
type Dev struct {
	opts    Opts
	output1 gpio.PinOut
	output2 gpio.PinOut
	output3 gpio.PinOut
	input1  gpio.PinIn
	input2  gpio.PinIn
	input3  gpio.PinIn
	leds    *sn3218.Dev
}

// NewAutomationHat returns a automationhat driver.
func NewAutomationHat(opts *Opts) (*Dev, error) {
	i2cPort, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		return nil, err
	}

	leds, err := sn3218.New(i2cPort)
	if err != nil {
		return nil, err
	}

	dev := &Dev{
		opts:    *opts,
		output1: rpi.P1_29, // GPIO 5
		output2: rpi.P1_32, // GPIO 12
		output3: rpi.P1_31, // GPIO 6
		input1:  rpi.P1_37, // GPIO 26
		input2:  rpi.P1_38, // GPIO 20
		input3:  rpi.P1_40, // GPIO 21
		leds:    leds,
	}

	if dev.opts.AutoLeds {
		if err := dev.leds.WakeUp(); err != nil {
			return nil, err
		}

		if err := dev.leds.SwitchAll(true); err != nil {
			return nil, err
		}

		if err := dev.leds.Brightness(LedPower, 0x01); err != nil {
			return nil, err
		}
	}

	return dev, nil
}

// GetOutput1 returns gpio.PinOut corresponding to output 1
func (d *Dev) GetOutput1() gpio.PinOut {
	return d.output1
}

// GetOutput2 returns gpio.PinOut corresponding to output 2
func (d *Dev) GetOutput2() gpio.PinOut {
	return d.output2
}

// GetOutput3 returns gpio.PinOut corresponding to output 3
func (d *Dev) GetOutput3() gpio.PinOut {
	return d.output3
}

// GetInput1 returns gpio.PinIn corresponding to input 1
func (d *Dev) GetInput1() gpio.PinIn {
	return d.input1
}

// GetInput2 returns gpio.PinIn corresponding to input 2
func (d *Dev) GetInput2() gpio.PinIn {
	return d.input2
}

// GetInput3 returns gpio.PinIn corresponding to input 3
func (d *Dev) GetInput3() gpio.PinIn {
	return d.input3
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

	if err := d.input1.Halt(); err != nil {
		return err
	}

	if err := d.input2.Halt(); err != nil {
		return err
	}

	if err := d.input3.Halt(); err != nil {
		return err
	}

	if err := d.leds.Halt(); err != nil {
		return err
	}

	return nil
}
