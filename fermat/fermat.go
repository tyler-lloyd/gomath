// Package fermat works with Fermat numbers and Fermat primes.
//
// The Fermat number F(n) is 2^(2^n) + 1. Fermat conjectured that all of them
// are prime; in fact F(0) through F(4) are the only primes known, and every
// other Fermat number whose status has been settled is composite.
package fermat

import (
	"fmt"
	"math/big"
)

// MaxIndex is the largest index the package accepts. Fermat numbers double
// in length with each step: F(MaxIndex) occupies 128 MiB. Testing indexes
// much beyond 15 is impractical with this straightforward implementation.
const MaxIndex = 30

var (
	one   = big.NewInt(1)
	three = big.NewInt(3)
)

// Number returns the Fermat number F(n) = 2^(2^n) + 1.
// It returns an error if n exceeds MaxIndex.
func Number(n uint) (*big.Int, error) {
	if n > MaxIndex {
		return nil, fmt.Errorf("fermat: index %d is too large: the maximum is %d", n, MaxIndex)
	}
	f := new(big.Int).Lsh(one, 1<<n)
	return f.Add(f, one), nil
}

// IsPrime reports whether the Fermat number F(n) = 2^(2^n) + 1 is prime.
// It returns an error if n exceeds MaxIndex.
//
// It uses Pepin's test: for n >= 1, F(n) is prime if and only if
//
//	3^((F(n)-1)/2) = -1  (mod F(n))
//
// This is Proth's theorem for the Proth number 1*2^(2^n) + 1, with the base
// fixed in advance: 3 is a quadratic non-residue modulo every Fermat prime
// beyond F(0). The answer is a proof, not a probabilistic result. Each step
// up in n makes the test roughly eight times slower.
func IsPrime(n uint) (bool, error) {
	f, err := Number(n)
	if err != nil {
		return false, err
	}
	if n == 0 {
		return true, nil // F(0) = 3, where the base 3 is not usable.
	}
	fMinusOne := new(big.Int).Sub(f, one)
	half := new(big.Int).Rsh(fMinusOne, 1)
	return new(big.Int).Exp(three, half, f).Cmp(fMinusOne) == 0, nil
}
