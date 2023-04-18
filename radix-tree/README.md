# Benchmark

## Run(Once)

```
$ go test -bench . -benchmem
```

### Result

```
BenchmarkRegexFirst-10          11028405                97.70 ns/op            0 B/op          0 allocs/op
BenchmarkRadixFirst-10          31994204                36.53 ns/op           32 B/op          1 allocs/op
BenchmarkRegexMid-10             3229677               363.4 ns/op             0 B/op          0 allocs/op
BenchmarkRadixMid-10            17497344                67.78 ns/op           64 B/op          1 allocs/op
BenchmarkRegexLast-10            2487745               482.7 ns/op             0 B/op          0 allocs/op
BenchmarkRadixLast-10           12847460                94.01 ns/op           80 B/op          1 allocs/op
BenchmarkRegexRound-10           1506157               791.6 ns/op             0 B/op          0 allocs/op
BenchmarkRadixRound-10          19155417                61.50 ns/op           51 B/op          1 allocs/op
```

## Run(Stats)

```
$ go test -bench . -benchmem -count 10 > out
```

### Result

```
$ benchstat out
name           time/op
RegexFirst-10  99.3ns ± 0%
RadixFirst-10  37.0ns ± 1%
RegexMid-10     369ns ± 1%
RadixMid-10    69.5ns ± 1%
RegexLast-10    492ns ± 3%
RadixLast-10   95.4ns ± 3%
RegexRound-10   803ns ± 2%
RadixRound-10  62.4ns ± 1%

name           alloc/op
RegexFirst-10   0.00B
RadixFirst-10   32.0B ± 0%
RegexMid-10     0.00B
RadixMid-10     64.0B ± 0%
RegexLast-10    0.00B
RadixLast-10    80.0B ± 0%
RegexRound-10   0.00B
RadixRound-10   51.0B ± 0%

name           allocs/op
RegexFirst-10    0.00
RadixFirst-10    1.00 ± 0%
RegexMid-10      0.00
RadixMid-10      1.00 ± 0%
RegexLast-10     0.00
RadixLast-10     1.00 ± 0%
RegexRound-10    0.00
RadixRound-10    1.00 ± 0%
```
