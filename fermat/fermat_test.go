package fermat

import "testing"

func TestNumber(t *testing.T) {
	// OEIS A000215.
	want := []string{"3", "5", "17", "257", "65537", "4294967297", "18446744073709551617"}
	for n, w := range want {
		f, err := Number(uint(n))
		if err != nil {
			t.Fatalf("Number(%d): %v", n, err)
		}
		if f.String() != w {
			t.Errorf("Number(%d) = %s, want %s", n, f, w)
		}
	}
}

func TestIsPrime(t *testing.T) {
	// F(0)..F(4) are prime. F(5) = 641 * 6700417 (Euler, 1732), and every
	// Fermat number from F(5) to F(32) is known to be composite.
	limit := uint(13)
	if testing.Short() {
		limit = 10
	}
	for n := uint(0); n <= limit; n++ {
		got, err := IsPrime(n)
		if err != nil {
			t.Fatalf("IsPrime(%d): %v", n, err)
		}
		if want := n <= 4; got != want {
			t.Errorf("IsPrime(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestIsPrimeAgreesWithProbablyPrime(t *testing.T) {
	for n := uint(0); n <= 8; n++ {
		f, _ := Number(n)
		got, _ := IsPrime(n)
		if want := f.ProbablyPrime(20); got != want {
			t.Errorf("IsPrime(%d) = %v, but ProbablyPrime = %v", n, got, want)
		}
	}
}

func TestIndexTooLarge(t *testing.T) {
	if _, err := Number(MaxIndex + 1); err == nil {
		t.Error("Number: expected an error past MaxIndex")
	}
	if _, err := IsPrime(MaxIndex + 1); err == nil {
		t.Error("IsPrime: expected an error past MaxIndex")
	}
}
