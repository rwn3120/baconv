// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package baconv

import (
	"fmt"
	"log"
)

func ExampleAppendBool() {
	b := []byte("bool:")
	b = AppendBool(b, true)
	fmt.Println(string(b))

	// Output:
	// bool:true
}

func ExampleAppendFloat() {
	b32 := []byte("float32:")
	b32 = AppendFloat(b32, 3.1415926535, 'E', -1, 32)
	fmt.Println(string(b32))

	b64 := []byte("float64:")
	b64 = AppendFloat(b64, 3.1415926535, 'E', -1, 64)
	fmt.Println(string(b64))

	// Output:
	// float32:3.1415927E+00
	// float64:3.1415926535E+00
}

func ExampleAppendInt() {
	b10 := []byte("int (base 10):")
	b10 = AppendInt(b10, -42, 10)
	fmt.Println(string(b10))

	b16 := []byte("int (base 16):")
	b16 = AppendInt(b16, -42, 16)
	fmt.Println(string(b16))

	// Output:
	// int (base 10):-42
	// int (base 16):-2a
}

func ExampleAppendQuote() {
	b := []byte("quote:")
	b = AppendQuote(b, `"Fran & Freddie's Diner"`)
	fmt.Println(string(b))

	// Output:
	// quote:"\"Fran & Freddie's Diner\""
}

func ExampleAppendQuoteRune() {
	b := []byte("rune:")
	b = AppendQuoteRune(b, '☺')
	fmt.Println(string(b))

	// Output:
	// rune:'☺'
}

func ExampleAppendQuoteRuneToASCII() {
	b := []byte("rune (ascii):")
	b = AppendQuoteRuneToASCII(b, '☺')
	fmt.Println(string(b))

	// Output:
	// rune (ascii):'\u263a'
}

func ExampleAppendQuoteToASCII() {
	b := []byte("quote (ascii):")
	b = AppendQuoteToASCII(b, `"Fran & Freddie's Diner"`)
	fmt.Println(string(b))

	// Output:
	// quote (ascii):"\"Fran & Freddie's Diner\""
}

func ExampleAppendUint() {
	b10 := []byte("uint (base 10):")
	b10 = AppendUint(b10, 42, 10)
	fmt.Println(string(b10))

	b16 := []byte("uint (base 16):")
	b16 = AppendUint(b16, 42, 16)
	fmt.Println(string(b16))

	// Output:
	// uint (base 10):42
	// uint (base 16):2a
}

func ExampleBatoi() {
	v := "10"
	if s, err := Batoi([]byte(v)); err == nil {
		fmt.Printf("%T, %v", s, s)
	}

	// Output:
	// int, 10
}

func ExampleCanBackquote() {
	fmt.Println(CanBackquote("Fran & Freddie's Diner ☺"))
	fmt.Println(CanBackquote("`can't backquote this`"))

	// Output:
	// true
	// false
}

func ExamplestringFormatBool() {
	v := true
	s := string(FormatBool(v))
	fmt.Printf("%T, %v\n", s, s)

	// Output:
	// string, true
}

func ExampleFormatFloat() {
	v := 3.1415926535

	s32 := string(FormatFloat(v, 'E', -1, 32))
	fmt.Printf("%T, %v\n", s32, s32)

	s64 := string(FormatFloat(v, 'E', -1, 64))
	fmt.Printf("%T, %v\n", s64, s64)

	// Output:
	// string, 3.1415927E+00
	// string, 3.1415926535E+00
}

func ExampleFormatInt() {
	v := int64(-42)

	s10 := FormatInt(v, 10)
	fmt.Printf("%T, %v\n", s10, s10)

	s16 := FormatInt(v, 16)
	fmt.Printf("%T, %v\n", s16, s16)

	// Output:
	// string, -42
	// string, -2a
}

func ExampleFormatUint() {
	v := uint64(42)

	s10 := FormatUint(v, 10)
	fmt.Printf("%T, %v\n", s10, s10)

	s16 := FormatUint(v, 16)
	fmt.Printf("%T, %v\n", s16, s16)

	// Output:
	// string, 42
	// string, 2a
}

func ExampleIsGraphic() {
	shamrock := IsGraphic('☘')
	fmt.Println(shamrock)

	a := IsGraphic('a')
	fmt.Println(a)

	bel := IsGraphic('\007')
	fmt.Println(bel)

	// Output:
	// true
	// true
	// false
}

func ExampleIsPrint() {
	c := IsPrint('\u263a')
	fmt.Println(c)

	bel := IsPrint('\007')
	fmt.Println(bel)

	// Output:
	// true
	// false
}

