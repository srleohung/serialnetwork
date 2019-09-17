package encoding

import (
	"strconv"
)

func BCDEncode(x int) (y int) {
	s := strconv.Itoa(x)
	for i, r := range s {
		a, _ := strconv.Atoi(string(r))
		if l := (len(s) - 1 - i); l > 0 {
			y = y + a<<4*l
		} else {
			y = y + a
		}
	}
	return y
}
