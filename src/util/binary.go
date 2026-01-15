package util

import "math/bits"

func NextPowerOfTwo(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	if (n & (n - 1)) == 0 {
		return n
	}
	bitLen := bits.Len64(n)
	if bitLen >= 64 {
		return 1 << 63
	}
	return 1 << bitLen
}
