package tests_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	calltree "github.com/xaionaro-go/rpn/implementations/calltree"
	compile "github.com/xaionaro-go/rpn/implementations/compile"
	exprtree "github.com/xaionaro-go/rpn/implementations/exprtree"
	tokenslice "github.com/xaionaro-go/rpn/implementations/tokenslice"
	"github.com/xaionaro-go/rpn/types"
)

type dummyResolver struct {
	t *testing.T
}

func (r dummyResolver) Resolve(sym string) (types.ValueLoader, error) {
	switch sym {
	case "x0":
		return types.FuncValue(func() float64 {
			return 2
		}), nil
	case "x1":
		return types.FuncValue(func() float64 {
			return 3
		}), nil
	case "y":
		return types.StaticValue(4), nil
	case "z":
		return types.FuncValue(func() float64 {
			return 1
		}), nil
	}
	require.FailNow(r.t, fmt.Sprintf("should not happen: '%s'", sym))
	return nil, nil
}

var implementations = map[string]func(string, types.SymbolResolver) (types.Expr, error){
	"calltree": func(s string, resolver types.SymbolResolver) (types.Expr, error) {
		return calltree.Parse(s, resolver)
	},
	"exprtree": func(s string, resolver types.SymbolResolver) (types.Expr, error) {
		return exprtree.Parse(s, resolver)
	},
	"compile": func(s string, resolver types.SymbolResolver) (types.Expr, error) {
		return compile.Parse(s, resolver)
	},
	"tokenslice": func(s string, resolver types.SymbolResolver) (types.Expr, error) {
		return tokenslice.Parse(s, resolver)
	},
}

func TestExpr(t *testing.T) {
	for implName, impl := range implementations {
		t.Run(implName, func(t *testing.T) {
			t.Run("const", func(t *testing.T) {
				expr, err := impl("b10 3.5 4 + *", nil)
				require.NoError(t, err)
				require.Equal(t, float64(15), expr.Eval(), fmt.Sprintf("%s: '%s'", implName, expr.String()))
			})
			t.Run("syms", func(t *testing.T) {
				expr, err := impl("y x0 x1 + *", dummyResolver{t: t})
				require.NoError(t, err)
				require.Equal(t, float64(20), expr.Eval(), fmt.Sprintf("%s: '%s'", implName, expr.String()))
			})
			if implName != "compile" {
				t.Run("large_expression", func(t *testing.T) {
					for _, sym := range []string{"1", "z"} {
						var description string
						if sym == "1" {
							description = "const"
						} else {
							description = "variable"
						}
						t.Run(description, func(t *testing.T) {
							rpn := strings.Repeat(sym+" ", 10000) + strings.Repeat("+ ", 9999)
							expr, err := impl(rpn, &dummyResolver{t: t})
							require.NoError(t, err)
							require.Equal(t, float64(10000), expr.Eval(), implName)
						})
					}
				})
			}
		})
	}

	t.Run("full_test", func(t *testing.T) {
		for _, args := range []string{
			"0 1",
			"0 x0",
			"1 1",
			"1 x0",
			"x0 1",
			"x0 y",
		} {
			for _, op := range []string{"+", "-", "*", "/", "^", "if"} {
				rpn := "0x1 " + args + " " + op + " +"
				resultMap := map[string]float64{}
				for implName, impl := range implementations {
					if (op == "^" || op == "if") && implName == "compile" {
						continue
					}
					expr, err := impl(rpn, dummyResolver{t: t})
					require.NoError(t, err, fmt.Sprintf("%s: '%s'", implName, rpn))
					resultMap[implName] = expr.Eval()
				}

				reference := resultMap["calltree"]
				for _, value := range resultMap {
					require.Equal(t, reference, value, fmt.Sprintf("'%s' -> %v", rpn, resultMap))
				}
			}
		}
	})
}

var (
	dummyA, dummyB, dummyC, dummyD     = float64(3.5), float64(4), float64(2), float64(0)
	dummyFuncA, dummyFuncB, dummyFuncC = func() float64 { return 3.5 }, func() float64 { return 4 }, func() float64 { return 2 }
)

func BenchmarkExpr_Eval(b *testing.B) {
	b.Run("ideal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dummyD = (dummyA + dummyB) * dummyC
		}
	})
	b.Run("idealFuncs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dummyD = (dummyFuncA() + dummyFuncB()) * dummyFuncC()
		}
	})

	for implName, impl := range implementations {
		b.Run(implName, func(b *testing.B) {
			b.Run("const", func(b *testing.B) {
				expr, _ := impl("b10 3.5 4 + *", nil)
				eval := expr.Eval
				expr.EnableMemoization(true)
				b.Run("with_cache", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						eval()
					}
				})
				expr.EnableMemoization(false)
				b.Run("without_cache", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						eval()
					}
				})
			})
			b.Run("variable", func(b *testing.B) {
				expr, _ := impl("z x0 x1 + *", dummyResolver{t: nil})
				eval := expr.Eval
				expr.EnableMemoization(true)
				b.Run("with_cache", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						eval()
					}
				})
				expr.EnableMemoization(false)
				b.Run("without_cache", func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						eval()
					}
				})
			})
			if implName != "compile" {
				b.Run("large_expression", func(b *testing.B) {
					for _, sym := range []string{"1", "z"} {
						var description string
						if sym == "1" {
							description = "const"
						} else {
							description = "variable"
						}
						b.Run(description, func(b *testing.B) {
							rpn := strings.Repeat(sym+" ", 10000) + strings.Repeat("+ ", 9999)
							expr, _ := impl(rpn, &dummyResolver{t: nil})
							eval := expr.Eval
							expr.EnableMemoization(false)
							b.Run("without_cache", func(b *testing.B) {
								b.SetBytes(10000)
								b.ReportAllocs()
								b.ResetTimer()
								for i := 0; i < b.N; i++ {
									eval()
								}
							})
						})
					}
				})
			}
		})
	}
}
