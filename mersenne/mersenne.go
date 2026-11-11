// Package mersenne works with Mersenne numbers and Mersenne primes.
//
// The Mersenne number M(p) is 2^p - 1. It can only be prime when p itself is
// prime, and the Lucas-Lehmer test decides the remaining cases exactly.
package mersenne

import (
	"fmt"
	"math/big"
)

// MaxExponent is the largest exponent the package accepts. M(MaxExponent)
// occupies 16 MiB; testing exponents anywhere near this size is far beyond
// what this straightforward implementation can finish in reasonable time.
const MaxExponent = 1 << 27

var (
	one = big.NewInt(1)
	two = big.NewInt(2)
)

// Number returns the Mersenne number M(p) = 2^p - 1.
// It returns an error if p exceeds MaxExponent.
func Number(p uint) (*big.Int, error) {
	if p > MaxExponent {
		return nil, fmt.Errorf("mersenne: exponent %d is too large: the maximum is %d", p, MaxExponent)
	}
	m := new(big.Int).Lsh(one, p)
	return m.Sub(m, one), nil
}

// IsPrime reports whether the Mersenne number M(p) = 2^p - 1 is prime.
// It returns an error if p exceeds MaxExponent.
//
// If p is composite, say p = a*b, then 2^a - 1 divides 2^p - 1, so M(p) is
// composite. For an odd prime p the Lucas-Lehmer test applies: define
//
//	s(0) = 4,  s(i+1) = s(i)^2 - 2  (mod M(p))
//
// Then M(p) is prime if and only if s(p-2) = 0. The answer is a proof, not a
// probabilistic result. The running time grows roughly as p^3.
func IsPrime(p uint) (bool, error) {
	m, err := Number(p)
	if err != nil {
		return false, err
	}
	if p == 2 {
		return true, nil // M(2) = 3; Lucas-Lehmer needs an odd exponent.
	}
	if !isPrimeExponent(p) {
		return false, nil
	}

	s := big.NewInt(4)
	hi := new(big.Int)
	for i := uint(0); i < p-2; i++ {
		s.Mul(s, s)
		s.Sub(s, two)
		if s.Sign() < 0 {
			s.Add(s, m)
		}

		// Reduce modulo 2^p - 1 without dividing. Since 2^p = 1 (mod M(p)),
		// the bits above position p can be folded back onto the low p bits.
		// s < M(p)^2 here, so one fold and one subtraction are enough.
		hi.Rsh(s, p)
		s.And(s, m)
		s.Add(s, hi)
		if s.Cmp(m) >= 0 {
			s.Sub(s, m)
		}
	}
	return s.Sign() == 0, nil
}

// isPrimeExponent reports whether p is prime, by trial division.
// p is at most MaxExponent, so the loop runs at most a few thousand times.
func isPrimeExponent(p uint) bool {
	if p < 2 {
		return false
	}
	if p%2 == 0 {
		return p == 2
	}
	for d := uint(3); d*d <= p; d += 2 {
		if p%d == 0 {
			return false
		}
	}
	return true
}
