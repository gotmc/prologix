// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"io"
	"log"
	"time"

	"github.com/gotmc/prologix"
	"github.com/tarm/serial"
)

func main() {
	// Open a serial port.
	cfg := serial.Config{
		Name:        "/dev/tty.usbserial-PX8X3YR6",
		Baud:        115200,
		ReadTimeout: time.Millisecond * 500,
	}
	port, err := serial.OpenPort(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GPIB controller using the aforementioned serial port
	// communicating with the instrument at the given GPIB address.
	gpib, err := prologix.NewController(port, 5, true)

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

	// Send the Selected Device Clear (SDC) message
	err = gpib.ClearDevice()
	if err != nil {
		log.Printf("error clearing device: %s", err)
	}

	cmds := []string{
		"apply p6v, 5.0, 1.0",
		"outp on",
	}

	for _, cmd := range cmds {
		_, err := gpib.WriteString(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Query the identification of the function generator.
	idn, err := gpib.Query("*idn?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	log.Printf("query idn = %s", idn)

	// Query the identification of the function generator again.
	idn, err = gpib.Query("*idn?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	log.Printf("query idn = %s", idn)

	// Return local control to the front panel.
	err = gpib.FrontPanel(true)
	if err != nil {
		log.Fatalf("error setting local control for front panel: %s", err)
	}

	// Discard any unread data on the serial port and then close.
	err = port.Flush()
	if err != nil {
		log.Printf("error flushing serial port: %s", err)
	}
	err = port.Close()
	if err != nil {
		log.Printf("error closing serial port: %s", err)
	}
}
