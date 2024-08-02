// Copyright (c) 2020–2024 The prologix developers. All rights reserved.
// Project site: https://github.com/gotmc/prologix
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package prologix

import (
	"fmt"
	"testing"
)

func TestIsPrimaryAddressValid(t *testing.T) {
	tests := []struct {
		given int
		want  bool
	}{
		{-2, false},
		{-1, false},
		{0, true},
		{1, true},
		{15, true},
		{30, true},
		{31, false},
		{96, false},
		{126, false},
		{131, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Test: %d", test.given), func(t *testing.T) {
			if got := isPrimaryAddressValid(test.given); got != test.want {
				t.Errorf(
					"Error getting Date string\n\tgot %v; want %v",
					got,
					test.want,
				)
			}
		})
	}
}

func TestIsSecondaryAddressValid(t *testing.T) {
	tests := []struct {
		given int
		want  bool
	}{
		{-2, false},
		{-1, false},
		{0, false},
		{1, false},
		{15, false},
		{30, false},
		{31, false},
		{95, false},
		{96, true},
		{126, true},
		{127, false},
		{131, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Test: %d", test.given), func(t *testing.T) {
			if got := isSecondaryAddressValid(test.given); got != test.want {
				t.Errorf(
					"Error getting Date string\n\tgot %v; want %v",
					got,
					test.want,
				)
			}
		})
	}
}
