// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package baconv

// decimal to binary floating point conversion.
// Algorithm:
//   1) Store input in multiprecision decimal.
//   2) Multiply/divide decimal by powers of two until in range [0.5, 1)
//   3) Multiply by 2^precision and round to get mantissa.

import (
	"bytes"
	"math"
)

var optimize = true // can change for testing

func equalIgnoreCase(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		if 'A' <= c1 && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		c2 := s2[i]
		if 'A' <= c2 && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}

func (d *decimal) set(ba []byte) (ok bool) {
	i := 0
	d.neg = false
	d.trunc = false

	// optional sign
	if i >= len(ba) {
		return
	}
	switch {
	case ba[i] == '+':
		i++
	case ba[i] == '-':
		d.neg = true
		i++
	}

	// digits
	sawdot := false
	sawdigits := false
	for ; i < len(ba); i++ {
		switch {
		case ba[i] == '.':
			if sawdot {
				return
			}
			sawdot = true
			d.dp = d.nd
			continue

		case '0' <= ba[i] && ba[i] <= '9':
			sawdigits = true
			if ba[i] == '0' && d.nd == 0 { // ignore leading zeros
				d.dp--
				continue
			}
			if d.nd < len(d.d) {
				d.d[d.nd] = ba[i]
				d.nd++
			} else if ba[i] != '0' {
				d.trunc = true
			}
			continue
		}
		break
	}
	if !sawdigits {
		return
	}
	if !sawdot {
		d.dp = d.nd
	}

	// optional exponent moves decimal point.
	// if we read a very large, very long number,
	// just be sure to move the decimal point by
	// a lot (say, 100000).  it doesn't matter if it's
	// not the exact number.
	if i < len(ba) && (ba[i] == 'e' || ba[i] == 'E') {
		i++
		if i >= len(ba) {
			return
		}
		esign := 1
		if ba[i] == '+' {
			i++
		} else if ba[i] == '-' {
			i++
			esign = -1
		}
		if i >= len(ba) || ba[i] < '0' || ba[i] > '9' {
			return
		}
		e := 0
		for ; i < len(ba) && '0' <= ba[i] && ba[i] <= '9'; i++ {
			if e < 10000 {
				e = e*10 + int(ba[i]) - '0'
			}
		}
		d.dp += e * esign
	}

	if i != len(ba) {
		return
	}

	ok = true
	return
}

// readFloat reads a decimal mantissa and exponent from a float
// string representation. It sets ok to false if the number could
// not fit return types or is invalid.
func readFloat(ba []byte) (mantissa uint64, exp int, neg, trunc, ok bool) {
	const uint64digits = 19
	i := 0

	// optional sign
	if i >= len(ba) {
		return
	}
	switch {
	case ba[i] == '+':
		i++
	case ba[i] == '-':
		neg = true
		i++
	}

	// digits
	sawdot := false
	sawdigits := false
	nd := 0
	ndMant := 0
	dp := 0
	for ; i < len(ba); i++ {
		switch c := ba[i]; true {
		case c == '.':
			if sawdot {
				return
			}
			sawdot = true
			dp = nd
			continue

		case '0' <= c && c <= '9':
			sawdigits = true
			if c == '0' && nd == 0 { // ignore leading zeros
				dp--
				continue
			}
			nd++
			if ndMant < uint64digits {
				mantissa *= 10
				mantissa += uint64(c - '0')
				ndMant++
			} else if ba[i] != '0' {
				trunc = true
			}
			continue
		}
		break
	}
	if !sawdigits {
		return
	}
	if !sawdot {
		dp = nd
	}

	// optional exponent moves decimal point.
	// if we read a very large, very long number,
	// just be sure to move the decimal point by
	// a lot (say, 100000).  it doesn't matter if it's
	// not the exact number.
	if i < len(ba) && (ba[i] == 'e' || ba[i] == 'E') {
		i++
		if i >= len(ba) {
			return
		}
		esign := 1
		if ba[i] == '+' {
			i++
		} else if ba[i] == '-' {
			i++
			esign = -1
		}
		if i >= len(ba) || ba[i] < '0' || ba[i] > '9' {
			return
		}
		e := 0
		for ; i < len(ba) && '0' <= ba[i] && ba[i] <= '9'; i++ {
			if e < 10000 {
				e = e*10 + int(ba[i]) - '0'
			}
		}
		dp += e * esign
	}

	if i != len(ba) {
		return
	}

	if mantissa != 0 {
		exp = dp - ndMant
	}
	ok = true
	return

}

