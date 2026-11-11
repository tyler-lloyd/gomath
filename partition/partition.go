// Package partition counts integer partitions.
//
// A partition of n is a way of writing n as a sum of positive integers,
// where the order of the terms does not matter. The number of partitions is
// written p(n): p(4) = 5, because 4 = 4 = 3+1 = 2+2 = 2+1+1 = 1+1+1+1.
package partition

import "math/big"

// Count returns p(n), the number of partitions of n. p(0) is 1.
//
// It is equivalent to Counts(n)[n]; see Counts for the method and cost.
func Count(n uint) *big.Int {
	return Counts(n)[n]
}

// Counts returns p(0), p(1), ..., p(n) as a slice of length n+1.
//
// It uses the recurrence that follows from Euler's pentagonal number theorem:
//
//	p(m) = sum over k >= 1 of (-1)^(k+1) * ( p(m - k(3k-1)/2) + p(m - k(3k+1)/2) )
//
// where p of a negative number is 0. The offsets k(3k-1)/2 and k(3k+1)/2 are
// the generalized pentagonal numbers 1, 2, 5, 7, 12, 15, ... Only about
// sqrt(m) of them are at most m, so the whole table takes on the order of
// n^1.5 additions, where enumerating the partitions themselves would take
// exponential time.
func Counts(n uint) []*big.Int {
	p := make([]*big.Int, n+1)
	p[0] = big.NewInt(1)

	for m := uint(1); m <= n; m++ {
		sum := new(big.Int)
		for k := uint(1); ; k++ {
			first := k * (3*k - 1) / 2
			if first > m {
				break
			}
			second := first + k // k(3k+1)/2

			if k%2 == 1 {
				sum.Add(sum, p[m-first])
				if second <= m {
					sum.Add(sum, p[m-second])
				}
			} else {
				sum.Sub(sum, p[m-first])
				if second <= m {
					sum.Sub(sum, p[m-second])
				}
			}
		}
		p[m] = sum
	}
	return p
}
