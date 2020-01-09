// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package vcp

import (
	"io"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// VCP models a Prologix GPIB-USB controller communicating using a Virtual COM
// Port (VCP).
type VCP struct {
	port *serial.Port
}

// NewVCP creates a new Virtual COM Port (VCP).
func NewVCP(serialPort string) (*VCP, error) {
	cfg := &serial.Config{
		Name:        serialPort,
		Baud:        115200,
		ReadTimeout: time.Millisecond * 500,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	vcp := VCP{
		port: port,
	}
	return &vcp, nil
}

// Write writes the given data to the serial port.
func (vcp *VCP) Write(p []byte) (n int, err error) {
	return vcp.port.Write(p)
}

// Read reads from the serial port into the given byte slice.
func (vcp *VCP) Read(p []byte) (n int, err error) {
	return vcp.port.Read(p)
}

// Close closes the underlying serial port.
func (vcp *VCP) Close() error {
	return vcp.port.Close()
}

// WriteString trims all whitespace, adds a newline, and then writes the
// string using the underlying serial port.
func (vcp *VCP) WriteString(s string) (n int, err error) {
	s = strings.TrimSpace(s) + "\n"
	return io.WriteString(vcp.port, s)
}