// decimal power of ten to binary power of two.
var powtab = []int{1, 3, 6, 9, 13, 16, 19, 23, 26}

func (d *decimal) floatBits(flt *floatInfo) (b uint64, overflow bool) {
	var exp int
	var mant uint64

	// Zero is always a special case.
	if d.nd == 0 {
		mant = 0
		exp = flt.bias
		goto out
	}

	// Obvious overflow/underflow.
	// These bounds are for 64-bit floats.
	// Will have to change if we want to support 80-bit floats in the future.
	if d.dp > 310 {
		goto overflow
	}
	if d.dp < -330 {
		// zero
		mant = 0
		exp = flt.bias
		goto out
	}

	// Scale by powers of two until in range [0.5, 1.0)
	exp = 0
	for d.dp > 0 {
		var n int
		if d.dp >= len(powtab) {
			n = 27
		} else {
			n = powtab[d.dp]
		}
		d.Shift(-n)
		exp += n
	}
	for d.dp < 0 || d.dp == 0 && d.d[0] < '5' {
		var n int
		if -d.dp >= len(powtab) {
			n = 27
		} else {
			n = powtab[-d.dp]
		}
		d.Shift(n)
		exp -= n
	}

	// Our range is [0.5,1) but floating point range is [1,2).
	exp--

	// Minimum representable exponent is flt.bias+1.
	// If the exponent is smaller, move it up and
	// adjust d accordingly.
	if exp < flt.bias+1 {
		n := flt.bias + 1 - exp
		d.Shift(-n)
		exp += n
	}

	if exp-flt.bias >= 1<<flt.expbits-1 {
		goto overflow
	}

	// Extract 1+flt.mantbits bits.
	d.Shift(int(1 + flt.mantbits))
	mant = d.RoundedInteger()

	// Rounding might have added a bit; shift down.
	if mant == 2<<flt.mantbits {
		mant >>= 1
		exp++
		if exp-flt.bias >= 1<<flt.expbits-1 {
			goto overflow
		}
	}

	// Denormalized?
	if mant&(1<<flt.mantbits) == 0 {
		exp = flt.bias
	}
	goto out

overflow:
	// ±Inf
	mant = 0
	exp = 1<<flt.expbits - 1 + flt.bias
	overflow = true

out:
	// Assemble bits.
	bits := mant & (uint64(1)<<flt.mantbits - 1)
	bits |= uint64((exp-flt.bias)&(1<<flt.expbits-1)) << flt.mantbits
	if d.neg {
		bits |= 1 << flt.mantbits << flt.expbits
	}
	return bits, overflow
}

// Exact powers of 10.
var float64pow10 = []float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
var float32pow10 = []float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}

// If possible to convert decimal representation to 64-bit float f exactly,
// entirely in floating-point math, do so, avoiding the expense of decimalToFloatBits.
// Three common cases:
//	value is exact integer
//	value is exact integer * exact power of ten
//	value is exact integer / exact power of ten
// These all produce potentially inexact but correctly rounded answers.
func batof64exact(mantissa uint64, exp int, neg bool) (f float64, ok bool) {
	if mantissa>>float64info.mantbits != 0 {
		return
	}
	f = float64(mantissa)
	if neg {
		f = -f
	}
	switch {
	case exp == 0:
		// an integer.
		return f, true
	// Exact integers are <= 10^15.
	// Exact powers of ten are <= 10^22.
	case exp > 0 && exp <= 15+22: // int * 10^k
		// If exponent is big but number of digits is not,
		// can move a few zeros into the integer part.
		if exp > 22 {
			f *= float64pow10[exp-22]
			exp = 22
		}
		if f > 1e15 || f < -1e15 {
			// the exponent was really too large.
			return
		}
		return f * float64pow10[exp], true
	case exp < 0 && exp >= -22: // int / 10^k
		return f / float64pow10[-exp], true
	}
	return
}

// If possible to compute mantissa*10^exp to 32-bit float f exactly,
// entirely in floating-point math, do so, avoiding the machinery above.
func batof32exact(mantissa uint64, exp int, neg bool) (f float32, ok bool) {
	if mantissa>>float32info.mantbits != 0 {
		return
	}
	f = float32(mantissa)
	if neg {
		f = -f
	}
	switch {
	case exp == 0:
		return f, true
	// Exact integers are <= 10^7.
	// Exact powers of ten are <= 10^10.
	case exp > 0 && exp <= 7+10: // int * 10^k
		// If exponent is big but number of digits is not,
		// can move a few zeros into the integer part.
		if exp > 10 {
			f *= float32pow10[exp-10]
			exp = 10
		}
		if f > 1e7 || f < -1e7 {
			// the exponent was really too large.
			return
		}
		return f * float32pow10[exp], true
	case exp < 0 && exp >= -10: // int / 10^k
		return f / float32pow10[-exp], true
	}
	return
}

