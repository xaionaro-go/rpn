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

4 approaches were implemented, and the fastest is `tokenslice`:

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/rpn/tests
BenchmarkExpr_Eval/ideal-4                                              1000000000               1.16 ns/op
BenchmarkExpr_Eval/idealFuncs-4                                         1000000000               4.60 ns/op
BenchmarkExpr_Eval/cache_true/const/compile-4                           1000000000               4.01 ns/op            0 B/op          0 allocs/op
BenchmarkExpr_Eval/cache_true/const/tokenslice-4                        1000000000               3.99 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/const/calltree-4                          1000000000               3.99 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/const/exprtree-4                          1000000000               3.76 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/variable/calltree-4                       1000000000               4.00 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/variable/exprtree-4                       1000000000               3.77 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/variable/compile-4                        1000000000               4.00 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/variable/tokenslice-4                     1000000000               3.99 ns/op            0 B/op        0 allocs/op
BenchmarkExpr_Eval/cache_true/tons_of_variables/tokenslice-4            1000000000               4.00 ns/op            0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_true/tons_of_variables/calltree-4              1000000000               4.05 ns/op            0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_true/tons_of_variables/exprtree-4              1000000000               3.78 ns/op            0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/const/calltree-4                         1000000000               5.25 ns/op            0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/const/exprtree-4                         369319456               16.2 ns/op             0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/const/compile-4                          449603368               13.4 ns/op             0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/const/tokenslice-4                       810901628                7.41 ns/op            0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/variable/calltree-4                      489537325               12.6 ns/op             0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/variable/exprtree-4                      258942452               22.8 ns/op             0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/variable/compile-4                       232903180               24.8 ns/op             0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/variable/tokenslice-4                    309449833               19.4 ns/op             0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/tons_of_variables/calltree-4                65892             90933 ns/op               0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/tons_of_variables/exprtree-4                59335            100775 ns/op               0 B/op         0 allocs/op
BenchmarkExpr_Eval/cache_false/tons_of_variables/tokenslice-4             136764             43815 ns/op               0 B/op         0 allocs/op
PASS
ok      github.com/xaionaro-go/rpn/tests        134.308s
```

The default approach is also the `tokenslice`, so if you will import
`github.com/xaionaro-go/rpn` then if you will use it.  

Approach `compile` has potential if the assembly code will be re-written by somebody who
is good at optimizing on an amd64 assembly language. Right now is more like
an unsafe proof of concept.

