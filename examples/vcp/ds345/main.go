// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"io"
	"log"

	"github.com/gotmc/prologix"
	"github.com/gotmc/prologix/driver/vcp"
)

func main() {
	// Open virtual comm port.
	serialPort := "/dev/tty.usbserial-PX8X3YR6"
	vcp, err := vcp.NewVCP(serialPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GPIB controller using the aforementioned serial port and
	// communicating with the instrument at GPIB address 4.
	gpib, err := prologix.NewController(vcp, 4, true)
	if err != nil {
		log.Fatalf("NewController error: %s", err)
	}

	// Query the GPIB instrument address.
	addr, err := gpib.InstrumentAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GPIB instrument address = %d", addr)

	// Query the Prologix controller version.
	ver, err := gpib.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", ver)

	// Query the auto mode (i.e., read after write).
	auto, err := gpib.ReadAfterWrite()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read after write = %t", auto)

	// Query the read timeout
	timeout, err := gpib.ReadTimeout()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read timeout = %d ms", timeout)

	// Determine if the SRQ is asserted.
	srq, err := gpib.ServiceRequest()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Service request asserted = %t", srq)

	// Query the GPIB Termination
	term, err := gpib.GPIBTermination()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", term)

	// Query the identification of the function generator.
	idn, err := gpib.Query("*idn?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	log.Printf("query idn = %s", idn)

	// Send commands to the function generator required to create a coded carrier
	// operating at 100 Hz with 400 ms on time and 200 ms off time.
	cmds := []string{
		"MENA0",      // Disable modulation
		"FUNC 0",     // Set output function to sine wave
		"FREQ 100.0", // Set frequency to 100 Hz
		"AMPL 0.5VP", // Set the amplitude to 0.5 Vpp
		"OFFS 0.0",   // Set the offset to 0 Vdc
		"PHSE 0.0",   // Set the phase to 0 degrees
		"BCNT 40",    // Set burst count to 40 (rounded to nearest integer)
		"TSRC1",      // Use internal trigger
		"TRAT 1.667", // Set the trigger rate to 1.667 Hz
		"MTYP5",      // Set the modulation to burst
		"MENA1",      // Enable modulation
	}
	for _, cmd := range cmds {
		log.Printf("Sending command: %s", cmd)
		err = gpib.Command(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Return local control to the front panel.
	err = gpib.FrontPanel(true)
	if err != nil {
		log.Fatalf("error setting local control for front panel: %s", err)
	}

	// Discard any unread data on the serial port and then close.
	err = vcp.Flush()
	if err != nil {
		log.Printf("error flushing serial port: %s", err)
	}
	err = vcp.Close()
	if err != nil {
		log.Printf("error closing serial port: %s", err)
	}
}