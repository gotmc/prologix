// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"io"
	"log"

	"github.com/gotmc/prologix"
	"github.com/gotmc/prologix/driver/vcp"
	"github.com/gotmc/query"
)

var (
	serialPort  string
	gpibAddress int
)

func init() {
	// Get Virtual COM Port (VCP) serial port for Prologix.
	flag.StringVar(
		&serialPort,
		"port",
		"/dev/tty.usbserial-PX8X3YR6",
		"Serial port for Prologix VCP GPIB controller",
	)

	flag.IntVar(&gpibAddress, "gpib", 5, "GPIB address for the E3631A")
}

func main() {
	// Parse the flags
	flag.Parse()

	// Open virtual comm port.
	log.Printf("Serial port = %s", serialPort)
	vcp, err := vcp.NewVCP(serialPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GPIB controller using the aforementioned serial port
	// communicating with the instrument at the given GPIB address.
	log.Printf("GPIB address = %d", gpibAddress)
	gpib, err := prologix.NewController(vcp, gpibAddress, true)
	if err != nil {
		log.Fatal(err)
	}

	// Query the GPIB instrument address.
	addr, _, err := gpib.InstrumentAddress()
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
	log.Println("Sending the Selected Device Clear (SDC) message")
	err = gpib.ClearDevice()
	if err != nil {
		log.Printf("error clearing device: %s", err)
	}

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

	cmds := []string{
		"outp off",
	}

	for _, cmd := range cmds {
		err = gpib.Command(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Query the output state
	state, err := query.Bool(gpib, "OUTP:STAT?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	if state {
		log.Println("output is enabled")
	} else {
		log.Println("output is disabled")
	}

	cmds = []string{
		"apply p6v,4.1,1.2",
		"outp on",
	}

	for _, cmd := range cmds {
		err = gpib.Command(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Query the voltage and current at the output
	vc, err := gpib.Query("appl? p6v")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	log.Printf("voltage, current = %s", vc)

	// Query the voltage at the output
	volt, err := gpib.Query("meas? p6v")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	log.Printf("voltage = %s", volt)

	// Query the output state
	state, err = query.Bool(gpib, "OUTP:STAT?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	if state {
		log.Println("output is enabled")
	} else {
		log.Println("output is disabled")
	}

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
	err = vcp.Flush()
	if err != nil {
		log.Printf("error flushing serial port: %s", err)
	}
	err = vcp.Close()
	if err != nil {
		log.Printf("error closing serial port: %s", err)
	}
}
