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

# Benchmark

4 approaches were implemented, and the fastest is `tokenslice`:

```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/rpn/tests
BenchmarkExpr_Eval/ideal-8                                              	1000000000	         0.675 ns/op
BenchmarkExpr_Eval/idealFuncs-8    	                                        1000000000	         5.07 ns/op
BenchmarkExpr_Eval/exprtree/const/with_cache-8         	                        1000000000	         4.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/compile/const/with_cache-8                              	1000000000	         4.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/tokenslice/const/with_cache-8                           	1000000000	         4.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/calltree/const/with_cache-8                             	1000000000	         4.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/exprtree/variable/with_cache-8      	                        1000000000	         4.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/compile/variable/with_cache-8                           	1000000000	         4.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/tokenslice/variable/with_cache-8                        	1000000000	         4.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/calltree/variable/with_cache-8                          	1000000000	         4.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/exprtree/const/without_cache-8      	                        288729301	        21.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/compile/const/without_cache-8                           	201709616	        29.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/tokenslice/const/without_cache-8                        	658899862	         9.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/calltree/const/without_cache-8                          	1000000000	         6.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/exprtree/variable/without_cache-8   	                        93561038	        56.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/compile/variable/without_cache-8                        	82377315	        71.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/tokenslice/variable/without_cache-8                     	114324116	        53.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/calltree/variable/without_cache-8                       	100000000	        57.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/exprtree/large_expression/const/without_cache-8         	   48636	    119083 ns/op	  83.97 MB/s	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/tokenslice/large_expression/const/without_cache-8       	600973670	         9.17 ns/op	1089965.99 MB/s	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/calltree/large_expression/const/without_cache-8         	977376576	         6.35 ns/op	1574754.60 MB/s	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/exprtree/large_expression/variable/without_cache-8      	   39877	    145013 ns/op	  68.96 MB/s	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/tokenslice/large_expression/variable/without_cache-8    	   82354	     68934 ns/op	 145.07 MB/s	       0 B/op	       0 allocs/op
BenchmarkExpr_Eval/calltree/large_expression/variable/without_cache-8      	   52723	    113742 ns/op	  87.92 MB/s	       0 B/op	       0 allocs/op
PASS
ok  	github.com/xaionaro-go/rpn/tests	149.776s
```
If compare pessimistic case `calltree/with_syms_cache/without_cache` with
ideal function-based reference benchmark `idealFuncs` then the overhead is
pretty small (according to the benchmark above).

The default approach is also the `tokenslice`, so if you will import
`github.com/xaionaro-go/rpn` then if you will use it.  

Approach `compile` has potential if the assembly code will be re-written by somebody who
is good at optimizing on an amd64 assembly language. Right now is more like
an unsafe proof of concept.

