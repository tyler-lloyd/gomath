# gomath

A Go library for computing difficult or obscure mathematics problems.

| Package | What it does |
| --- | --- |
| [`bbp`](bbp) | Hexadecimal digits of pi starting at any position, without computing the digits before it |
| [`proth`](proth) | Proth numbers and Proth primes `k*2^n + 1`, with a deterministic primality test |

Requires Go 1.21 or newer. There are no dependencies outside the standard library.

```bash
go get github.com/tyler-lloyd/gomath
```

## Bailey-Borwein-Plouffe algorithm

The Bailey-Borwein-Plouffe (BBP) formula, discovered in 1995, is

```
pi = sum over k >= 0 of  16^-k * ( 4/(8k+1) - 2/(8k+4) - 1/(8k+5) - 1/(8k+6) )
```

Because of the `16^-k` factor, it can be turned into a *digit-extraction*
algorithm: one that produces the hexadecimal digits of pi at position `n`
directly, in roughly `O(n log n)` time and constant memory.

Multiplying by `16^n` moves the wanted digits to just after the hexadecimal
point, so they are the fractional part of `4*S(1) - 2*S(4) - S(5) - S(6)`, where

```
S(j) = sum over k >= 0 of  16^(n-k) / (8k+j)
```

- For `k < n`, `16^(n-k)` is an integer, and only `16^(n-k) mod (8k+j)` affects
  the fractional part. That remainder is cheap to compute with modular
  exponentiation, so the huge power is never formed.
- For `k >= n`, the terms shrink by a factor of 16 each time, and a handful of
  them is enough.

This implementation keeps each sum as a 64-bit fixed-point fraction, so the
integer parts fall away through unsigned overflow and there is no
floating-point rounding. Each term is truncated by less than `2^-64`, which
bounds the total error and therefore the number of digits that can be trusted
at a given position: 13 at position 0, 9 at position one million, 6 at the
maximum position of 536,870,911. `PiHexDigits` refuses to return more than that.

### Usage

```go
import "github.com/tyler-lloyd/gomath/bbp"

// pi = 3.243F6A88... in hexadecimal.
digits, err := bbp.PiHexDigits(0, 8) // "243F6A88"

// Skip the first million digits and compute the next eight directly.
digits, err = bbp.PiHexDigits(1000000, 8) // "6C65E52C"
```

| Function | Purpose |
| --- | --- |
| `PiHexDigits(position, count)` | `count` hex digits of pi, starting with the digit at `position+1` after the point |
| `MaxDigits(position)` | How many digits can be guaranteed at `position` |
| `MaxPosition` | The largest supported position, 536,870,911 |

`PiHexDigits` returns an error if the position is past `MaxPosition` or if
more digits are requested than `MaxDigits` allows.

## Proth primes

A **Proth number** is an integer of the form `N = k*2^n + 1` where `k` is odd
and `2^n > k`. A **Proth prime** is a Proth number that is prime. The first few
are 3, 5, 13, 17, 41, 97, 113, 193, 241, 257
([OEIS A080076](https://oeis.org/A080076)).

The condition `2^n > k` matters. Without it every odd number would qualify,
and it is what makes Proth numbers easy to test. **Proth's theorem** states
that a Proth number `N` is prime if and only if there is an integer `a` with

```
a^((N-1)/2) = -1  (mod N)
```

If `N` is prime, the congruence holds exactly when `a` is a quadratic
non-residue mod `N` (Euler's criterion). So `IsPrime` picks the smallest `a`
whose Jacobi symbol `(a/N)` is `-1` and performs one modular exponentiation.
The result is a proof of primality or compositeness, not a probabilistic
answer. Several of the largest known primes are Proth primes found this way.

Some `k` never produce a Proth prime. These are the **Sierpinski numbers**; the
smallest known is 78557, for which `78557*2^n + 1` is always divisible by one
of 3, 5, 7, 13, 19, 37 or 73.

### Usage

```go
import "github.com/tyler-lloyd/gomath/proth"

// Smallest Proth prime 99*2^n + 1 with n <= 1000.
n, prime, ok := proth.FindPrime(99, 1000) // 10, 101377, true

proth.IsPrime(big.NewInt(113)) // true:  7*2^4 + 1
proth.IsPrime(big.NewInt(25))  // false: a Proth number, but composite
proth.IsPrime(big.NewInt(7))   // false: prime, but not a Proth number
```

| Function | Purpose |
| --- | --- |
| `FindPrime(k, maxExponent)` | Smallest `n <= maxExponent` for which `k*2^n + 1` is a Proth prime |
| `IsPrime(p)` | Whether `p` is a Proth prime, proved by Proth's theorem |
| `IsNumber(p)` | Whether `p` is a Proth number |
| `Number(k, n)` | `k*2^n + 1` as a `*big.Int` |
| `MinExponent(k)` | Smallest `n` with `2^n > k` |

`FindPrime` reports `ok == false` when no prime exists in the range or when
`k` is even. Because the search needs arbitrary precision, primes such as
`47*2^583 + 1` (177 digits) are found correctly.

## Testing

```bash
go test ./...
```

## References

- D. H. Bailey, P. B. Borwein, S. Plouffe, "On the Rapid Computation of Various
  Polylogarithmic Constants", *Mathematics of Computation* 66 (1997).
- D. H. Bailey, "The BBP Algorithm for Pi" (2006).
- [Proth's theorem](https://en.wikipedia.org/wiki/Proth%27s_theorem) and
  [Proth prime](https://en.wikipedia.org/wiki/Proth_prime) on Wikipedia.
