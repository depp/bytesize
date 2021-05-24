package bytesize

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

var (
	errMissingNumber   = errors.New("missing number")
	errMultipleDecimal = errors.New("multiple decimal points")
	errUnknownUnits    = errors.New("unknown units")

	// ErrRange indicates that the byte size is too large to fit in a uint64.
	ErrRange = errors.New("byte size out of range")
)

// An ParseError is an error parsing a byte size.
type ParseError struct {
	Value string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse %q: %v", e.Value, e.Err)
}

var prefixes = [...]uint32{
	'K': 1,
	'k': 1,
	'M': 2,
	'm': 2,
	'G': 3,
	'g': 3,
	'T': 4,
	't': 4,
	'P': 5,
	'p': 5,
	'E': 6,
	'e': 6,
	'Z': 7,
	'z': 7,
	'Y': 8,
	'y': 8,
}

// Parse parses a string as a byte size, which may contain units and may use a
// decimal point. If the result would overflow a uint64, returns a ParseError
// with an ErrRange inside. Results are rounded to even, although results for
// binary prefixes may not be exactly rounded. The units are the standard SI
// units as well as their binary counterparts, optionally followed by "b". ASCII
// space characters may appear between the number and the units.
//
//     Parse("1") = 1
//     Parse("555k") = 555000
//     Parse("15 EiB") = 17293822569102704640
//     Parse("1.5 mb") = 1500000
//     Parse("2gi") = 2147483648
//     Parse("0.001 zb") = 1000000000000000000
func Parse(s string) (uint64, error) {
	const max = ^uint64(0)

	// Split into number and units.
	spos := -1
	ppos := -1
	var hasd bool
scan:
	for i, c := range s {
		switch {
		case '0' <= c && c <= '9':
			hasd = true
		case c == '.':
			if ppos != -1 {
				return 0, &ParseError{s, errMultipleDecimal}
			}
			ppos = i
		default:
			spos = i
			break scan
		}
	}
	if !hasd {
		return 0, &ParseError{s, errMissingNumber}
	}
	var num, units string
	if spos == -1 {
		num = s
	} else {
		num = s[:spos]
		units = s[spos:]
	}
	for len(units) > 0 && units[0] == ' ' {
		units = units[1:]
	}

	// Parse the units.
	if len(units) > 0 {
		if c := units[len(units)-1]; c == 'b' || c == 'B' {
			units = units[:len(units)-1]
		}
	}
	var binary bool
	var scale int
	if len(units) > 0 {
		if len(units) == 2 {
			if c := units[1]; c == 'i' || c == 'I' {
				binary = true
			}
		} else if len(units) != 1 {
			return 0, &ParseError{s, errUnknownUnits}
		}
		u := units[0]
		if int(u) >= len(prefixes) {
			return 0, &ParseError{s, errUnknownUnits}
		}
		scale = int(prefixes[u])
		if scale == 0 {
			return 0, &ParseError{s, errUnknownUnits}
		}
	}

	// Parse a binary unit.
	if binary {
		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0, &ParseError{s, err}
		}
		f = math.Ldexp(f, scale*10)
		if f >= 1<<64 {
			return max, &ParseError{s, ErrRange}
		}
		return uint64(math.RoundToEven(f)), nil
	}

	// Parse a decimal unit.
	place := scale * 3
	if ppos != -1 {
		place += ppos
	} else {
		place += len(num)
	}
	var v uint64
	var frac string
	for i, c := range num {
		if c == '.' {
			continue
		}
		if place <= 0 {
			frac = num[i:]
			break
		}
		place--
		if v > max/10 {
			return max, &ParseError{s, ErrRange}
		}
		v *= 10
		o := v
		v += uint64(c - '0')
		if o > v {
			return max, &ParseError{s, ErrRange}
		}
	}
	for i := 0; i < place; i++ {
		if v > max/10 {
			return max, &ParseError{s, ErrRange}
		}
		v *= 10
	}
	// Just for consistency, round to even.
	if len(frac) != 0 {
		var roundUp bool
		if frac[0] > '5' {
			roundUp = true
		} else if frac[0] == '5' {
			if v&1 == 1 {
				roundUp = true
			} else {
				for _, c := range frac[1:] {
					if c != '0' {
						roundUp = true
						break
					}
				}
			}
		}
		if roundUp {
			v++
			if v == 0 {
				return max, &ParseError{s, ErrRange}
			}
		}
	}
	return v, nil
}
