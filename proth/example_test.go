package proth_test

import (
	"fmt"
	"math/big"

	"github.com/tyler-lloyd/gomath/proth"
)

func ExampleFindPrime() {
	n, prime, ok := proth.FindPrime(99, 1000)
	fmt.Println(n, prime, ok)
	// Output:
	// 10 101377 true
}

func ExampleIsPrime() {
	fmt.Println(proth.IsPrime(big.NewInt(113))) // 7*2^4 + 1
	fmt.Println(proth.IsPrime(big.NewInt(25)))  // 3*2^3 + 1, composite
	fmt.Println(proth.IsPrime(big.NewInt(7)))   // prime, but 3*2^1 + 1 has 2^1 < 3
	// Output:
	// true
	// false
	// false
}