const fnParseFloat = "ParseFloat"

var infinity = [][]byte{
	[]byte{'i', 'n', 'f'},
	[]byte{'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'},
	[]byte{'+', 'i', 'n', 'f'},
	[]byte{'+', 'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'}}
var negInfinity = [][]byte{
	[]byte{'-', 'i', 'n', 'f'},
	[]byte{'-', 'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'}}
var nan = []byte{'n', 'a', 'n'}

func special(ba []byte) (f float64, ok bool) {
	if len(ba) == 0 {
		return
	}
	switch ba[0] {
	default:
		return
	case '+', 'i', 'I':
		for _, val := range infinity {
			if bytes.EqualFold(ba, val) {
				return math.Inf(1), true
			}
		}
	case '-':
		for _, val := range negInfinity {
			if bytes.EqualFold(ba, val) {
				return math.Inf(-1), true
			}
		}
	case 'n', 'N':
		if bytes.EqualFold(ba, nan) {
			return math.NaN(), true
		}
	}
	return
}

func Batof32(ba []byte) (f float32, err error) {
	if val, ok := special(ba); ok {
		return float32(val), nil
	}

	if optimize {
		// Parse mantissa and exponent.
		mantissa, exp, neg, trunc, ok := readFloat(ba)
		if ok {
			// Try pure floating-point arithmetic conversion.
			if !trunc {
				if f, ok := batof32exact(mantissa, exp, neg); ok {
					return f, nil
				}
			}
			// Try another fast path.
			ext := new(extFloat)
			if ok := ext.AssignDecimal(mantissa, exp, neg, trunc, &float32info); ok {
				bits, ovf := ext.floatBits(&float32info)
				f = math.Float32frombits(uint32(bits))
				if ovf {
					err = rangeError(fnParseFloat, string(ba))
				}
				return f, err
			}
		}
	}
	var d decimal
	if !d.set(ba) {
		return 0, syntaxError(fnParseFloat, string(ba))
	}
	bits, ovf := d.floatBits(&float32info)
	f = math.Float32frombits(uint32(bits))
	if ovf {
		err = rangeError(fnParseFloat, string(ba))
	}
	return f, err
}

func Batof64(ba []byte) (f float64, err error) {
	if val, ok := special(ba); ok {
		return val, nil
	}

	if optimize {
		// Parse mantissa and exponent.
		mantissa, exp, neg, trunc, ok := readFloat(ba)
		if ok {
			// Try pure floating-point arithmetic conversion.
			if !trunc {
				if f, ok := batof64exact(mantissa, exp, neg); ok {
					return f, nil
				}
			}
			// Try another fast path.
			ext := new(extFloat)
			if ok := ext.AssignDecimal(mantissa, exp, neg, trunc, &float64info); ok {
				bits, ovf := ext.floatBits(&float64info)
				f = math.Float64frombits(bits)
				if ovf {
					err = rangeError(fnParseFloat, string(ba))
				}
				return f, err
			}
		}
	}
	var d decimal
	if !d.set(ba) {
		return 0, syntaxError(fnParseFloat, string(ba))
	}
	bits, ovf := d.floatBits(&float64info)
	f = math.Float64frombits(bits)
	if ovf {
		err = rangeError(fnParseFloat, string(ba))
	}
	return f, err
}

// ParseFloat converts the string s to a floating-point number
// with the precision specified by bitSize: 32 for float32, or 64 for float64.
// When bitSize=32, the result still has type float64, but it will be
// convertible to float32 without changing its value.
//
// If s is well-formed and near a valid floating point number,
// ParseFloat returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
// The errors that ParseFloat returns have concrete type *NumError
// and include err.Num = s.
//
// If s is not syntactically well-formed, ParseFloat returns err.Err = ErrSyntax.
//
// If s is syntactically well-formed but is more than 1/2 ULP
// away from the largest floating point number of the given size,
// ParseFloat returns f = ±Inf, err.Err = ErrRange.
func ParseFloat(ba []byte, bitSize int) (float64, error) {
	if bitSize == 32 {
		f, err := Batof32(ba)
		return float64(f), err
	}
	return Batof64(ba)
}
