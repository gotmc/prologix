// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gotmc/prologix"
	"github.com/gotmc/prologix/driver/vcp"
)

func main() {
	// Open virtual comm port.
	serialPort := "/dev/tty.usbserial-PXFJL0WD"
	vcp, err := vcp.NewVCP(serialPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GPIB controller using the aforementioned serial port
	// communicating with the instrument at the given GPIB address.
	gpib, err := prologix.NewController(vcp, 10, true)
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

	// Determine if ReadAfterWrite is enabled.
	enabled, err := gpib.ReadAfterWrite()
	if err != nil {
		log.Printf("error determining ReadAfterWrite: %s", err)
	}
	log.Printf("ReadAfterWrite = %t", enabled)

	// Determine if EOI is asserted
	asserted, err := gpib.AssertEOI()
	if err != nil {
		log.Printf("error determining AssertEOI: %s", err)
	}
	log.Printf("AssertEOI = %t", asserted)

	// Set the GPIB termination
	err = gpib.SetGPIBTermination(prologix.AppendCRLF)
	if err != nil {
		log.Printf("error setting the GPIB termination: %s", err)
	}

	// Determine GPIB termination
	term, err := gpib.GPIBTermination()
	if err != nil {
		log.Printf("error determining GPIB termination: %s", err)
	}
	log.Printf("GPIB termination = %s", term)

	// Send the Selected Device Clear (SDC) message
	err = gpib.ClearDevice()
	if err != nil {
		log.Printf("error clearing device: %s", err)
	}
	time.Sleep(time.Millisecond * 500)

	cmds := []string{
		"vac",
		"vdc",
		"ohms",
		"rate s",
		"range 1",
	}

	for _, cmd := range cmds {
		_, err = gpib.WriteString(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Query the identification of the DMM
	log.Printf("Querying with ReadAfterWrite disabled")
	err = gpib.SetReadAfterWrite(false)
	if err != nil {
		log.Fatal(err)
	}
	idn, err := gpib.Query("*idn?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	} else if err == io.EOF {
		log.Printf("received EOF")
	}
	log.Printf("query idn = %s", idn)

	// Query the identification of the DMM
	log.Printf("Querying with ReadAfterWrite enabled")
	err = gpib.SetReadAfterWrite(true)
	if err != nil {
		log.Fatal(err)
	}
	idn, err = gpib.Query("*idn?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	} else if err == io.EOF {
		log.Printf("received EOF")
	}
	log.Printf("query idn = %s", idn)

	// Measure the resistance
	resString, err := gpib.Query("meas1?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	} else if err == io.EOF {
		log.Printf("received EOF")
	}
	resString = strings.TrimSpace(resString)
	res, err := strconv.ParseFloat(resString, 64)
	if err != nil {
		log.Printf("error converting %s into a float64", resString)
	}
	log.Printf("resistance = %s ohms (string)", resString)
	log.Printf("resistance = %f ohms", res)

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
