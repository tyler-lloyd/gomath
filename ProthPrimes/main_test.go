package main

import (
	"math/big"
	"testing"
)

func TestIsProthPrimeAgreesWithProbablyPrime(t *testing.T) {
	for k := uint64(1); k < 300; k += 2 {
		for n := MinExponent(k); n <= 64; n++ {
			p := ProthNumber(k, n)
			if got, want := IsProthPrime(p), p.ProbablyPrime(20); got != want {
				t.Errorf("IsProthPrime(%d*2^%d+1 = %s) = %v, want %v", k, n, p, got, want)
			}
		}
	}
}

func TestProthPrimeSequence(t *testing.T) {
	// OEIS A080076: the Proth primes in increasing order.
	want := []int64{3, 5, 13, 17, 41, 97, 113, 193, 241, 257, 353, 449, 577, 641, 673, 769, 929, 1153}

	seen := map[int64]bool{}
	for k := uint64(1); k < 1200; k += 2 {
		for n := MinExponent(k); n <= 12; n++ {
			if p := ProthNumber(k, n); p.Cmp(big.NewInt(1200)) < 0 && IsProthPrime(p) {
				seen[p.Int64()] = true
			}
		}
	}
	if len(seen) != len(want) {
		t.Errorf("found %d Proth primes below 1200, want %d: %v", len(seen), len(want), seen)
	}
	for _, p := range want {
		if !seen[p] {
			t.Errorf("Proth prime %d was not found", p)
		}
	}
}

func TestFindProthPrime(t *testing.T) {
	tests := []struct {
		k     uint64
		wantN uint
	}{
		{1, 1},    // 3
		{3, 2},    // 13; 3*2+1 = 7 is prime but not a Proth number
		{7, 4},    // 113; 7*2^2+1 = 29 is prime but not a Proth number
		{47, 583}, // needs far more than 200 bits of precision
	}
	for _, tt := range tests {
		n, p, found := FindProthPrime(tt.k, defaultMaxExponent)
		if !found || n != tt.wantN {
			t.Errorf("FindProthPrime(%d) = n %d, found %v; want n %d", tt.k, n, found, tt.wantN)
			continue
		}
		if !p.ProbablyPrime(20) {
			t.Errorf("FindProthPrime(%d) returned composite %s", tt.k, p)
		}
	}
}

func TestSierpinskiNumberHasNoProthPrime(t *testing.T) {
	// 78557 is a Sierpinski number: 78557*2^n + 1 is composite for every n.
	if n, _, found := FindProthPrime(78557, defaultMaxExponent); found {
		t.Errorf("FindProthPrime(78557) claimed a prime at n = %d", n)
	}
}
