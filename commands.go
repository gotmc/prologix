// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package prologix

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"go.uber.org/multierr"
)

// AssertEOI determines if the Prologix controller is configured to assert the
// EOI signal at the end of any command sent over the GPIB port.
func (c *Controller) AssertEOI() (bool, error) {
	s, err := c.QueryController("eoi")
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

// ClearDevice sends the `clr` command to the Prologix controller which sends
// the Selected Device Clear (SDC) message to the currently selected GPIB
// address.
func (c *Controller) ClearDevice() error {
	return c.CommandController("clr")
}

// ClearInterface sends the `ifc` command to the Prologix controller which
// asserts the GPIB Interface Clear (IFC) signal for 150 microseconds making
// the Prologix GPIB controller the Controller-In-Charge.
func (c *Controller) ClearInterface() error {
	return c.CommandController("ifc")
}

// FrontPanel enables or disables front panel operation for the instrument at
// the current GPIB address. This is accomplished by either sending the
// Prologix `loc` local command to enable the front panel or by sending the
// Prologix `llo` local lockout command to disable the front panel.
func (c *Controller) FrontPanel(enable bool) error {
	cmd := "llo" // Local lockout
	if enable {
		cmd = "loc" // Local
	}
	return c.CommandController(cmd)
}

// GPIBTermination uses the Prologix `eos` command to query the GPIB
// terminator.
func (c *Controller) GPIBTermination() (GpibTerm, error) {
	s, err := c.QueryController("eos")
	if err != nil {
		return 0, err
	}
	term, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, err
	}
	return GpibTerm(term), nil
}

// InstrumentAddress returns the GPIB address for the instrument under control.
func (c *Controller) InstrumentAddress() (int, int, error) {
	s, err := c.QueryController("addr")
	if err != nil {
		return 0, 0, err
	}
	s = strings.TrimSpace(s)
	idx := 0
	for i, c := range s {
		if !unicode.IsNumber(c) {
			idx = i
			break
		}
	}
	if idx == 0 {
		idx = len(s)
	}
	pri, err := strconv.ParseInt(s[:idx], 10, 8)
	if err != nil {
		return 0, 0, err
	}

	var errs []error

	if int(pri) != c.primaryAddr {
		errs = append(errs,
			fmt.Errorf("internal state mismatch, primary address was %d now %d", c.primaryAddr, pri))
	}
	c.primaryAddr = int(pri)

	// secondary address?
	remain := strings.TrimLeft(s[idx:], " :,\t")
	var sec int64
	if len(remain) > 0 {
		sec, err = strconv.ParseInt(remain, 10, 8)
		if err != nil {
			errs = append(errs, fmt.Errorf("parsing secondary address in %s: %w", s, err))
			return 0, 0, multierr.Combine(errs...)
		}
	}
	if c.secondaryAddr != int(sec) {
		errs = append(errs,
			fmt.Errorf("internal state mismatch, secondary address was %d now %d", c.secondaryAddr, sec))
	}
	c.secondaryAddr = int(sec)

	merr := multierr.Combine(errs...) // result is nil if errs is empty
	return c.primaryAddr, c.secondaryAddr, merr
}

// ReadAfterWrite determines if the Prologix controller is configured to
// automatically read after a write.
func (c *Controller) ReadAfterWrite() (bool, error) {
	s, err := c.QueryController("auto")
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

// ReadTimeout queries the read timeout value in milliseconds from the Prologix
// GPIB controller.
func (c *Controller) ReadTimeout() (int, error) {
	s, err := c.QueryController("read_tmo_ms")
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		return 0, err
	}
	readTimeout := int(i)
	if readTimeout < 1 || readTimeout > 3000 {
		return 0, fmt.Errorf(
			"read timeout must be between 1 and 3000 ms was set to %d",
			readTimeout,
		)
	}
	return readTimeout, nil
}

// Reset performs a power-on reset of the controller. The process takes about 5
// seconds. All input received during this time are ignored.
func (c *Controller) Reset() error {
	return c.CommandController("rst")
}

// ServiceRequest sends the `srq` command to the Prologix controller to
// determine if the GPIB SRQ signal is asserted or not.
func (c *Controller) ServiceRequest() (bool, error) {
	s, err := c.QueryController("srq")
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

// SetAssertEOI sets the Prologix controller to assert the EOI signal after the
// last character of a command sent over the GPIB port.
func (c *Controller) SetAssertEOI(enable bool) error {
	cmd := "eoi 0"
	if enable {
		cmd = "eoi 1"
	}
	err := c.CommandController(cmd)
	// As long as there wasn't an error setting the eoi mode on the Prologix
	// controller, set the eoi status in the controller struct.
	if err != nil {
		return err
	}
	c.eoi = enable
	return nil
}

// SetGPIBTermination uses the Prologix `eos` command to set the character to
// be appended as the GPIB terminator to all data sent from the Prologix
// Controller to the instrument.
func (c *Controller) SetGPIBTermination(term GpibTerm) error {
	return c.CommandController(fmt.Sprintf("eos %d", term))
}

// SetInstrumentAddress sets the GPIB address for the instrument under control.
func (c *Controller) SetInstrumentAddress(addr int) error {
	cmd := fmt.Sprintf("addr %d", addr)
	err := c.CommandController(cmd)
	if err != nil {
		return err
	}
	c.primaryAddr = addr
	return nil
}

// SetReadAfterWrite sets the Proglogix controller to automatically read after write.
func (c *Controller) SetReadAfterWrite(enable bool) error {
	// Send the proper command based on whether enabling or disabling.
	cmd := "auto 0"
	if enable {
		cmd = "auto 1"
	}
	err := c.CommandController(cmd)
	// As long as there wasn't an error setting the auto mode on the Prologix
	// controller, set the auto status in the controller struct.
	if err != nil {
		return err
	}
	c.auto = enable
	return nil
}

// SetReadTimeout sets the Proglogix controller's read timeout in milliseconds.
// The timeout must be between 1 and 3000 milliseconds.
func (c *Controller) SetReadTimeout(timeout int) error {
	if timeout < 1 || timeout > 3000 {
		return fmt.Errorf("read timeout outside 1 to 3000 ms; attempted to set to %d", timeout)
	}
	return c.CommandController(fmt.Sprintf("read_tmo_ms %d", timeout))
}

// Version returns the version string from the Prologix GPIB controller.
func (c *Controller) Version() (string, error) {
	return c.QueryController("ver")
}
