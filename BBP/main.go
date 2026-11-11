// Command BBP prints hexadecimal digits of pi starting at an arbitrary
// position, without computing any of the digits before it, using the
// Bailey-Borwein-Plouffe digit-extraction algorithm.
package main

import (
	"flag"
	"fmt"
	"math/bits"
	"os"
	"strconv"
)

const (
	// maxPosition keeps every denominator 8k+j below 2^32, so that the
	// products in powMod16 cannot overflow a uint64.
	maxPosition = 1<<29 - 1

	defaultPosition = 1000000
	defaultDigits   = 8

	// tailTerms is the number of terms summed past the requested position.
	// Term m of the tail is smaller than 16^-m, so it stops contributing to
	// a 64-bit fraction after 16 terms.
	tailTerms = 16
)

func main() {
	digits := flag.Int("digits", 0, "number of hex digits to print (default 8, or fewer if the position is too large to guarantee 8)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [-digits n] [position]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Prints the hex digits of pi that follow the given position after the\nhexadecimal point (default %d). Position 0 gives the first digits.\n\n", defaultPosition)
		flag.PrintDefaults()
	}
	flag.Parse()

	position := uint64(defaultPosition)
	if flag.NArg() > 0 {
		var err error
		position, err = strconv.ParseUint(flag.Arg(0), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid position %q: must be a non-negative integer\n", flag.Arg(0))
			os.Exit(2)
		}
	}
	if position > maxPosition {
		fmt.Fprintf(os.Stderr, "position %d is too large: the maximum is %d\n", position, maxPosition)
		os.Exit(2)
	}

	count := *digits
	if count == 0 {
		count = min(defaultDigits, MaxDigits(position))
	}

	hex, err := PiHexDigits(position, count)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Printf("Hex digits of pi starting at position %d: %s\n", position+1, hex)
}

// PiHexDigits returns count hexadecimal digits of pi, beginning with the
// digit at position+1 after the hexadecimal point. PiHexDigits(0, 8) is
// "243F6A88".
func PiHexDigits(position uint64, count int) (string, error) {
	if position > maxPosition {
		return "", fmt.Errorf("position %d is too large: the maximum is %d", position, maxPosition)
	}
	if limit := MaxDigits(position); count < 1 || count > limit {
		return "", fmt.Errorf("digit count must be between 1 and %d at position %d", limit, position)
	}

	// The BBP formula is
	//
	//	pi = sum over k >= 0 of 16^-k * (4/(8k+1) - 2/(8k+4) - 1/(8k+5) - 1/(8k+6))
	//
	// so the fractional part of 16^position * pi is the fractional part of
	// 4*S(1) - 2*S(4) - S(5) - S(6). Each S is a 64-bit binary fraction, and
	// unsigned overflow discards the integer part for us.
	frac := 4*series(1, position) - 2*series(4, position) - series(5, position) - series(6, position)

	return fmt.Sprintf("%016X", frac)[:count], nil
}

// MaxDigits reports how many hex digits PiHexDigits can be trusted to produce
// at the given position.
//
// Every term of a series is truncated to 64 bits, an error of less than
// 2^-64. The four series hold position+tailTerms terms each and are combined
// with weights 4, 2, 1 and 1, so the total error is below
// 8*(position+tailTerms)*2^-64. One further digit is held back as a guard.
func MaxDigits(position uint64) int {
	errorBits := bits.Len64(8 * (position + tailTerms))
	return (64-errorBits)/4 - 1
}

// series returns the fractional part of
//
//	S(j) = sum over k >= 0 of 16^(n-k) / (8k+j)
//
// as a 64-bit binary fraction, i.e. floor(frac(S(j)) * 2^64) up to
// truncation error.
func series(j, n uint64) uint64 {
	var sum uint64

	// For k < n the power of 16 is an integer. Only its remainder modulo the
	// denominator affects the fractional part, which is what makes it
	// possible to skip straight to digit n.
	for k := uint64(0); k < n; k++ {
		d := 8*k + j
		sum += fraction(powMod16(n-k, d), d)
	}

	// For k >= n the terms are 16^-(k-n) / (8k+j) and shrink rapidly.
	for m := uint64(0); m < tailTerms; m++ {
		d := 8*(n+m) + j
		if m == 0 {
			// 1/d. When d is 1 the term is an integer and contributes nothing.
			if d > 1 {
				sum += fraction(1, d)
			}
			continue
		}
		sum += (1 << (64 - 4*m)) / d
	}

	return sum
}

// fraction returns floor(r/d * 2^64). It requires r < d.
func fraction(r, d uint64) uint64 {
	q, _ := bits.Div64(r, 0, d)
	return q
}

// powMod16 returns 16^exp mod m by left-to-right binary exponentiation.
// It requires 0 < m < 2^32.
func powMod16(exp, m uint64) uint64 {
	if m == 1 {
		return 0
	}
	result := uint64(1)
	for bit := bits.Len64(exp) - 1; bit >= 0; bit-- {
		result = result * result % m
		if exp>>uint(bit)&1 == 1 {
			result = result * 16 % m
		}
	}
	return result
}
