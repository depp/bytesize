package bytesize

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	const (
		sizek = 1000
		sizem = 1000 * sizek
		sizeg = 1000 * sizem
		sizet = 1000 * sizeg
		sizep = 1000 * sizet
		sizee = 1000 * sizep
	)

	const (
		sizeki = 1 << (10*iota + 10)
		sizemi
		sizegi
		sizeti
		sizepi
		sizeei
	)

	type testcase struct {
		num  string
		unit string
		out  uint64
	}
	cases := []testcase{
		// Integers + decimal units.
		{"0", "", 0},
		{"23", "", 23},
		{"103", "k", 103 * sizek},
		{"12", "m", 12 * sizem},
		{"3", "g", 3 * sizeg},
		{"715", "t", 715 * sizet},
		{"9", "p", 9 * sizep},
		{"5", "e", 5 * sizee},
		{"18446744073709551615", "", ^uint64(0)},

		// Fractions + decimal units.
		{"1.205", "k", 1205},
		{"16.6", "m", 16600 * sizek},
		{"18.446744073709551615", "e", ^uint64(0)},
		{"0.018", "z", 18 * sizee},

		// Fractions of a byte round to even.
		{"1.4", "", 1},
		{"1.5", "", 2},
		{"1.9", "", 2},
		{"2.1", "", 2},
		{"2.5000", "", 2},
		{"2.50001", "", 3},
		{"1.2306", "k", 1231},

		// Integers + binary units.
		{"103", "ki", 103 * sizeki},
		{"99", "mi", 99 * sizemi},
		{"15", "ei", 15 * sizeei},

		// Fractions + binary units.
		{"103.5", "ki", 103*sizeki + sizeki/2},
		{"593.2", "mi", 622015283},
	}
	for _, c := range cases {
		units := [...]string{
			c.unit,
			c.unit + "b",
			strings.ToUpper(c.unit) + "B",
			" " + c.unit,
		}
		for _, u := range units {
			if strings.HasSuffix(u, " ") {
				continue
			}
			in := c.num + u
			out, err := Parse(in)
			if err != nil {
				t.Errorf("Parse(%q): %v", in, err)
			} else if out != c.out {
				t.Errorf("Parse(%q) = %d, expect %d", in, out, c.out)
			}
		}
	}
}
