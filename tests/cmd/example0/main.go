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
