// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gotmc/prologix"
	"github.com/tarm/serial"
)

func main() {
	// Open a serial port.
	cfg := serial.Config{
		Name:        "/dev/tty.usbserial-PXFJL0WD",
		Baud:        115200,
		ReadTimeout: time.Millisecond * 1000,
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

	cmds := []string{
		"*rst",
		"*cls",
		"apply p6v, 5.0, 1.0",
		"outp on",
	}

	for _, cmd := range cmds {
		_, err := gpib.WriteString(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// queries := []string{"*idn"}
	// queryRange(gpib, queries)

	io.WriteString(port, "*idn?\n")
	io.WriteString(port, "++read eoi\n")
	buf := make([]byte, 8)
	n, err := port.Read(buf)
	if err != nil {
		log.Printf("error = %s", err)
	}
	if err == io.EOF {
		log.Printf("EOF error; read %d bytes = %s", n, buf[:n])
	} else {
		log.Printf("read %d bytes = %s", n, buf[:n])
	}

	// Close the serial port and check for errors.
	err = port.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
}

func queryRange(gpib *prologix.Controller, r []string) {
	for _, q := range r {
		ws := fmt.Sprintf("%s?", q)
		log.Printf("Querying %s", ws)
		s, err := gpib.Query(ws)
		log.Printf("Completed %s query", ws)
		if err != nil && err != io.EOF {
			log.Printf("Error reading: %v", err)
		} else if err == io.EOF {
			log.Printf("got EOF")
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
