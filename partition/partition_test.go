package partition

import (
	"math/big"
	"testing"
)

func TestCountsSmallValues(t *testing.T) {
	// OEIS A000041.
	want := []int64{1, 1, 2, 3, 5, 7, 11, 15, 22, 30, 42, 56, 77, 101, 135, 176, 231, 297, 385, 490, 627}
	got := Counts(uint(len(want) - 1))
	if len(got) != len(want) {
		t.Fatalf("Counts returned %d values, want %d", len(got), len(want))
	}
	for n, w := range want {
		if got[n].Cmp(big.NewInt(w)) != 0 {
			t.Errorf("p(%d) = %s, want %d", n, got[n], w)
		}
	}
}

func TestCountLargeValues(t *testing.T) {
	tests := []struct {
		n    uint
		want string
	}{
		{0, "1"},
		{100, "190569292"},
		{200, "3972999029388"}, // MacMahon's hand computation, which Hardy and Ramanujan checked their formula against
		{1000, "24061467864032622473692149727991"},
	}
	for _, tt := range tests {
		if got := Count(tt.n).String(); got != tt.want {
			t.Errorf("Count(%d) = %s, want %s", tt.n, got, tt.want)
		}
	}
}

// bruteForce counts the partitions of n into parts of size at most max by
// direct recursion, independently of the pentagonal recurrence.
func bruteForce(n, max int) int64 {
	if n == 0 {
		return 1
	}
	var total int64
	for part := min(n, max); part >= 1; part-- {
		total += bruteForce(n-part, part)
	}
	return total
}

func TestCountsAgainstBruteForce(t *testing.T) {
	got := Counts(45)
	for n := 0; n <= 45; n++ {
		if want := bruteForce(n, n); got[n].Cmp(big.NewInt(want)) != 0 {
			t.Errorf("p(%d) = %s, want %d", n, got[n], want)
		}
	}
}

func TestRamanujanCongruences(t *testing.T) {
	// p(5k+4) = 0 (mod 5), p(7k+5) = 0 (mod 7), p(11k+6) = 0 (mod 11).
	p := Counts(2000)
	congruences := []struct{ modulus, offset uint }{{5, 4}, {7, 5}, {11, 6}}
	for _, c := range congruences {
		mod := big.NewInt(int64(c.modulus))
		for n := c.offset; n <= 2000; n += c.modulus {
			if new(big.Int).Mod(p[n], mod).Sign() != 0 {
				t.Errorf("p(%d) = %s is not divisible by %d", n, p[n], c.modulus)
			}
		}
	}
}
