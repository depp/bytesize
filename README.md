# ByteSize

**Simple** and **correct.**

ByteSize is a library for formatting and parsing numbers as byte sizes for humans. Numbers are formatted with three digits of precision, using SI prefixes, and the unit "B" for bytes. For example, 1000 is formatted as "1.00 KB". The numbers are rounded using round-to-even, which is the familiar method used by `fmt.Sprintf`, `math.Round`, and other functions in the standard library.

There are no choices to make. ByteSize is not configurable.

## Formatting

The precision cannot be changed. Non-decimal prefixes are not produced: no kibibytes, no powers of two. There is only one function to call.

    func Format(size uint64) string

All corner cases should be handled correctly and you should never see unusual or unexpected output. You should always see exactly three digits, except for inputs under 100.

Some test cases:

```
0 => "0 B"
999 => "999 B"
1000 => "1.00 kB"
1005 => "1.00 kB"
1006 => "1.01 kB"
1014 => "1.01 kB"
1015 => "1.02 kB"
9995 => "10.0 kB"
314000 => "314 kB"
math.MaxInt64 => "9.22 EB"
math.MaxUint64 => "18.4 EB"
```

## Parsing

The parser understands decimal numbers, decimal SI prefixes, and binary SI prefixes. It is not case-sensitive. ASCII space (0x20) may appear between the number and the units. This will recognize all positve prefixes which are powers of 1,000, including prefixes that will overflow (yotta and zetta).

    func Parse(s string) (uint64, error)

Some test cases:

    "0" => 0
    "1" => 1
    "555k" => 555000
    "15 EiB" => 17293822569102704640
    "1.5 mb" => 1500000
    "2gi" => 2147483648
    "0.001 zb" => 1000000000000000000

On overflow, this will return `^uint64(0)` and a `ParseError` containing `ErrRange`.

## License

This is licensed under the MIT license. See LICENSE.txt.
