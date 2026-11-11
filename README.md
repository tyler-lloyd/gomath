# gomath

A Go library for computing difficult or obscure mathematics problems.

| Package | What it does |
| --- | --- |
| [`bbp`](bbp) | Hexadecimal digits of pi starting at any position, without computing the digits before it |
| [`proth`](proth) | Proth numbers and Proth primes `k*2^n + 1`, with a deterministic primality test |
| [`mersenne`](mersenne) | Mersenne primes `2^p - 1`, proved by the Lucas-Lehmer test |
| [`fermat`](fermat) | Fermat primes `2^(2^n) + 1`, proved by Pepin's test |
| [`partition`](partition) | The number of integer partitions `p(n)`, via Euler's pentagonal number theorem |

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

## Mersenne primes

A **Mersenne number** is `M(p) = 2^p - 1`. If `p` is composite, say `p = a*b`,
then `2^a - 1` divides `M(p)`, so only prime exponents can give primes. Most
still do not: `M(11) = 2047 = 23 * 89`.

For an odd prime `p`, the **Lucas-Lehmer test** settles the question. Define

```
s(0) = 4,    s(i+1) = s(i)^2 - 2  (mod M(p))
```

Then `M(p)` is prime if and only if `s(p-2) = 0`. The test needs only `p - 2`
squarings, and reducing modulo `2^p - 1` needs no division: because
`2^p = 1 (mod M(p))`, the high bits are simply added back onto the low ones.
This efficiency is why the largest known prime has almost always been a
Mersenne prime.

### Usage

```go
import "github.com/tyler-lloyd/gomath/mersenne"

prime, err := mersenne.IsPrime(127) // true:  2^127 - 1 is prime
prime, err = mersenne.IsPrime(11)   // false: 2047 = 23 * 89
m, err := mersenne.Number(31)       // 2147483647
```

| Function | Purpose |
| --- | --- |
| `IsPrime(p)` | Whether `2^p - 1` is prime, proved by the Lucas-Lehmer test |
| `Number(p)` | `2^p - 1` as a `*big.Int` |
| `MaxExponent` | The largest supported exponent, 2^27 |

The running time grows roughly as `p^3`. Exponents in the low thousands take
well under a second; the record-setting exponents in the tens of millions need
specialised FFT multiplication and are out of reach here.

## Fermat primes

A **Fermat number** is `F(n) = 2^(2^n) + 1`. The first five, 3, 5, 17, 257 and
65537, are prime, and Fermat conjectured that they all are. Euler disproved
this in 1732 by finding `F(5) = 641 * 6700417`. No further Fermat prime has
ever been found, and `F(5)` through `F(32)` are all known to be composite.

**Pepin's test** states that for `n >= 1`, `F(n)` is prime if and only if

```
3^((F(n)-1)/2) = -1  (mod F(n))
```

This is Proth's theorem applied to the Proth number `1*2^(2^n) + 1`, with the
base known in advance: 3 is a quadratic non-residue modulo every Fermat prime
beyond `F(0)`.

### Usage

```go
import "github.com/tyler-lloyd/gomath/fermat"

prime, err := fermat.IsPrime(4) // true:  65537
prime, err = fermat.IsPrime(5)  // false: 4294967297 = 641 * 6700417
f, err := fermat.Number(6)      // 18446744073709551617
```

| Function | Purpose |
| --- | --- |
| `IsPrime(n)` | Whether `2^(2^n) + 1` is prime, proved by Pepin's test |
| `Number(n)` | `2^(2^n) + 1` as a `*big.Int` |
| `MaxIndex` | The largest supported index, 30 |

Fermat numbers double in length at every step, so each increase in `n` makes
the test roughly eight times slower. Indexes up to about 13 finish within a
second or so.

## Integer partitions

A **partition** of `n` is a way of writing it as a sum of positive integers,
ignoring order. There are five partitions of 4:

```
4 = 3+1 = 2+2 = 2+1+1 = 1+1+1+1
```

so `p(4) = 5`. The partition function grows very quickly:
`p(100) = 190,569,292` and `p(1000)` has 32 digits.

Listing partitions to count them is hopeless, but **Euler's pentagonal number
theorem** gives a recurrence that uses only earlier values:

```
p(n) = p(n-1) + p(n-2) - p(n-5) - p(n-7) + p(n-12) + p(n-15) - ...
```

The offsets 1, 2, 5, 7, 12, 15, ... are the generalized pentagonal numbers
`k(3k-1)/2` and `k(3k+1)/2`, and the signs alternate in pairs. Only about
`sqrt(n)` terms are non-zero, so a table up to `n` costs on the order of
`n^1.5` additions.

The values also obey **Ramanujan's congruences**, which the tests verify:
`p(5k+4)` is divisible by 5, `p(7k+5)` by 7, and `p(11k+6)` by 11.

### Usage

```go
import "github.com/tyler-lloyd/gomath/partition"

partition.Count(100) // 190569292
partition.Counts(10) // [1 1 2 3 5 7 11 15 22 30 42]
```

| Function | Purpose |
| --- | --- |
| `Count(n)` | `p(n)` as a `*big.Int` |
| `Counts(n)` | The table `p(0), ..., p(n)` |

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
- [Lucas-Lehmer primality test](https://en.wikipedia.org/wiki/Lucas%E2%80%93Lehmer_primality_test)
  and [Pepin's test](https://en.wikipedia.org/wiki/P%C3%A9pin%27s_test) on Wikipedia.
- [Pentagonal number theorem](https://en.wikipedia.org/wiki/Pentagonal_number_theorem)
  on Wikipedia; OEIS [A000043](https://oeis.org/A000043) (Mersenne exponents)
  and [A000041](https://oeis.org/A000041) (partition numbers).
