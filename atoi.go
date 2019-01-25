// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bconv

import "errors"

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat)
	Num  string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *NumError) Error() string {
	return "bconv." + e.Func + ": " + "parsing " + Quote(e.Num) + ": " + e.Err.Error()
}

func syntaxError(fn, str string) *NumError {
	return &NumError{fn, str, ErrSyntax}
}

func rangeError(fn, str string) *NumError {
	return &NumError{fn, str, ErrRange}
}

func baseError(fn, str string, base int) *NumError {
	return &NumError{fn, str, errors.New("invalid base " + Itoba(base))}
}

func bitSizeError(fn, str string, bitSize int) *NumError {
	return &NumError{fn, str, errors.New("invalid bit size " + Itoba(bitSize))}
}

const intSize = 32 << (^uint(0) >> 63)

// IntSize is the size in bits of an int or uint value.
const IntSize = intSize

const maxUint64 = 1<<64 - 1

// ParseUint is like ParseInt but for unsigned numbers.
func ParseUint(ba []byte, base int, bitSize int) (uint64, error) {
	const fnParseUint = "ParseUint"

	if len(ba) == 0 {
		return 0, syntaxError(fnParseUint, string(ba))
	}

	ba0 := ba
	switch {
	case 2 <= base && base <= 36:
		// valid base; nothing to do

	case base == 0:
		// Look for octal, hex prefix.
		switch {
		case ba[0] == '0' && len(ba) > 1 && (ba[1] == 'x' || ba[1] == 'X'):
			if len(ba) < 3 {
				return 0, syntaxError(fnParseUint, string(ba0))
			}
			base = 16
			ba = ba[2:]
		case ba[0] == '0':
			base = 8
			ba = ba[1:]
		default:
			base = 10
		}

	default:
		return 0, baseError(fnParseUint, string(ba0), base)
	}

	if bitSize == 0 {
		bitSize = int(IntSize)
	} else if bitSize < 0 || bitSize > 64 {
		return 0, bitSizeError(fnParseUint, string(ba0), bitSize)
	}

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	// Use compile-time constants for common cases.
	var cutoff uint64
	switch base {
	case 10:
		cutoff = maxUint64/10 + 1
	case 16:
		cutoff = maxUint64/16 + 1
	default:
		cutoff = maxUint64/uint64(base) + 1
	}

	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
	for _, c := range []byte(ba) {
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= c && c <= 'z':
			d = c - 'a' + 10
		case 'A' <= c && c <= 'Z':
			d = c - 'A' + 10
		default:
			return 0, syntaxError(fnParseUint, string(ba0))
		}

		if d >= byte(base) {
			return 0, syntaxError(fnParseUint, string(ba0))
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, rangeError(fnParseUint, string(ba0))
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return maxVal, rangeError(fnParseUint, string(ba0))
		}
		n = n1
	}

	return n, nil
}

// ParseInt interprets a string s in the given base (0, 2 to 36) and
// bit size (0 to 64) and returns the corresponding value i.
//
// If base == 0, the base is implied by the string's prefix:
// base 16 for "0x", base 8 for "0", and base 10 otherwise.
// For bases 1, below 0 or above 36 an error is returned.
//
// The bitSize argument specifies the integer type
// that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
// correspond to int, int8, int16, int32, and int64.
// For a bitSize below 0 or above 64 an error is returned.
//
// The errors that ParseInt returns have concrete type *NumError
// and include err.Num = s. If s is empty or contains invalid
// digits, err.Err = ErrSyntax and the returned value is 0;
// if the value corresponding to s cannot be represented by a
// signed integer of the given size, err.Err = ErrRange and the
// returned value is the maximum magnitude integer of the
// appropriate bitSize and sign.
func ParseInt(ba []byte, base int, bitSize int) (i int64, err error) {
	const fnParseInt = "ParseInt"

	// Empty string bad.
	if len(ba) == 0 {
		return 0, syntaxError(fnParseInt, string(ba))
	}

	// Pick off leading sign.
	ba0 := ba
	neg := false
	if ba[0] == '+' {
		ba = ba[1:]
	} else if ba[0] == '-' {
		neg = true
		ba = ba[1:]
	}

	// Convert unsigned and check range.
	var un uint64
	un, err = ParseUint(ba, base, bitSize)
	if err != nil && err.(*NumError).Err != ErrRange {
		err.(*NumError).Func = fnParseInt
		err.(*NumError).Num = string(ba0)
		return 0, err
	}

	if bitSize == 0 {
		bitSize = int(IntSize)
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), rangeError(fnParseInt, string(ba0))
	}
	if neg && un > cutoff {
		return -int64(cutoff), rangeError(fnParseInt, string(ba0))
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

// Batoi is equivalent to ParseInt(s, 10, 0), converted to type int.
func Batoi(ba []byte) (int, error) {
	const fnBatoi = "Batoi"

	sLen := len(ba)
	if intSize == 32 && (0 < sLen && sLen < 10) ||
		intSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.
		ba0 := ba
		if ba[0] == '-' || ba[0] == '+' {
			ba = ba[1:]
			if len(ba) < 1 {
				return 0, &NumError{fnBatoi, string(ba0), ErrSyntax}
			}
		}

		n := 0
		for _, ch := range []byte(ba) {
			ch -= '0'
			if ch > 9 {
				return 0, &NumError{fnBatoi, string(ba0), ErrSyntax}
			}
			n = n*10 + int(ch)
		}
		if ba0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Slow path for invalid or big integers.
	i64, err := ParseInt(ba, 10, 0)
	if nerr, ok := err.(*NumError); ok {
		nerr.Func = fnBatoi
	}
	return int(i64), err
}
