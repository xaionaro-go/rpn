package tests_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/rpn"
	callslice "github.com/xaionaro-go/rpn/implementations/callslice"
	calltree "github.com/xaionaro-go/rpn/implementations/calltree"
	compile "github.com/xaionaro-go/rpn/implementations/compile"
	exprtree "github.com/xaionaro-go/rpn/implementations/exprtree"
	tokenslice "github.com/xaionaro-go/rpn/implementations/tokenslice"
	"github.com/xaionaro-go/rpn/tests"
	"github.com/xaionaro-go/rpn/types"
)

var implementations = map[string]func(string, types.SymbolResolver) (types.Expr, error){
	"callslice": func(s string, resolver types.SymbolResolver) (types.Expr, error) {
		return callslice.Parse(s, resolver)
	},
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
	"default": func(s string, resolver types.SymbolResolver) (types.Expr, error) {
		return rpn.Parse(s, resolver)
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
				expr, err := impl("y x0 x1 + *", tests.DummyResolver{T: t})
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
							expr, err := impl(rpn, tests.DummyResolver{T: t})
							require.NoError(t, err)
							require.Equal(t, float64(10000), expr.Eval(), implName)
						})
					}
				})
			}
		})
	}

	t.Run("extra_tests", func(t *testing.T) {
		for _, args := range []string{
			"h2 1",
			"0x2 b11",
			"0 1",
			"0 x0",
			"1 1",
			"1 x0",
			"x0 1",
			"x0 y",
		} {
			for _, memoization := range []bool{false, true} {
				for _, op := range []string{"+", "-", "*", "/", "^", "if"} {
					rpn := "0x1 " + args + " " + op + " +"
					resultMap := map[string]float64{}
					for implName, impl := range implementations {
						if (op == "^" || op == "if") && implName == "compile" {
							continue
						}
						expr, err := impl(rpn, tests.DummyResolver{T: t})
						require.NoError(t, err, fmt.Sprintf("%s: '%s'", implName, rpn))
						expr.EnableMemoization(memoization)
						require.NotEmpty(t, expr.String())
						resultMap[implName] = expr.Eval()
					}

					reference := resultMap["default"]
					for _, value := range resultMap {
						require.Equal(t, reference, value, fmt.Sprintf("'%s' -> %v", rpn, resultMap))
					}
				}
			}
		}
	})

	t.Run("random_expressions", func(t *testing.T) {
		randGen := rand.New(rand.NewSource(0))
		for i := 0; i < 10000; i++ {
			exprString := randExpression(randGen)
			resultMap := map[string]float64{}
			for implName, impl := range implementations {
				if implName == "compile" {
					continue
				}

				expr, err := impl(exprString, tests.DummyResolver{T: t})
				if err != nil {
					resultMap[implName] = -1
					return
				}
				resultMap[implName] = expr.Eval()
			}
			reference := resultMap["default"]
			for _, value := range resultMap {
				require.Equal(t, reference, value, fmt.Sprintf("'%s' -> %v", exprString, resultMap))
			}
		}
	})
}

func randExpression(randGen *rand.Rand) string {
	valDict := []string{
		"x0", "x1", "y", "z",
		"0", "1", "-1", "0.5", "1e2",
	}
	amountOfOps := randGen.Intn(10)
	collection := []string{
		"",
		valDict[randGen.Intn(len(valDict))],
	}
	for i := 0; i < amountOfOps; i++ {
		collection = append(collection, (types.OpPlus + types.Op(randGen.Intn(int(types.BoundaryOp-types.OpPlus)))).String())
		collection = append(collection, valDict[randGen.Intn(len(valDict))])
	}
	rand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})
	return strings.Join(collection, " ")
}

func BenchmarkExpr_Eval(b *testing.B) {
	for _, enableMemoization := range []bool{true, false} {
		b.Run(fmt.Sprintf("cache_%v", enableMemoization), func(b *testing.B) {
			for _, exprName := range []string{"3_constant_values", "3_variables", "10_variables", "100_variables", "1000_variables", "10000_variables"} {
				b.Run(exprName, func(b *testing.B) {
					exprString := func() string {
						switch exprName {
						case "3_constant_values":
							return "b10 3.5 4 + *"
						case "3_variables":
							return "z x0 x1 + *"
						case "10_variables":
							return strings.Repeat("z ", 10) + strings.Repeat("+ ", 9)
						case "100_variables":
							return strings.Repeat("z ", 100) + strings.Repeat("+ ", 99)
						case "1000_variables":
							return strings.Repeat("z ", 1000) + strings.Repeat("+ ", 999)
						case "10000_variables":
							return strings.Repeat("z ", 10000) + strings.Repeat("+ ", 9999)
						}
						panic("should not happen")
					}()
					for implName, impl := range implementations {
						expr, err := impl(exprString, tests.DummyResolver{})
						if err != nil {
							panic(err)
						}
						eval := expr.Eval
						expr.EnableMemoization(enableMemoization)
						b.Run(implName, func(b *testing.B) {
							b.ReportAllocs()
							b.ResetTimer()
							for i := 0; i < b.N; i++ {
								eval()
							}
						})
					}
				})
			}
		})
	}
}
