package rpn_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	rpn "github.com/xaionaro-go/rpn/implementations/tokenslice"
	"github.com/xaionaro-go/rpn/tests"
)

func TestBugCase0(t *testing.T) {
	exprString := "x0 1e2 if z - x1 0.5 + -"
	expr, err := rpn.Parse(exprString, tests.DummyResolver{})
	require.NoError(t, err)
	require.Equal(t, 95.5, expr.Eval(), exprString+": "+expr.String())
}

func TestBugCase1(t *testing.T) {
	exprString := "0.5 0 -1 * x1 * y - -1 if + y 1 +  /"
	expr, err := rpn.Parse(exprString, tests.DummyResolver{})
	require.NoError(t, err)
	require.Equal(t, 0.1, expr.Eval(), exprString+": "+expr.String())
}
