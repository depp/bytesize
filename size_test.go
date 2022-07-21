package bytesize

import (
	"math"
	"testing"
)

func TestFormat(t *testing.T) {
	type tcase struct {
		value  uint64
		output string
	}
	cases := []tcase{
		{0, "0 B"},
		{5, "5 B"},
		{20, "20 B"},
		{100, "100 B"},
		{500, "500 B"},
		{999, "999 B"},
		{1000, "1.00 kB"},
		{1005, "1.00 kB"},
		{1006, "1.01 kB"},
		{2334, "2.33 kB"},
		{2335, "2.34 kB"},
		{2995, "3.00 kB"},
		{9994, "9.99 kB"},
		{9995, "10.0 kB"},
		{10000, "10.0 kB"},
		{10050, "10.0 kB"},
		{10061, "10.1 kB"},
		{99949, "99.9 kB"},
		{99950, "100 kB"},
		{999499, "999 kB"},
		{999500, "1.00 MB"},
		{1000000, "1.00 MB"},
		{952500000, "952 MB"},
		{952500001, "953 MB"},
		{1000000000, "1.00 GB"},
		{2300000000000, "2.30 TB"},
		{15500000000000000, "15.5 PB"},
		{math.MaxInt64, "9.22 EB"},
		{math.MaxUint64, "18.4 EB"},
	}
	for _, c := range cases {
		out := Format(c.value)
		if out != c.output {
			t.Errorf("Format(%d): got %q, want %q", c.value, out, c.output)
		}
	}
}
