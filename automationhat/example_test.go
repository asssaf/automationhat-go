package automationhat_test

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"periph.io/x/host/v3"

	"automationhat"
)

func Example() {

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	hat, err := automationhat.NewAutomationHat()
	if err != nil {
		log.Fatal(err)
	}
	defer hat.Halt()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigs // Wait for signal
		log.Println(sig)
		done <- true
	}()

	log.Println("Press ctrl+c to stop...")
	<-done // Wait
}
