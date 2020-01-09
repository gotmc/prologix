// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"
	"time"

	"github.com/gotmc/prologix"
	"github.com/tarm/serial"
)

func main() {
	cfg := serial.Config{
		Name:        "/dev/tty.usbserial-PXFJL0WD",
		Baud:        115200,
		ReadTimeout: time.Millisecond * 500,
	}
	port, err := serial.OpenPort(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	gpib, err := prologix.NewController(port, 4, true)
	addr, err := gpib.InstrumentAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GPIB instrument address = %d", addr)
	ver, err := gpib.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", ver)
}
