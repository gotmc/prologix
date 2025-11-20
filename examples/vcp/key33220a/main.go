// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"io"
	"log"
	"time"

	"github.com/gotmc/prologix"
	"github.com/gotmc/prologix/driver/vcp"
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

	flag.IntVar(&gpibAddress, "gpib", 6, "GPIB address for the Keysight 33220A")
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
	log.Printf("Create new Prolgoix controller using GPIB address %d",
		gpibAddress)
	gpib, err := prologix.NewController(vcp, gpibAddress, false)
	if err != nil {
		log.Fatalf("NewController error: %s", err)
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

	// Send commands to the function generator required to create a coded carrier
	// operating at 100 Hz with 400 ms on time and 200 ms off time.
	cmds := []string{
		"SYST:REM", // Set the instrument state to remote
		"OUTP OFF", // Disable the output
		// "BURS:STAT OFF",        // Disable burst mode (i.e., modulation)
		"APPL:SIN 100,0.5,0.0", // 100 Hz, 0.5 Vpp amp, 0.0Vdc offset
		"BURS:MODE TRIG",       // Select the triggered burst mode
		"BURS:NCYC 40",         // Set burst count to 40 (nearest integer)
		"BURS:INT:PER 0.6",     // Set burst period to 600 ms
		"BURS:PHAS 0",          // Set starting phase of burst to 0 degrees
		"BURS:STAT ON",         // Enable burst mode (i.e., modulation)
		"OUTP ON",              // Enable the output
		"SYST:LOC",             // Set the instrument state to local
	}
	for _, cmd := range cmds {
		log.Printf("Sending command: %s", cmd)
		err = gpib.Command(cmd)
		if err != nil {
			log.Fatal(err)
		}
		// FIXME: I had to enter this to get the function generator to accept all
		// commands.
		time.Sleep(250 * time.Millisecond)
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
