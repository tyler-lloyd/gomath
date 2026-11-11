# gomath

Small Go programs for computing difficult or obscure mathematics problems.

| Program | What it does |
| --- | --- |
| [BBP](BBP) | Prints hexadecimal digits of pi starting at any position, without computing the digits before it |
| [ProthPrimes](ProthPrimes) | Finds the smallest Proth prime `k*2^n + 1` for a given odd `k` |

Requires Go 1.21 or newer. There are no dependencies outside the standard library.

```bash
go test ./...
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
maximum position of 536,870,911. The program refuses to print more than that.

### Usage

```
BBP [-digits n] [position]
```

`position` is the number of hex digits to skip after the point (default
1,000,000), so `0` gives the very first digits of `pi = 3.243F6A88...`

```
$ go run ./BBP 0
Hex digits of pi starting at position 1: 243F6A88
$ go run ./BBP
Hex digits of pi starting at position 1000001: 6C65E52C
$ go run ./BBP -digits 13 0
Hex digits of pi starting at position 1: 243F6A8885A30
```

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
non-residue mod `N` (Euler's criterion). So the program picks the smallest `a`
whose Jacobi symbol `(a/N)` is `-1` and performs one modular exponentiation.
The result is a proof of primality or compositeness, not a probabilistic
answer. Several of the largest known primes are Proth primes found this way.

For a given odd `k`, the program tries each `n` from the smallest allowed value
up to a limit (1000 by default) and reports the first Proth prime.

Some `k` never produce one. These are the **Sierpinski numbers**; the smallest
known is 78557, for which `78557*2^n + 1` is always divisible by one of
3, 5, 7, 13, 19, 37 or 73.

### Usage

```
ProthPrimes [-max n] k
```

```
$ go run ./ProthPrimes 99
Proth prime has been found!
99*2^10 + 1 = 101377
$ go run ./ProthPrimes 47
Proth prime has been found!
47*2^583 + 1 = 1487939695262196876907983166454197495251350196192890428923003345454869706240895712896623468784438158657419591298913094265537812046389415279164757669092989298186306341246574002177
$ go run ./ProthPrimes 78557
No Proth prime 78557*2^n + 1 found for 17 <= n <= 1000.
```

The exit status is 0 when a prime is found, 1 when none is found within the
limit, and 2 for invalid input.

## References

- D. H. Bailey, P. B. Borwein, S. Plouffe, "On the Rapid Computation of Various
  Polylogarithmic Constants", *Mathematics of Computation* 66 (1997).
- D. H. Bailey, "The BBP Algorithm for Pi" (2006).
- [Proth's theorem](https://en.wikipedia.org/wiki/Proth%27s_theorem) and
  [Proth prime](https://en.wikipedia.org/wiki/Proth_prime) on Wikipedia.
