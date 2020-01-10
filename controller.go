// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package prologix

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// Controller models a GPIB controller-in-charge.
type Controller struct {
	rw             io.ReadWriter
	instrumentAddr int
	auto           bool
	eoi            bool
}

// NewController creates a GPIB controller-in-charge at the given address using
// the given Prologix driver, which can either be a Virtual COM Port (VCP), USB
// direct, or Ethernet. Enable clear to send the Selected Device Clear (SDC)
// message to the GPIB address.
func NewController(rw io.ReadWriter, addr int, clear bool) (*Controller, error) {
	c := Controller{
		rw:             rw,
		instrumentAddr: addr,
		auto:           false,
	}
	// Configure the Prologix GPIB controller.
	addrCmd := fmt.Sprintf("addr %d", addr)
	cmds := []string{
		"savecfg 0",       // Disable saving of configuration parameters in EPROM
		addrCmd,           // Set the primary address.
		"mode 1",          // Switch to controller mode.
		"auto 0",          // Turn off read-after-write and address instrument to listen.
		"eoi 1",           // Enable EOI assertion with last character.
		"eos 0",           // Set GPIB termination.
		"read_tmo_ms 500", // Set the read timeout to 500 ms.
		"eot_char 10",     // Set the EOT char
		"eot_enable 1",    // Append character when EOI detected?
		"savecfg 1",       // Enable saving of configuration parameters in EPROM
	}
	if clear {
		cmds = append(cmds, "clr")
	}
	for _, cmd := range cmds {
		if err := c.Command(cmd); err != nil {
			return nil, err
		}
	}
	return &c, nil
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
	cmd := strings.TrimSpace(s) + "\n"
	return c.Write([]byte(cmd))
}

// Query queries the instrument at the currently assigned GPIB using the Read
// and Write methods.
func (c *Controller) Query(s string) (string, error) {
	n, err := fmt.Fprintf(c.rw, "%s\n", s)
	if err != nil {
		return "", fmt.Errorf("error writing command: %s", err)
	}
	log.Printf("wrote %d bytes sending command `%s`", n, s)
	if !c.auto {
		readCmd := "++read"
		log.Printf("%s for query %s", readCmd, s)
		n, err := fmt.Fprintf(c.rw, "%s\n", readCmd)
		if err != nil {
			return "", fmt.Errorf("Whoa! %s", err)
		}
		log.Printf("wrote %d bytes sending command `%s`", n, readCmd)
	}
	return bufio.NewReader(c.rw).ReadString('\n')
}

// QueryCommand sends the given command to the Prologix controller and returns
// its response as a string. To indicate this is a command for the Prologix
// controller, thereby not transmitting over GPIB, two plus signs `++` are
// prepended. Addtionally, a new line is appended to act as the USB termination
// character.
func (c *Controller) QueryCommand(cmd string) (string, error) {
	_, err := fmt.Fprintf(c.rw, "++%s\n", strings.ToLower(cmd))
	if err != nil {
		return "", err
	}
	return bufio.NewReader(c.rw).ReadString('\n')
}

// Command sends the given command to the Prologix controller. To indicate this
// is a command for the Prologix controller, thereby not transmitting over
// GPIB, two plus signs `++` are prepended. Addtionally, a new line is appended
// to act as the USB termination character.
func (c *Controller) Command(cmd string) error {
	_, err := fmt.Fprintf(c.rw, "++%s\n", strings.ToLower(cmd))
	return err
}

// InstrumentAddress returns the GPIB address for the instrument under control.
func (c *Controller) InstrumentAddress() (int, error) {
	// FIXME(mdr): Need to update so that it can handle a secondary address.
	s, err := c.QueryCommand("addr")
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		return 0, err
	}
	addr := int(i)
	if addr != c.instrumentAddr {
		c.instrumentAddr = addr
		return addr, fmt.Errorf("internal state mismatch, address is now %d", addr)
	}
	return addr, nil
}

// SetInstrumentAddress sets the GPIB address for the instrument under control.
func (c *Controller) SetInstrumentAddress(addr int) error {
	cmd := fmt.Sprintf("addr %d", addr)
	err := c.Command(cmd)
	if err != nil {
		return err
	}
	c.instrumentAddr = addr
	return nil
}

