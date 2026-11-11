package fermat_test

import (
	"fmt"

	"github.com/tyler-lloyd/gomath/fermat"
)

func ExampleIsPrime() {
	for n := uint(3); n <= 6; n++ {
		f, _ := fermat.Number(n)
		prime, _ := fermat.IsPrime(n)
		fmt.Printf("F(%d) = %s prime: %v\n", n, f, prime)
	}
	// Output:
	// F(3) = 257 prime: true
	// F(4) = 65537 prime: true
	// F(5) = 4294967297 prime: false
	// F(6) = 18446744073709551617 prime: false
}
