/*
Binary coded decimal (BCD) is a system of writing numerals that assigns a four-digit binary code to each digit 0 through 9 in a decimal (base-10) numeral.
The four-bit BCD code for any particular single base-10 digit is its representation in binary notation, as follows: 0 = 0000. 1 = 0001. 2 = 0010.
*/
package bcd

import (
	"strconv"
)

func Encode(x int) (y int) {
	s := strconv.Itoa(x)
	for i, r := range s {
		_r, _ := strconv.Atoi(string(r))
		y = y + _r<<uint(4*(len(s)-1-i))
	}
	return y
}

func Decode(x int) (y int) {
	var _s int
	var _y string
	s := strconv.Itoa(x)
	for i, _ := range s {
		_x := (x - _s) >> uint(4*(len(s)-1-i))
		_s = _s + _x<<uint(4*(len(s)-1-i))
		_y = _y + strconv.Itoa(_x)
	}
	y, _ = strconv.Atoi(_y)
	return y
}