// ReadAfterWrite determines if the Prologix controller is configured to
// automatically read after a write.
func (c *Controller) ReadAfterWrite() (bool, error) {
	s, err := c.QueryCommand("auto")
	if err != nil {
		return false, err
	}
	var auto bool
	if strings.TrimSpace(s) == "1" {
		auto = true
	} else if strings.TrimSpace(s) == "0" {
		auto = false
	} else {
		return false, fmt.Errorf("auto mode not determinable; received %s", s)
	}
	if auto != c.auto {
		c.auto = auto
		return false, fmt.Errorf("internal state mismatch, auto is now %t", auto)
	}
	return auto, nil
}

// SetReadAfterWrite sets the Proglogix controller to automatically read after write.
func (c *Controller) SetReadAfterWrite(enable bool) error {
	// Send the proper command based on whether enabling or disabling.
	cmd := "auto 0"
	if enable {
		cmd = "auto 1"
	}
	err := c.Command(cmd)
	// As long as there wasn't an error setting the auto mode on the Prologix
	// controller, set the auto status in the controller struct.
	if err != nil {
		return err
	}
	c.auto = enable
	return nil
}

// SendEOI determines if the Prologix controller is configured to append the
// EOI signal to the end of any command sent over the GPIB port.
func (c *Controller) SendEOI() (bool, error) {
	s, err := c.QueryCommand("eoi")
	if err != nil {
		return false, err
	}
	var eoi bool
	if strings.TrimSpace(s) == "1" {
		eoi = true
	} else if strings.TrimSpace(s) == "0" {
		eoi = false
	} else {
		return false, fmt.Errorf("eoi mode not determinable; received %s", s)
	}
	if eoi != c.eoi {
		c.eoi = eoi
		return false, fmt.Errorf("internal state mismatch, eoi is now %t", eoi)
	}
	return eoi, nil
}

// SetSendEOI sets the Proglogix controller to append the EOI signal after the
// last character of a command sent over the GPIB port.
func (c *Controller) SetSendEOI(enable bool) error {
	cmd := "eoi 0"
	if enable {
		cmd = "eoi 1"
	}
	err := c.Command(cmd)
	// As long as there wasn't an error setting the eoi mode on the Prologix
	// controller, set the eoi status in the controller struct.
	if err != nil {
		return err
	}
	c.eoi = enable
	return nil
}

// ReadTimeout queries the read timeout value in milliseconds from the Prologix
// GPIB controller.
func (c *Controller) ReadTimeout() (int, error) {
	s, err := c.QueryCommand("read_tmo_ms")
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		return 0, err
	}
	readTimeout := int(i)
	if readTimeout < 1 || readTimeout > 3000 {
		return 0, fmt.Errorf("read timeout must be between 1 and 3000 ms was set to %d", readTimeout)
	}
	return readTimeout, nil
}

// SetReadTimeout sets the Proglogix controllers read timeout in milliseconds.
// The timeout must be between 1 and 3000 milliseconds.
func (c *Controller) SetReadTimeout(timeout int) error {
	if timeout < 1 || timeout > 3000 {
		return fmt.Errorf("read timeout outside 1 to 3000 ms; attempted to set to %d", timeout)
	}
	return c.Command(fmt.Sprintf("read_tmo_ms %d", timeout))
}

// FrontPanel enables (local mode) or disables (local lockout) front panel
// operation for the instrument at the current GPIB address.
func (c *Controller) FrontPanel(enable bool) error {
	cmd := "llo" // Local lockout
	if enable {
		cmd = "loc" // Local
	}
	return c.Command(cmd)
}

// InterfaceClear asserts the GPIB IFC singal for 150 microseconds making the
// Prologix GPIB controller the Controller-In-Charge on the GPIB.
func (c *Controller) InterfaceClear() error {
	return c.Command("ifc")
}

// ServiceRequest determines if the GPIB SRQ signal is asserted or not.
func (c *Controller) ServiceRequest() (bool, error) {
	s, err := c.QueryCommand("srq")
	if err != nil {
		return false, err
	}
	var srq bool
	if strings.TrimSpace(s) == "1" {
		srq = true
	} else if strings.TrimSpace(s) == "0" {
		srq = false
	} else {
		return false, fmt.Errorf("srq not determinable; received %s", s)
	}
	return srq, nil
}

// Version returns the version string from the Prologix GPIB-USB controller.
func (c *Controller) Version() (string, error) {
	return c.QueryCommand("ver")
}

// Clear sends the Selected Device Clear (SDC) message to the currently
// selected GPIB address.
func (c *Controller) Clear() error {
	return c.Command("clr")
}
