// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package vcp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// Querier provides the interface to query using a given string and provide the
// resultant string.
type Querier interface {
	Query(s string) (value string, err error)
}

// Device models a GPIB device connected using the Virtual COM Port (VCP) with
// the Prologix GPIB-USB Controller.
type Device struct {
	primaryAddress int
	port           *serial.Port
}

// NewDevice opens a GPIB device at the given primary address using the given
// serial port for the Prologix GPIB-USB controller using the Virtual COM port.
func NewDevice(serialPort string, primaryAddress int) (*Device, error) {
	c := &serial.Config{Name: serialPort, Baud: 115200, ReadTimeout: time.Millisecond * 500}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	d := Device{
		port: port,
	}
	addrCmd := fmt.Sprintf("++addr %d", primaryAddress)
	cmds := []string{
		"++savecfg 0",    // Disable saving of configuration parameters in EPROM
		addrCmd,          // Set the primary address.
		"++mode 1",       // Switch to controller mode.
		"++auto 0",       // Turn off read-after-write and address instrument to listen.
		"++eoi 1",        // Enable EOI assertion with last character.
		"++eos 0",        // Set GPIB termination to append CR+LF to instrument commands.
		"++eot_enable 0", // Do not append character when EOI detected
		"++savecfg 1",    // Enable saving of configuration parameters in EPROM
		"++clr",
	}
	for _, cmd := range cmds {
		_, err = d.WriteString(cmd)
		if err != nil {
			return nil, err
		}
	}
	return &d, nil
}

// Write writes the given data to the serial port, which will be sent to the
// GPIB using the current GPIB address.
func (d *Device) Write(p []byte) (n int, err error) {
	return d.port.Write(p)
}

// Read reads from the serial port into the given byte slice.
func (d *Device) Read(p []byte) (n int, err error) {
	return d.port.Read(p)
}

// Close closes the underlying serial port.
func (d *Device) Close() error {
	return d.port.Close()
}

// WriteString trimes all whitespace, adds a newline, and then writes the
// string using the underlying serial port.
func (d *Device) WriteString(s string) (n int, err error) {
	s = strings.TrimSpace(s) + "\n"
	return io.WriteString(d.port, s)
}

// Query queries the device sending the given string and then reading the
// response.
func (d *Device) Query(s string) (string, error) {
	_, err := d.WriteString(s)
	if err != nil {
		return "", err
	}
	return bufio.NewReader(d.port).ReadString('\n')
}

// QueryInt is used to query a Querier interface and return an int.
func QueryInt(q Querier, query string) (int, error) {
	s, err := q.Query(query)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return int(i), err
}

// SetAddress sets the primary address of the GPIB controller.
func (d *Device) SetAddress(primaryAddress int) error {
	s := fmt.Sprintf("++addr %d", primaryAddress)
	_, err := d.WriteString(s)
	return err
}

// Address displays the GPIB address of the GPIB controller.
func (d *Device) Address() (int, error) {
	return QueryInt(d, "++addr")
}

// FrontPanelControl enables or disables front panel operation.
func (d *Device) FrontPanelControl(enable bool) error {
	cmd := "++llo"
	if enable {
		cmd = "++loc"
	}
	_, err := d.WriteString(cmd)
	return err
}

// Version returns the version string of the Prologix GPUB-USB controller.
func (d *Device) Version() (string, error) {
	return d.Query("++ver")
}
