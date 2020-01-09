// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package prologix

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Controller models a GPIB controller-in-charge.
type Controller struct {
	rw             io.ReadWriter
	instrumentAddr int
	auto           bool
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
		"savecfg 0",    // Disable saving of configuration parameters in EPROM
		addrCmd,        // Set the primary address.
		"mode 1",       // Switch to controller mode.
		"auto 0",       // Turn off read-after-write and address instrument to listen.
		"eoi 1",        // Enable EOI assertion with last character.
		"eos 0",        // Set GPIB termination to append CR+LF to instrument commands.
		"eot_enable 0", // Do not append character when EOI detected
		"savecfg 1",    // Enable saving of configuration parameters in EPROM
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

// QueryCommand sends the given command to the Prologix controller and returns
// its response as a string. To indicate this is a command for the Prologix
// controller, thereby not transmitting over GPIB, two plus signs `++` are
// prepended. Addtionally, a new line is appended to act as the USB termination
// character.
func (c *Controller) QueryCommand(cmd string) (string, error) {
	_, err := fmt.Fprintf(c.rw, "++%s\n", cmd)
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
	_, err := fmt.Fprintf(c.rw, "++%s\n", cmd)
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
	c.instrumentAddr = addr
	return addr, nil
}

// SetAuto sets the Proglogix controller to automatically read after write.
func (c *Controller) SetAuto(enable bool) error {
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

// FrontPanel enables or disables front panel operation for the instrument at
// the current GPIB address.
func (c *Controller) FrontPanel(enable bool) error {
	cmd := "llo"
	if enable {
		cmd = "loc"
	}
	return c.Command(cmd)
}

// Version returns the version string of the Prologix GPUB-USB controller.
func (c *Controller) Version() (string, error) {
	return c.QueryCommand("ver")
}

// Clear sends the Selected Device Clear (SDC) message to the currently
// selected GPIB address.
func (c *Controller) Clear() error {
	return c.Command("clr")
}
