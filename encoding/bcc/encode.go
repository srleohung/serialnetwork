/*
Block check character.
In telecommunications, a block check character (BCC) is a character added to a transmission block to facilitate error detection.
In longitudinal redundancy checking and cyclic redundancy checking, block check characters are computed for, and added to, each message block transmitted.
*/
package bcc

const BLOCK_CHECK_CHARACTER string = "bcc"

func Encode(bytes []byte) (sum byte) {
	for _, aByte := range bytes {
		sum ^= aByte
	}
	return sum
}
