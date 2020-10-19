// Package bytesize formats data sizes as strings.
package bytesize

import (
	"fmt"
)

// Format formats a data size for humans. The size is rounded to three
// significant digits, but never more than integer accuracy, and printed with SI
// prefixes, with the unit B for bytes.
//
// For example, 1000 is formatted as "1.00 kB".
//
// This uses unbiased, round-to-even rounding just like the standard library
// does (math.Round, strconv.FormatFloat). So 1005 is "1.00 kB", and 1015 is
// "1.02 kB".
func Format(size uint64) string {
	const prefixes = "kMGTPEZY"
	if size <= 0 {
		return "0 B"
	}
	if size < 1000 {
		return fmt.Sprintf("%d B", size)
	}
	var hasRem bool
	var rem, pfx uint32
	for pfx = 0; ; pfx++ {
		rem = uint32(size % 1000)
		size /= 1000
		if size < 1000 || int(pfx)+1 == len(prefixes) {
			break
		}
		if rem > 0 {
			hasRem = true
		}
	}
	n := int(size)
	if n < 10 {
		m := rem / 10
		rem = rem % 10
		if rem > 5 || (rem == 5 && (m&1 != 0 || hasRem)) {
			m++
			if m == 100 {
				m = 0
				n++
				if n == 10 {
					return fmt.Sprintf("10.0 %cB", prefixes[pfx])
				}
			}
		}
		return fmt.Sprintf("%d.%02d %cB", n, m, prefixes[pfx])
	}
	if n < 100 {
		m := rem / 100
		rem = rem % 100
		if rem > 50 || (rem == 50 && (m&1 != 0 || hasRem)) {
			m++
			if m == 10 {
				m = 0
				n++
				if n == 100 {
					return fmt.Sprintf("100 %cB", prefixes[pfx])
				}
			}
		}
		return fmt.Sprintf("%d.%d %cB", n, m, prefixes[pfx])
	}
	if rem > 500 || (rem == 500 && (n&1 != 0 || hasRem)) {
		n++
	}
	if n >= 1000 && int(pfx)+1 < len(prefixes) {
		return fmt.Sprintf("1.00 %cB", prefixes[pfx+1])
	}
	return fmt.Sprintf("%d %cB", n, prefixes[pfx])
}
