// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

//go:build noasm || (!amd64 && !arm64)
// +build noasm !amd64,!arm64

package p751

import (
	"math/bits"

	"github.com/JI-0/circl/dh/sidh/internal/common"
)

// Compute z = x + y (mod p).
func addP751(z, x, y *common.Fp) {
	var carry uint64

	// z=x+y % P751
	for i := 0; i < FpWords; i++ {
		z[i], carry = bits.Add64(x[i], y[i], carry)
	}

	// z = z - P751x2
	carry = 0
	for i := 0; i < FpWords; i++ {
		z[i], carry = bits.Sub64(z[i], P751x2[i], carry)
	}

	// if z<0 add P751x2 back
	mask := uint64(0 - carry)
	carry = 0
	for i := 0; i < FpWords; i++ {
		z[i], carry = bits.Add64(z[i], P751x2[i]&mask, carry)
	}
}

// Compute z = x - y (mod p).
func subP751(z, x, y *common.Fp) {
	var borrow uint64

	for i := 0; i < FpWords; i++ {
		z[i], borrow = bits.Sub64(x[i], y[i], borrow)
	}

	mask := uint64(0 - borrow)
	borrow = 0

	for i := 0; i < FpWords; i++ {
		z[i], borrow = bits.Add64(z[i], P751x2[i]&mask, borrow)
	}
}

// If choice = 0, leave x unchanged. If choice = 1, sets x to y.
// If choice is neither 0 nor 1 then behaviour is undefined.
// This function executes in constant time.
func cmovP751(x, y *common.Fp, choice uint8) {
	mask := 0 - uint64(choice)
	for i := 0; i < FpWords; i++ {
		x[i] ^= mask & (x[i] ^ y[i])
	}
}

// Conditionally swaps bits in x and y in constant time.
// mask indicates bits to be swapped (set bits are swapped)
// For details see "Hackers Delight, 2.20"
//
// Implementation doesn't actually depend on a prime field.
func cswapP751(x, y *common.Fp, mask uint8) {
	var tmp, mask64 uint64

	mask64 = 0 - uint64(mask)
	for i := 0; i < FpWords; i++ {
		tmp = mask64 & (x[i] ^ y[i])
		x[i] = tmp ^ x[i]
		y[i] = tmp ^ y[i]
	}
}

// Perform Montgomery reduction: set z = x R^{-1} (mod 2*p)
// with R=2^(FpWords*64). Destroys the input value.
func rdcP751(z *common.Fp, x *common.FpX2) {
	var carry, t, u, v uint64
	var hi, lo uint64
	var count int

	count = P751p1Zeros

	for i := 0; i < FpWords; i++ {
		for j := 0; j < i; j++ {
			if j < (i - count + 1) {
				hi, lo = bits.Mul64(z[j], P751p1[i-j])
				v, carry = bits.Add64(lo, v, 0)
				u, carry = bits.Add64(hi, u, carry)
				t += carry
			}
		}
		v, carry = bits.Add64(v, x[i], 0)
		u, carry = bits.Add64(u, 0, carry)
		t += carry

		z[i] = v
		v = u
		u = t
		t = 0
	}

	for i := FpWords; i < 2*FpWords-1; i++ {
		if count > 0 {
			count--
		}
		for j := i - FpWords + 1; j < FpWords; j++ {
			if j < (FpWords - count) {
				hi, lo = bits.Mul64(z[j], P751p1[i-j])
				v, carry = bits.Add64(lo, v, 0)
				u, carry = bits.Add64(hi, u, carry)
				t += carry
			}
		}
		v, carry = bits.Add64(v, x[i], 0)
		u, carry = bits.Add64(u, 0, carry)

		t += carry
		z[i-FpWords] = v
		v = u
		u = t
		t = 0
	}
	v, _ = bits.Add64(v, x[2*FpWords-1], 0)
	z[FpWords-1] = v
}

// Compute z = x * y.
func mulP751(z *common.FpX2, x, y *common.Fp) {
	var u, v, t uint64
	var hi, lo uint64
	var carry uint64

	for i := uint64(0); i < FpWords; i++ {
		for j := uint64(0); j <= i; j++ {
			hi, lo = bits.Mul64(x[j], y[i-j])
			v, carry = bits.Add64(lo, v, 0)
			u, carry = bits.Add64(hi, u, carry)
			t += carry
		}
		z[i] = v
		v = u
		u = t
		t = 0
	}

	for i := FpWords; i < (2*FpWords)-1; i++ {
		for j := i - FpWords + 1; j < FpWords; j++ {
			hi, lo = bits.Mul64(x[j], y[i-j])
			v, carry = bits.Add64(lo, v, 0)
			u, carry = bits.Add64(hi, u, carry)
			t += carry
		}
		z[i] = v
		v = u
		u = t
		t = 0
	}
	z[2*FpWords-1] = v
}

// Compute z = x + y, without reducing mod p.
func adlP751(z, x, y *common.FpX2) {
	var carry uint64
	for i := 0; i < 2*FpWords; i++ {
		z[i], carry = bits.Add64(x[i], y[i], carry)
	}
}

// Reduce a field element in [0, 2*p) to one in [0,p).
func modP751(x *common.Fp) {
	var borrow, mask uint64
	for i := 0; i < FpWords; i++ {
		x[i], borrow = bits.Sub64(x[i], P751[i], borrow)
	}

	// Sets all bits if borrow = 1
	mask = 0 - borrow
	borrow = 0
	for i := 0; i < FpWords; i++ {
		x[i], borrow = bits.Add64(x[i], P751[i]&mask, borrow)
	}
}

// Compute z = x - y, without reducing mod p.
func sulP751(z, x, y *common.FpX2) {
	var borrow, mask uint64
	for i := 0; i < 2*FpWords; i++ {
		z[i], borrow = bits.Sub64(x[i], y[i], borrow)
	}

	// Sets all bits if borrow = 1
	mask = 0 - borrow
	borrow = 0
	for i := FpWords; i < 2*FpWords; i++ {
		z[i], borrow = bits.Add64(z[i], P751[i-FpWords]&mask, borrow)
	}
}
