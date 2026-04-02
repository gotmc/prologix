// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package vcp

import (
	"context"
	"io"
	"strings"
	"time"

	"go.bug.st/serial"
)

// VCP models a Prologix GPIB-USB controller communicating using a Virtual COM
// Port (VCP).
type VCP struct {
	port serial.Port
}

// NewVCP creates a new Virtual COM Port (VCP).
func NewVCP(serialPort string) (*VCP, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(serialPort, mode)
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

// ReadContext reads from the serial port into the given byte slice. If the
// context has a deadline, the serial port's read timeout is set accordingly. If
// the context is cancelled before the read begins, the context error is
// returned immediately.
func (vcp *VCP) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		if timeout <= 0 {
			return 0, context.DeadlineExceeded
		}
		if err := vcp.port.SetReadTimeout(timeout); err != nil {
			return 0, err
		}
		defer func() { _ = vcp.port.SetReadTimeout(serial.NoTimeout) }()
	}
	n, err = vcp.port.Read(p)
	if err != nil {
		return n, err
	}
	if ctx.Err() != nil {
		return n, ctx.Err()
	}
	return n, nil
}

// WriteContext writes the given data to the serial port. If the context is
// already done, the context error is returned immediately. Because the serial
// port interface does not support write timeouts, the write is performed in a
// goroutine so the context cancellation is respected.
func (vcp *VCP) WriteContext(ctx context.Context, p []byte) (n int, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	type result struct {
		n   int
		err error
	}
	ch := make(chan result, 1)
	go func() {
		n, err := vcp.port.Write(p)
		ch <- result{n, err}
	}()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case r := <-ch:
		return r.n, r.err
	}
}

// Close closes the underlying serial port.
func (vcp *VCP) Close() error {
	return vcp.port.Close()
}

func (vcp *VCP) Flush() error {
	err := vcp.port.ResetInputBuffer()
	if err != nil {
		return err
	}
	return vcp.port.ResetOutputBuffer()
}

// WriteString trims all whitespace, adds a newline, and then writes the
// string using the underlying serial port.
func (vcp *VCP) WriteString(s string) (n int, err error) {
	s = strings.TrimSpace(s) + "\n"
	return io.WriteString(vcp.port, s)
}
