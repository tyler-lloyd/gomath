package mersenne_test

import (
	"fmt"

	"github.com/tyler-lloyd/gomath/mersenne"
)

func ExampleIsPrime() {
	for _, p := range []uint{11, 13, 127} {
		prime, _ := mersenne.IsPrime(p)
		fmt.Printf("2^%d - 1 prime: %v\n", p, prime)
	}
	// Output:
	// 2^11 - 1 prime: false
	// 2^13 - 1 prime: true
	// 2^127 - 1 prime: true
}

func ExampleNumber() {
	m, _ := mersenne.Number(31)
	fmt.Println(m)
	// Output:
	// 2147483647
}
