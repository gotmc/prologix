// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"

	"github.com/gotmc/prologix/usb/vcp"
)

func main() {
	// dev, err := vcp.NewDevice("/dev/tty.usbserial-PX8X3YR6", 4)
	dev, err := vcp.NewDevice("/dev/tty.usbserial-PXFJL0WD", 4)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := dev.Address()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GPIB address = %d", addr)
	ver, err := dev.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", ver)
}
