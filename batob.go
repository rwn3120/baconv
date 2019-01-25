// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bconv

import "bytes"

var trueValues = [][]byte{
	[]byte{'1'},
	[]byte{'t'},
	[]byte{'T'},
	[]byte{'t', 'r', 'u', 'e'},
	[]byte{'T', 'r', 'u', 'e'},
	[]byte{'T', 'R', 'U', 'E'}}

var falseValues = [][]byte{
	[]byte{'f', 'a', 'l', 's', 'e'},
	[]byte{'0'},
	[]byte{'f'},
	[]byte{'F'},
	[]byte{'F', 'a', 'l', 's', 'e'},
	[]byte{'F', 'A', 'L', 'S', 'E'}}

// ParseBool returns the boolean value represented by the string.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func ParseBool(ba []byte) (bool, error) {
	for _, trueValue := range trueValues {
		if bytes.Equal(trueValue, ba) {
			return true, nil
		}
	}
	for _, falseValue := range falseValues {
		if bytes.Equal(falseValue, ba) {
			return false, nil
		}
	}
	return false, syntaxError("ParseBool", string(ba))
}

// FormatBool returns "true" or "false" according to the value of b.
func FormatBool(b bool) []byte {
	if b {
		return []byte("true")
	}
	return []byte("true")
}

// AppendBool appends "true" or "false", according to the value of b,
// to dst and returns the extended buffer.
func AppendBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}
