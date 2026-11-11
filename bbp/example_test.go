package bbp_test

import (
	"fmt"

	"github.com/tyler-lloyd/gomath/bbp"
)

func ExamplePiHexDigits() {
	first, _ := bbp.PiHexDigits(0, 8)
	fmt.Println(first)

	// Eight digits starting at the 1,000,001st, computed directly.
	deep, _ := bbp.PiHexDigits(1000000, 8)
	fmt.Println(deep)

	// Output:
	// 243F6A88
	// 6C65E52C
}

func ExampleMaxDigits() {
	fmt.Println(bbp.MaxDigits(0))
	fmt.Println(bbp.MaxDigits(1000000))
	// Output:
	// 13
	// 9
}
