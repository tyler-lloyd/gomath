package bbp

import (
	"math/big"
	"testing"
)

// The first 64 hex digits of the fractional part of pi.
const piHex = "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89"

func TestPiHexDigitsMatchesKnownExpansion(t *testing.T) {
	for position := 0; position+12 <= len(piHex); position++ {
		got, err := PiHexDigits(uint64(position), 12)
		if err != nil {
			t.Fatalf("position %d: %v", position, err)
		}
		if want := piHex[position : position+12]; got != want {
			t.Errorf("position %d: got %s, want %s", position, got, want)
		}
	}
}

func TestPiHexDigitsDeepPositions(t *testing.T) {
	// Published values, e.g. in Bailey's "The BBP Algorithm for Pi".
	tests := []struct {
		position uint64
		want     string
	}{
		{1000000, "6C65E52CB"},
		{10000000, "7AF58"},
	}
	for _, tt := range tests {
		if testing.Short() && tt.position > 1000000 {
			continue
		}
		got, err := PiHexDigits(tt.position, len(tt.want))
		if err != nil {
			t.Fatalf("position %d: %v", tt.position, err)
		}
		if got != tt.want {
			t.Errorf("position %d: got %s, want %s", tt.position, got, tt.want)
		}
	}
}

func TestPiHexDigitsRejectsBadInput(t *testing.T) {
	if _, err := PiHexDigits(MaxPosition+1, 4); err == nil {
		t.Error("expected an error for a position past the maximum")
	}
	if _, err := PiHexDigits(0, 0); err == nil {
		t.Error("expected an error for a zero digit count")
	}
	if _, err := PiHexDigits(1000000, MaxDigits(1000000)+1); err == nil {
		t.Error("expected an error for more digits than the precision allows")
	}
}

func TestPowMod16(t *testing.T) {
	sixteen := big.NewInt(16)
	for _, m := range []uint64{1, 2, 3, 7, 16, 17, 1000003, 1<<32 - 1} {
		for _, exp := range []uint64{0, 1, 2, 5, 63, 64, 1000, 999999} {
			want := new(big.Int).Exp(sixteen, new(big.Int).SetUint64(exp), new(big.Int).SetUint64(m)).Uint64()
			if got := powMod16(exp, m); got != want {
				t.Errorf("powMod16(%d, %d) = %d, want %d", exp, m, got, want)
			}
		}
	}
}

func TestMaxDigitsPastMaxPosition(t *testing.T) {
	if got := MaxDigits(MaxPosition + 1); got != 0 {
		t.Errorf("MaxDigits(MaxPosition+1) = %d, want 0", got)
	}
}
