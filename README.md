[![GoDoc](https://godoc.org/github.com/xaionaro-go/rpn?status.svg)](https://pkg.go.dev/github.com/xaionaro-go/rpn?tab=doc)
[![go report](https://goreportcard.com/badge/github.com/xaionaro-go/rpn)](https://goreportcard.com/report/github.com/xaionaro-go/rpn)
[![Build Status](https://travis-ci.org/xaionaro-go/rpn.svg?branch=master)](https://travis-ci.org/xaionaro-go/rpn)
[![Coverage Status](https://coveralls.io/repos/github/xaionaro-go/rpn/badge.svg?branch=master)](https://coveralls.io/github/xaionaro-go/rpn?branch=master)
<p xmlns:dct="http://purl.org/dc/terms/" xmlns:vcard="http://www.w3.org/2001/vcard-rdf/3.0#">
  <a rel="license"
     href="http://creativecommons.org/publicdomain/zero/1.0/">
    <img src="http://i.creativecommons.org/p/zero/1.0/88x31.png" style="border-style: none;" alt="CC0" />
  </a>
  <br />
  To the extent possible under law,
  <a rel="dct:publisher"
     href="https://github.com/xaionaro-go/rpn">
    <span property="dct:title">Dmitrii Okunev</span></a>
  has waived all copyright and related or neighboring rights to
  <span property="dct:title">Reverse Polish Notation for Go</span>.
This work is published from:
<span property="vcard:Country" datatype="dct:ISO3166"
      content="IE" about="https://github.com/xaionaro-go/rpn">
  Ireland</span>.
</p>

# About

`github.com/xaionaro-go/rpn` is an implementation of [Reverse Polish Notation](https://en.wikipedia.org/wiki/Reverse_Polish_notation).
The implementation is focused on fast evaluation (but slow parsing).

# Quick start

```go
package main

import (
	"fmt"

	"github.com/xaionaro-go/rpn"
	"github.com/xaionaro-go/rpn/types"
)

type variables struct {
	X float64
}

func (r *variables) Resolve(sym string) (types.ValueLoader, error) {
	switch sym {
	case "x":
		return types.FuncValue(func() float64 {
			return r.X
		}), nil
	}
	return nil, fmt.Errorf("symbol '%s' not found", sym)
}

func main() {
	vars := &variables{}
	expr, err := rpn.Parse("x 2 *", vars)
	if err != nil {
		panic(err)
	}

	vars.X = 1
	_, _ = fmt.Println(expr.Eval())

	vars.X = 3
	_, _ = fmt.Println(expr.Eval())
}
```
The output is:
```
2
6
```

# Benchmark

There are 5 approaches implemented (`callslice`, `calltree`, `exprtree`, `compile` and `tokenslice`):

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/rpn/tests
BenchmarkExpr_Eval/cache_true/3_constant_values/callslice-4             1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_constant_values/calltree-4              1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_constant_values/exprtree-4              1000000000               3.72 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_constant_values/compile-4               1000000000               4.06 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_constant_values/tokenslice-4            1000000000               3.99 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_constant_values/default-4               1000000000               4.01 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_variables/default-4                     1000000000               3.99 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_variables/callslice-4                   1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_variables/calltree-4                    1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_variables/exprtree-4                    1000000000               3.73 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_variables/compile-4                     1000000000               4.08 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/3_variables/tokenslice-4                  1000000000               3.99 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10_variables/callslice-4                  1000000000               4.02 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10_variables/calltree-4                   1000000000               4.01 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10_variables/exprtree-4                   1000000000               3.72 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10_variables/compile-4                    1000000000               4.05 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10_variables/tokenslice-4                 1000000000               4.05 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10_variables/default-4                    1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/100_variables/callslice-4                 1000000000               4.03 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/100_variables/calltree-4                  1000000000               4.01 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/100_variables/exprtree-4                  1000000000               3.73 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/100_variables/compile-4                   1000000000               4.06 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/100_variables/tokenslice-4                1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/100_variables/default-4                   1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/1000_variables/exprtree-4                 1000000000               3.73 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/1000_variables/compile-4                  1000000000               4.04 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/1000_variables/tokenslice-4               1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/1000_variables/default-4                  1000000000               3.99 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/1000_variables/callslice-4                1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/1000_variables/calltree-4                 1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10000_variables/calltree-4                1000000000               4.06 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10000_variables/exprtree-4                1000000000               3.72 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10000_variables/compile-4                 1000000000               4.05 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10000_variables/tokenslice-4              1000000000               4.00 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10000_variables/default-4                 1000000000               3.99 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/10000_variables/callslice-4               1000000000               4.01 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_constant_values/callslice-4            579179106               10.4 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_constant_values/calltree-4             1000000000               5.51 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_constant_values/exprtree-4             376218780               15.9 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_constant_values/compile-4              451313248               13.3 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_constant_values/tokenslice-4           285631948               21.0 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_constant_values/default-4              1000000000               5.52 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_variables/callslice-4                  331783568               18.1 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_variables/calltree-4                   495298591               12.1 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_variables/exprtree-4                   287268447               20.9 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_variables/compile-4                    256272177               23.3 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_variables/tokenslice-4                 229319173               26.4 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/3_variables/default-4                    494261940               12.1 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10_variables/callslice-4                 100000000               52.9 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10_variables/calltree-4                  122009749               49.0 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10_variables/exprtree-4                  75151573                73.8 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10_variables/compile-4                   90459986                66.6 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10_variables/tokenslice-4                69911527                85.7 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10_variables/default-4                   121929210               48.9 ns/op             0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/100_variables/calltree-4                 10963188               548 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/100_variables/exprtree-4                  7401158               808 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/100_variables/compile-4                   8893951               678 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/100_variables/tokenslice-4                6937212               865 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/100_variables/default-4                  12024320               491 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/100_variables/callslice-4                11103020               529 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/1000_variables/callslice-4                1000000              5396 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/1000_variables/calltree-4                  885706              6453 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/1000_variables/exprtree-4                  650157              9065 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/1000_variables/compile-4                   876373              6829 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/1000_variables/tokenslice-4                694416              8591 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/1000_variables/default-4                  1000000              5400 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10000_variables/compile-4                   84819             70461 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10000_variables/tokenslice-4                68793             86768 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10000_variables/default-4                  104406             57272 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10000_variables/callslice-4                104284             57153 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10000_variables/calltree-4                  77030             77918 ns/op               0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_false/10000_variables/exprtree-4                  60081             99462 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/xaionaro-go/rpn/tests        407.949s
```

The default approach for small expressions is `calltree`, and
for larger expression is `callslice`, so if you will import
`github.com/xaionaro-go/rpn` then if you will use them.


