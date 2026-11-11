package proth

import (
	"math/big"
	"testing"
)

func TestIsPrimeAgreesWithProbablyPrime(t *testing.T) {
	for k := uint64(1); k < 300; k += 2 {
		for n := MinExponent(k); n <= 64; n++ {
			p := Number(k, n)
			if got, want := IsPrime(p), p.ProbablyPrime(20); got != want {
				t.Errorf("IsPrime(%d*2^%d+1 = %s) = %v, want %v", k, n, p, got, want)
			}
		}
	}
}

// OEIS A080076: the Proth primes in increasing order.
var prothPrimes = []int64{3, 5, 13, 17, 41, 97, 113, 193, 241, 257, 353, 449, 577, 641, 673, 769, 929, 1153}

func TestIsPrimeOnEveryIntegerBelow1200(t *testing.T) {
	want := map[int64]bool{}
	for _, p := range prothPrimes {
		want[p] = true
	}
	for i := int64(-5); i < 1200; i++ {
		if got := IsPrime(big.NewInt(i)); got != want[i] {
			t.Errorf("IsPrime(%d) = %v, want %v", i, got, want[i])
		}
	}
}

func TestIsNumber(t *testing.T) {
	// OEIS A080075: the Proth numbers in increasing order.
	want := map[int64]bool{}
	for _, p := range []int64{3, 5, 9, 13, 17, 25, 33, 41, 49, 57, 65, 81, 97, 113, 129, 145, 161, 177, 193} {
		want[p] = true
	}
	for i := int64(-5); i < 200; i++ {
		if got := IsNumber(big.NewInt(i)); got != want[i] {
			t.Errorf("IsNumber(%d) = %v, want %v", i, got, want[i])
		}
	}
}

func TestIsPrimeDoesNotModifyArgument(t *testing.T) {
	p := big.NewInt(113)
	IsPrime(p)
	IsNumber(p)
	if p.Cmp(big.NewInt(113)) != 0 {
		t.Errorf("argument was modified to %s", p)
	}
}

func TestFindPrime(t *testing.T) {
	tests := []struct {
		k     uint64
		wantN uint
	}{
		{1, 1},    // 3
		{3, 2},    // 13; 3*2+1 = 7 is prime but not a Proth number
		{7, 4},    // 113; 7*2^2+1 = 29 is prime but not a Proth number
		{99, 10},  // 101377
		{47, 583}, // needs far more than 200 bits of precision
	}
	for _, tt := range tests {
		n, p, ok := FindPrime(tt.k, 1000)
		if !ok || n != tt.wantN {
			t.Errorf("FindPrime(%d) = n %d, ok %v; want n %d", tt.k, n, ok, tt.wantN)
			continue
		}
		if !p.ProbablyPrime(20) {
			t.Errorf("FindPrime(%d) returned composite %s", tt.k, p)
		}
	}
}

func TestFindPrimeNotFound(t *testing.T) {
	// 78557 is a Sierpinski number: 78557*2^n + 1 is composite for every n.
	if n, _, ok := FindPrime(78557, 1000); ok {
		t.Errorf("FindPrime(78557) claimed a prime at n = %d", n)
	}
	if _, _, ok := FindPrime(47, 582); ok {
		t.Error("FindPrime(47, 582) should stop before the prime at n = 583")
	}
	for _, k := range []uint64{0, 2, 98} {
		if _, _, ok := FindPrime(k, 1000); ok {
			t.Errorf("FindPrime(%d) should reject an even k", k)
		}
	}
}
