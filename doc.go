// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bconv implements conversions to and from string representations
// of basic data types.
//
// Numeric Conversions
//
// The most common numeric conversions are Batoi (string to int) and Itoba (int to string).
//
//	i, err := bconv.Batoi("-42")
//	s := bconv.Itoba(-42)
//
// These assume decimal and the Go int type.
//
// ParseBool, ParseFloat, ParseInt, and ParseUint convert strings to values:
//
//	b, err := bconv.ParseBool("true")
//	f, err := bconv.ParseFloat("3.1415", 64)
//	i, err := bconv.ParseInt("-42", 10, 64)
//	u, err := bconv.ParseUint("42", 10, 64)
//
// The parse functions return the widest type (float64, int64, and uint64),
// but if the size argument specifies a narrower width the result can be
// converted to that narrower type without data loss:
//
//	s := "2147483647" // biggest int32
//	i64, err := bconv.ParseInt(s, 10, 32)
//	...
//	i := int32(i64)
//
// FormatBool, FormatFloat, FormatInt, and FormatUint convert values to strings:
//
//	s := bconv.FormatBool(true)
//	s := bconv.FormatFloat(3.1415, 'E', -1, 64)
//	s := bconv.FormatInt(-42, 16)
//	s := bconv.FormatUint(42, 16)
//
// AppendBool, AppendFloat, AppendInt, and AppendUint are similar but
// append the formatted value to a destination slice.
//
// String Conversions
//
// Quote and QuoteToASCII convert strings to quoted Go string literals.
// The latter guarantees that the result is an ASCII string, by escaping
// any non-ASCII Unicode with \u:
//
//	q := bconv.Quote("Hello, 世界")
//	q := bconv.QuoteToASCII("Hello, 世界")
//
// QuoteRune and QuoteRuneToASCII are similar but accept runes and
// return quoted Go rune literals.
//
// Unquote and UnquoteChar unquote Go string and rune literals.
//
package bconv
