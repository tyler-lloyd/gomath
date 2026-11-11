package partition_test

import (
	"fmt"

	"github.com/tyler-lloyd/gomath/partition"
)

func ExampleCount() {
	fmt.Println(partition.Count(4))
	fmt.Println(partition.Count(100))
	fmt.Println(partition.Count(1000))
	// Output:
	// 5
	// 190569292
	// 24061467864032622473692149727991
}

func ExampleCounts() {
	fmt.Println(partition.Counts(10))
	// Output:
	// [1 1 2 3 5 7 11 15 22 30 42]
}
