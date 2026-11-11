package mersenne

import "testing"

// OEIS A000043: the exponents p for which 2^p - 1 is prime.
var primeExponents = []uint{2, 3, 5, 7, 13, 17, 19, 31, 61, 89, 107, 127, 521, 607, 1279, 2203, 2281, 3217, 4253, 4423}

func TestIsPrimeMatchesKnownExponents(t *testing.T) {
	want := map[uint]bool{}
	for _, p := range primeExponents {
		want[p] = true
	}
	limit := uint(1300)
	if testing.Short() {
		limit = 200
	}
	for p := uint(0); p <= limit; p++ {
		got, err := IsPrime(p)
		if err != nil {
			t.Fatalf("IsPrime(%d): %v", p, err)
		}
		if got != want[p] {
			t.Errorf("IsPrime(%d) = %v, want %v", p, got, want[p])
		}
	}
}

func TestIsPrimeLargerExponents(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	tests := []struct {
		p    uint
		want bool
	}{
		{2203, true}, {2281, true}, {3217, true}, {4253, true}, {4423, true},
		{2207, false}, {4421, false}, // prime exponents, composite Mersenne numbers
	}
	for _, tt := range tests {
		got, err := IsPrime(tt.p)
		if err != nil {
			t.Fatalf("IsPrime(%d): %v", tt.p, err)
		}
		if got != tt.want {
			t.Errorf("IsPrime(%d) = %v, want %v", tt.p, got, tt.want)
		}
	}
}

func TestIsPrimeAgreesWithProbablyPrime(t *testing.T) {
	for p := uint(0); p <= 130; p++ {
		m, _ := Number(p)
		got, _ := IsPrime(p)
		if want := m.ProbablyPrime(20); got != want {
			t.Errorf("IsPrime(%d) = %v, but ProbablyPrime(%s) = %v", p, got, m, want)
		}
	}
}

func TestNumber(t *testing.T) {
	for p, want := range map[uint]string{0: "0", 1: "1", 2: "3", 7: "127", 61: "2305843009213693951"} {
		m, err := Number(p)
		if err != nil {
			t.Fatalf("Number(%d): %v", p, err)
		}
		if m.String() != want {
			t.Errorf("Number(%d) = %s, want %s", p, m, want)
		}
	}
}

func TestExponentTooLarge(t *testing.T) {
	if _, err := Number(MaxExponent + 1); err == nil {
		t.Error("Number: expected an error past MaxExponent")
	}
	if _, err := IsPrime(MaxExponent + 1); err == nil {
		t.Error("IsPrime: expected an error past MaxExponent")
	}
}
