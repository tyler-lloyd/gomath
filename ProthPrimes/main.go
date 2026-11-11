// Command ProthPrimes searches for the smallest Proth prime k*2^n + 1 for a
// given odd k.
package main

import (
	"flag"
	"fmt"
	"math/big"
	"math/bits"
	"os"
	"strconv"
)

const defaultMaxExponent = 1000

func main() {
	maxExponent := flag.Uint("max", defaultMaxExponent, "largest exponent n to try")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [-max n] k\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Finds the smallest Proth prime k*2^n + 1 for the odd, positive integer k.")
		fmt.Fprintln(flag.CommandLine.Output())
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	k, err := strconv.ParseUint(flag.Arg(0), 10, 64)
	if err != nil || k%2 == 0 {
		fmt.Fprintf(os.Stderr, "invalid k %q: must be an odd, positive integer\n", flag.Arg(0))
		os.Exit(2)
	}

	n, prime, found := FindProthPrime(k, *maxExponent)
	if !found {
		fmt.Printf("No Proth prime %d*2^n + 1 found for %d <= n <= %d.\n", k, MinExponent(k), *maxExponent)
		os.Exit(1)
	}
	fmt.Println("Proth prime has been found!")
	fmt.Printf("%d*2^%d + 1 = %s\n", k, n, prime)
}

// MinExponent returns the smallest n for which k*2^n + 1 is a Proth number,
// i.e. the smallest n with 2^n > k.
func MinExponent(k uint64) uint {
	return uint(bits.Len64(k))
}

// ProthNumber returns k*2^n + 1.
func ProthNumber(k uint64, n uint) *big.Int {
	p := new(big.Int).SetUint64(k)
	p.Lsh(p, n)
	return p.Add(p, big.NewInt(1))
}

// FindProthPrime returns the smallest exponent n <= maxExponent, and the
// corresponding prime, for which k*2^n + 1 is a Proth prime. k must be odd.
func FindProthPrime(k uint64, maxExponent uint) (n uint, prime *big.Int, found bool) {
	for n = MinExponent(k); n <= maxExponent; n++ {
		if p := ProthNumber(k, n); IsProthPrime(p) {
			return n, p, true
		}
	}
	return 0, nil, false
}

// IsProthPrime reports whether p is prime, given that p is a Proth number:
// p = k*2^n + 1 with k odd and 2^n > k.
//
// It applies Proth's theorem: such a p is prime if and only if
// a^((p-1)/2) = -1 (mod p) for some integer a. If p is prime, Euler's
// criterion says the congruence holds for exactly those a that are quadratic
// non-residues mod p, so it is enough to find one a whose Jacobi symbol
// (a/p) is -1 and test that single base. Unlike a probabilistic test, the
// answer is a proof either way.
func IsProthPrime(p *big.Int) bool {
	// A perfect square has no a with Jacobi symbol -1, so the search below
	// would only end on reaching a factor of p. Rule squares out up front.
	if root := new(big.Int).Sqrt(p); root.Mul(root, root).Cmp(p) == 0 {
		return false
	}

	pMinusOne := new(big.Int).Sub(p, big.NewInt(1))
	half := new(big.Int).Rsh(pMinusOne, 1)

	a := big.NewInt(2)
	for ; ; a.Add(a, big.NewInt(1)) {
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
