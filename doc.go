// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package prologix communicates with either the Prologix GPIB-USB or
// GPIB-ETHERNET controller. Communication with the GPIB-USB controller can
// either be via the FTDI Virtual COM Port (VCP) driver or via the FTDI D2XX
// direct driver. The Prologix GPIB controller can operate as a GPIB
// Controller-in-Charge (CIC), a GPIB Talker Device, or a GPIB Listener Device.
//
// This package is part of the gotmc ecosystem. The visa package
// (github.com/gotmc/visa) defines a common interface for instrument
// communication across different transports (GPIB, USB, TCP/IP, serial). The
// asrl package provides the serial transport implementation. The ivi package
// (github.com/gotmc/ivi) builds on top of visa to provide standardized,
// instrument-class-specific APIs following the IVI Foundation specifications.
package prologix
