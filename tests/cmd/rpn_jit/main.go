// This "main" was added to faster debug the JIT implementation.

package main

import (
	"fmt"

	rpn "github.com/xaionaro-go/rpn/implementations/jit"
)

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	expr, err := rpn.Parse("b10 3.5 4 + *", nil)
	assertNoError(err)

	fmt.Println(expr.Eval())
}
