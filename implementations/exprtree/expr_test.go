package rpn

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/rpn/tests"
)

func TestExpr_ClearCache(t *testing.T) {
	expr, err := Parse("y x0 x1 + +", tests.DummyResolver{T: t})
	require.NoError(t, err)

	expr.EnableMemoization(true)
	expr.Eval()
	require.Equal(t, true, expr.ResultCache.Valid)
	require.Equal(t, true, expr.RHS.ResultCache.Valid)
	expr.ClearCache()
	require.Equal(t, false, expr.ResultCache.Valid)
	require.Equal(t, false, expr.RHS.ResultCache.Valid)
	expr.EnableUpdateCache(false)
	expr.Eval()
	require.Equal(t, false, expr.ResultCache.Valid)
	require.Equal(t, false, expr.RHS.ResultCache.Valid)
}
