// Copyright (c) 2020 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package prologix

// Driver provides the interface for either a Virtual COM Port (VCP), USB
// direct, or Ethernet connection.
type Driver interface {
	Command(cmd string) error
	QueryCommand(cmd string) (string, error)
}
