// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package prologix

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// Controller models a GPIB controller-in-charge.
type Controller struct {
	rw               io.ReadWriter
	primaryAddr      int
	hasSecondaryAddr bool
	secondaryAddr    int
	auto             bool
	eoi              bool
	usbTerm          byte
	eotChar          byte
}

// ControllerOption applies an option to the controller.
type ControllerOption func(*Controller)

// NewController creates a GPIB controller-in-charge at the given address using
// the given Prologix driver, which can either be a Virtual COM Port (VCP), USB
// direct, or Ethernet. Enable clear to send the Selected Device Clear (SDC)
// message to the GPIB address. Optionally controller configuration can be
// included using a ControllerOption.
func NewController(
	rw io.ReadWriter,
	addr int,
	clear bool,
	opts ...ControllerOption,
) (*Controller, error) {
	c := Controller{
		rw:               rw,
		primaryAddr:      addr,
		hasSecondaryAddr: false,
		auto:             false,
		eoi:              true,
		usbTerm:          '\n',
		eotChar:          '\n',
	}

	// Apply options using the functional option pattern.
	for _, opt := range opts {
		opt(&c)
	}

	// Verify validate primary address.
	if !isPrimaryAddressValid(c.primaryAddr) {
		return nil, fmt.Errorf("invalid primary address %d (must by 0-30)", c.primaryAddr)
	}

	// Configure the Prologix GPIB controller.
	addrCmd := fmt.Sprintf("addr %d", c.primaryAddr)
	if c.hasSecondaryAddr {
		if !isSecondaryAddressValid(c.secondaryAddr) {
			return nil, fmt.Errorf("invalid secondary address %d (must be 96-126)", c.secondaryAddr)
		}
		addrCmd = fmt.Sprintf("addr %d %d", c.primaryAddr, c.secondaryAddr)
	}
	eotCharCmd := fmt.Sprintf("eot_char %d", c.eotChar)
	cmds := []string{
		"savecfg 0",       // Disable saving of configuration parameters in EPROM
		addrCmd,           // Set the primary address.
		"mode 1",          // Switch to controller mode.
		"auto 0",          // Turn off read-after-write and address instrument to listen.
		"eoi 1",           // Enable EOI assertion with last character.
		"eos 0",           // Set GPIB termination.
		"read_tmo_ms 500", // Set the read timeout to 500 ms.
		eotCharCmd,        // Set the EOT char
		"eot_enable 1",    // Append character when EOI detected?
		"savecfg 1",       // Enable saving of configuration parameters in EPROM
	}
	if clear {
		cmds = append(cmds, "clr")
	}
	for _, cmd := range cmds {
		if err := c.CommandController(cmd); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// WithSecondaryAddress sets a secondary address, which must be in the range of
// 96 and 126, inclusive.
func WithSecondaryAddress(addr int) ControllerOption {
	return func(c *Controller) {
		c.hasSecondaryAddr = true
		c.secondaryAddr = addr
	}
}

// Write writes the given data to the instrument at the currently assigned GPIB
// address.
func (c *Controller) Write(p []byte) (n int, err error) {
	return c.rw.Write(p)
}

// Read reads from the instrument at the currently assigned GPIB address into
// the given byte slice.
func (c *Controller) Read(p []byte) (n int, err error) {
	return c.rw.Read(p)
}

// WriteString writes a string to the instrument at the currently assigned GPIB
// address.
func (c *Controller) WriteString(s string) (n int, err error) {
	cmd := fmt.Sprintf("%s%c", strings.TrimSpace(s), c.usbTerm)
	log.Printf("prologix driver writing string: %s", cmd)
	return c.rw.Write([]byte(cmd))
}

// Command formats according to a format specifier if provided and sends a
// SCPI/ASCII command to the instrument at the currently assigned GPIB address.
// All leading and trailing whitespace is removed before appending the USB
// terminator to the command sent to the Prologix.
func (c *Controller) Command(format string, a ...any) error {
	cmd := format
	if a != nil {
		cmd = fmt.Sprintf(format, a...)
	}
	// log.Printf("sending cmd (terminator not yet added): %#v", cmd)
	// TODO: Why am I trimming whitespace and adding the USB terminator here if
	// I'm calling the WriteString method, which does that as well?
	cmd = fmt.Sprintf("%s%c", strings.TrimSpace(cmd), c.usbTerm)
	// log.Printf("sending cmd (with terminator added): %#v", cmd)
	_, err := fmt.Fprint(c.rw, cmd)
	return err
}

// Query queries the instrument at the currently assigned GPIB using the given
// SCPI/ASCII command. The cmd string does not need to include a new line
// character, since all leading and trailing whitespace is removed before
// appending the USB terminator to the command sent to the Prologix.  When data
// from host is received over USB, the Prologix controller removes all
// non-escaped LF, CR and ESC characters and appends the GPIB terminator, as
// specified by the `eos` command, before sending the data to instruments.  To
// change the GPIB terminator use the SetGPIBTermination method.
func (c *Controller) Query(cmd string) (string, error) {
	cmd = fmt.Sprintf("%s%c", strings.TrimSpace(cmd), c.usbTerm)
	// log.Printf("sending query cmd: %#v", cmd)
	_, err := fmt.Fprint(c.rw, cmd)
	if err != nil {
		return "", fmt.Errorf("error writing command: %s", err)
	}
	// If read-after-write is disabled, need to tell the Prologix controller to
	// read.
	if !c.auto {
		readCmd := "++read eoi"
		_, err = fmt.Fprintf(c.rw, "%s%c", readCmd, c.usbTerm)
		if err != nil {
			return "", fmt.Errorf("error sending `%s` command: %s", readCmd, err)
		}
	}
	s, err := bufio.NewReader(c.rw).ReadString(c.eotChar)
	if err == io.EOF {
		log.Printf("found EOF")
		return s, nil
	}
	return s, err
}

// QueryController sends the given command to the Prologix controller and
// returns its response as a string. To indicate this is a command for the
// Prologix controller, thereby not transmitting over GPIB, two plus signs `++`
// are prepended. Addtionally, a new line is appended to act as the USB
// termination character.
func (c *Controller) QueryController(cmd string) (string, error) {
	_, err := fmt.Fprintf(c.rw, "++%s%c", strings.ToLower(strings.TrimSpace(cmd)), c.usbTerm)
	if err != nil {
		return "", err
	}
	return bufio.NewReader(c.rw).ReadString(c.eotChar)
}

// CommandController sends the given command to the Prologix controller. To
// indicate this is a command for the Prologix controller, thereby not
// transmitting to the instrument over GPIB, two plus signs `++` are prepended.
// Addtionally, a new line is appended to act as the USB termination character.
func (c *Controller) CommandController(cmd string) error {
	_, err := fmt.Fprintf(c.rw, "++%s%c", strings.ToLower(strings.TrimSpace(cmd)), c.usbTerm)
	return err
}

// GpibTerm provides the type for the available GPIB terminators.
type GpibTerm int

// Available GPIB terminators for the Prologix Controller.
const (
	AppendCRLF GpibTerm = iota
	AppendCR
	AppendLF
	AppendNothing
)

var gpibTermDesc = map[GpibTerm]string{
	AppendCRLF:    `Append CR+LF (\r\n) to instrument commands`,
	AppendCR:      `Append CR (\r) to instrument commands`,
	AppendLF:      `Append LF (\n) to instrument commands`,
	AppendNothing: `Do not append anything to instrument commands`,
}

func (term GpibTerm) String() string {
	return gpibTermDesc[term]
}

// isPrimaryAddressValid checks that the primary GPIB address is between 0 and
// 30, inclusive.
func isPrimaryAddressValid(addr int) bool {
	if addr < 0 || addr > 30 {
		return false
	}
	return true
}

// isSecondaryAddressValid checks that the secondary GPIB address is between 96
// and 126, inclusive.
func isSecondaryAddressValid(addr int) bool {
	if addr < 96 || addr > 126 {
		return false
	}
	return true
}