func ExampleItoba() {
	i := 10
	s := Itoba(i)
	fmt.Printf("%T, %v\n", s, s)

	// Output:
	// string, 10
}

func ExampleParseBool() {
	v := "true"
	if s, err := ParseBool([]byte(v)); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	// Output:
	// bool, true
}

func ExampleParseFloat() {
	v := "3.1415926535"
	if s, err := ParseFloat([]byte(v), 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := ParseFloat([]byte(v), 64); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	// Output:
	// float64, 3.1415927410125732
	// float64, 3.1415926535
}

func ExampleParseInt() {
	v32 := "-354634382"
	if s, err := ParseInt([]byte(v32), 10, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := ParseInt([]byte(v32), 16, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	v64 := "-3546343826724305832"
	if s, err := ParseInt([]byte(v64), 10, 64); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := ParseInt([]byte(v64), 16, 64); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	// Output:
	// int64, -354634382
	// int64, -3546343826724305832
}

func ExampleParseUint() {
	v := "42"
	if s, err := ParseUint([]byte(v), 10, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := ParseUint([]byte(v), 10, 64); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	// Output:
	// uint64, 42
	// uint64, 42
}

func ExampleQuote() {
	s := Quote(`"Fran & Freddie's Diner	☺"`) // there is a tab character inside the string literal
	fmt.Println(s)

	// Output:
	// "\"Fran & Freddie's Diner\t☺\""
}

func ExampleQuoteRune() {
	s := QuoteRune('☺')
	fmt.Println(s)

	// Output:
	// '☺'
}

func ExampleQuoteRuneToASCII() {
	s := QuoteRuneToASCII('☺')
	fmt.Println(s)

	// Output:
	// '\u263a'
}

func ExampleQuoteRuneToGraphic() {
	s := QuoteRuneToGraphic('☺')
	fmt.Println(s)

	s = QuoteRuneToGraphic('\u263a')
	fmt.Println(s)

	s = QuoteRuneToGraphic('\u000a')
	fmt.Println(s)

	s = QuoteRuneToGraphic('	') // tab character
	fmt.Println(s)

	// Output:
	// '☺'
	// '☺'
	// '\n'
	// '\t'
}

func ExampleQuoteToASCII() {
	s := QuoteToASCII(`"Fran & Freddie's Diner	☺"`) // there is a tab character inside the string literal
	fmt.Println(s)

	// Output:
	// "\"Fran & Freddie's Diner\t\u263a\""
}

func ExampleQuoteToGraphic() {
	s := QuoteToGraphic("☺")
	fmt.Println(s)

	s = QuoteToGraphic("This is a \u263a	\u000a") // there is a tab character inside the string literal
	fmt.Println(s)

	s = QuoteToGraphic(`" This is a ☺ \n "`)
	fmt.Println(s)

	// Output:
	// "☺"
	// "This is a ☺\t\n"
	// "\" This is a ☺ \\n \""
}

func ExampleUnquote() {
	s, err := Unquote("You can't unquote a string without quotes")
	fmt.Printf("%q, %v\n", s, err)
	s, err = Unquote("\"The string must be either double-quoted\"")
	fmt.Printf("%q, %v\n", s, err)
	s, err = Unquote("`or backquoted.`")
	fmt.Printf("%q, %v\n", s, err)
	s, err = Unquote("'\u263a'") // single character only allowed in single quotes
	fmt.Printf("%q, %v\n", s, err)
	s, err = Unquote("'\u2639\u2639'")
	fmt.Printf("%q, %v\n", s, err)

	// Output:
	// "", invalid syntax
	// "The string must be either double-quoted", <nil>
	// "or backquoted.", <nil>
	// "☺", <nil>
	// "", invalid syntax
}

func ExampleUnquoteChar() {
	v, mb, t, err := UnquoteChar(`\"Fran & Freddie's Diner\"`, '"')
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("value:", string(v))
	fmt.Println("multibyte:", mb)
	fmt.Println("tail:", t)

	// Output:
	// value: "
	// multibyte: false
	// tail: Fran & Freddie's Diner\"
}

func ExampleNumError() {
	str := "Not a number"
	if _, err := ParseFloat([]byte(str), 64); err != nil {
		e := err.(*NumError)
		fmt.Println("Func:", e.Func)
		fmt.Println("Num:", e.Num)
		fmt.Println("Err:", e.Err)
		fmt.Println(err)
	}

	// Output:
	// Func: ParseFloat
	// Num: Not a number
	// Err: invalid syntax
	// bconv.ParseFloat: parsing "Not a number": invalid syntax
}
