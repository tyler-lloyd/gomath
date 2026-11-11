// Package proth works with Proth numbers and Proth primes.
//
// A Proth number is an integer of the form k*2^n + 1 where k is odd and
// 2^n > k. A Proth prime is a Proth number that is prime.
package proth

import (
	"math/big"
	"math/bits"
)

var one = big.NewInt(1)

// MinExponent returns the smallest n for which k*2^n + 1 is a Proth number,
// i.e. the smallest n with 2^n > k.
func MinExponent(k uint64) uint {
	return uint(bits.Len64(k))
}

// Number returns k*2^n + 1. The result is a Proth number only if k is odd
// and n >= MinExponent(k).
func Number(k uint64, n uint) *big.Int {
	p := new(big.Int).SetUint64(k)
	p.Lsh(p, n)
	return p.Add(p, one)
}

// IsNumber reports whether p is a Proth number: p = k*2^n + 1 for some odd k
// and some n with 2^n > k.
func IsNumber(p *big.Int) bool {
	if p.Cmp(one) <= 0 {
		return false
	}
	// Write p-1 as k*2^n with k odd. Then 2^n > k exactly when k has at
	// most n bits.
	m := new(big.Int).Sub(p, one)
	n := m.TrailingZeroBits()
	k := m.Rsh(m, n)
	return uint(k.BitLen()) <= n
}

// IsPrime reports whether p is a Proth prime, that is, both a Proth number
// and prime. It returns false for primes that are not Proth numbers, such
// as 7.
//
// It applies Proth's theorem: a Proth number p is prime if and only if
// a^((p-1)/2) = -1 (mod p) for some integer a. If p is prime, Euler's
// criterion says the congruence holds for exactly those a that are quadratic
// non-residues mod p, so it is enough to find one a whose Jacobi symbol
// (a/p) is -1 and test that single base. Unlike a probabilistic test, the
// answer is a proof either way.
func IsPrime(p *big.Int) bool {
	if !IsNumber(p) {
		return false
	}

	// A perfect square has no a with Jacobi symbol -1, so the search below
	// would only end on reaching a factor of p. Rule squares out up front.
	if root := new(big.Int).Sqrt(p); root.Mul(root, root).Cmp(p) == 0 {
		return false
	}

	pMinusOne := new(big.Int).Sub(p, one)
	half := new(big.Int).Rsh(pMinusOne, 1)

	for a := big.NewInt(2); ; a.Add(a, one) {
		switch big.Jacobi(a, p) {
		case 0:
			// a shares a factor with p. Every prime has a non-residue
			// below itself, so we only get here when p is composite.
			return false
		case -1:
			return new(big.Int).Exp(a, half, p).Cmp(pMinusOne) == 0
		}
	}
}

// FindPrime searches for the smallest exponent n, with
// MinExponent(k) <= n <= maxExponent, for which k*2^n + 1 is a Proth prime.
// It returns that exponent and the prime, with ok set to true.
//
// ok is false if no such prime exists in the range, or if k is not odd.
// Some k, the Sierpinski numbers, never yield a prime; 78557 is the smallest
// known.
func FindPrime(k uint64, maxExponent uint) (n uint, prime *big.Int, ok bool) {
	if k%2 == 0 {
		return 0, nil, false
	}
	for n = MinExponent(k); n <= maxExponent; n++ {
		if p := Number(k, n); IsPrime(p) {
			return n, p, true
		}
	}
	return 0, nil, false
}
