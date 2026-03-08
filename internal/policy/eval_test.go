package policy

import (
	"testing"

	"github.com/casbin/govaluate"
	"github.com/stretchr/testify/assert"
)

func makeTestEdge(cond string) *Edge {
	valid := cond == "" || isValidCond(cond)
	var compiled *govaluate.EvaluableExpression
	if valid && cond != "" {
		compiled, _ = govaluate.NewEvaluableExpression(cond)
	}
	return &Edge{Cond: cond, ValidCond: valid, CompiledCond: compiled}
}

func TestEvalEdgeCondition(t *testing.T) {
	t.Run("empty cond is true", func(t *testing.T) {
		edge := makeTestEdge("")
		vars := map[string]any{"x": 1}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("age>=18 true", func(t *testing.T) {
		edge := makeTestEdge("age>=18")
		vars := map[string]any{"age": 20}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("age>=18 false", func(t *testing.T) {
		edge := makeTestEdge("age>=18")
		vars := map[string]any{"age": 15}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.False(t, got)
	})
	t.Run("x==1 true", func(t *testing.T) {
		edge := makeTestEdge("x==1")
		vars := map[string]any{"x": 1.0}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("invalid expr returns false and error", func(t *testing.T) {
		edge := makeTestEdge("invalid!!!")
		vars := map[string]any{}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("bare number is invalid condition", func(t *testing.T) {
		edge := makeTestEdge("1")
		vars := map[string]any{}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic plus returns invalid", func(t *testing.T) {
		edge := makeTestEdge("age+1>=18")
		vars := map[string]any{"age": 17}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic star returns invalid", func(t *testing.T) {
		edge := makeTestEdge("score*2>100")
		vars := map[string]any{"score": 60}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("allowed && and ||", func(t *testing.T) {
		edge := makeTestEdge("age>=18 && score>700")
		vars := map[string]any{"age": 25, "score": 720}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("allowed comparison with negative number", func(t *testing.T) {
		edge := makeTestEdge("score> -1")
		vars := map[string]any{"score": 100}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("allowed string literal", func(t *testing.T) {
		edge := makeTestEdge(`role=="admin"`)
		vars := map[string]any{"role": "admin"}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("allowed boolean literal", func(t *testing.T) {
		edge := makeTestEdge("active==true")
		vars := map[string]any{"active": true}

		got, err := EvalEdgeCondition(edge, vars)

		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("arithmetic plus with spaces invalid", func(t *testing.T) {
		edge := makeTestEdge("age + 1 >= 18")
		vars := map[string]any{"age": 17}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic star invalid", func(t *testing.T) {
		edge := makeTestEdge("score*2 == 1400")
		vars := map[string]any{"score": 700}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic slash invalid", func(t *testing.T) {
		edge := makeTestEdge("points / 2 > 100")
		vars := map[string]any{"points": 250}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic minus invalid", func(t *testing.T) {
		edge := makeTestEdge("balance - 50 > 0")
		vars := map[string]any{"balance": 100}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("triple equals invalid", func(t *testing.T) {
		edge := makeTestEdge(`status === "ok"`)
		vars := map[string]any{"status": "ok"}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("angle bracket not equals invalid", func(t *testing.T) {
		edge := makeTestEdge("age <> 18")
		vars := map[string]any{"age": 20}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arrow equals invalid", func(t *testing.T) {
		edge := makeTestEdge("score => 800")
		vars := map[string]any{"score": 800}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("caret equals invalid", func(t *testing.T) {
		edge := makeTestEdge("active ^= true")
		vars := map[string]any{"active": true}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("parentheses not supported", func(t *testing.T) {
		edge := makeTestEdge("age>=18 && (score>700)")
		vars := map[string]any{"age": 25, "score": 720}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("parentheses around comparisons invalid", func(t *testing.T) {
		edge := makeTestEdge(`(name == "João") || (age < 20)`)
		vars := map[string]any{"name": "João", "age": 18}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("function call len invalid", func(t *testing.T) {
		edge := makeTestEdge("len(name) > 5")
		vars := map[string]any{"name": "Alice"}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("function call isAdult invalid", func(t *testing.T) {
		edge := makeTestEdge("isAdult(age) == true")
		vars := map[string]any{"age": 20}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("NULL value invalid", func(t *testing.T) {
		edge := makeTestEdge("height == NULL")
		vars := map[string]any{}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("undefined value invalid", func(t *testing.T) {
		edge := makeTestEdge("name == undefined")
		vars := map[string]any{}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("single ampersand invalid", func(t *testing.T) {
		edge := makeTestEdge("age>=18 & score>700")
		vars := map[string]any{"age": 25, "score": 720}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("triple pipe invalid", func(t *testing.T) {
		edge := makeTestEdge("age>=18 ||| score>700")
		vars := map[string]any{"age": 25, "score": 720}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("string without quotes invalid", func(t *testing.T) {
		edge := makeTestEdge("name == Joao")
		vars := map[string]any{"name": "Joao"}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("identifier starting with digit invalid", func(t *testing.T) {
		edge := makeTestEdge("1age == 10")
		vars := map[string]any{"1age": 10}

		got, err := EvalEdgeCondition(edge, vars)

		assert.Error(t, err)
		assert.False(t, got)
	})
}

func TestPreParseResult(t *testing.T) {
	t.Run("empty result returns nil", func(t *testing.T) {
		got := preParseResult("")

		assert.Nil(t, got)
	})
	t.Run("key=value sets string", func(t *testing.T) {
		got := preParseResult("name=foo")

		assert.Equal(t, "foo", got["name"])
	})
	t.Run("key=true sets bool", func(t *testing.T) {
		got := preParseResult("approved=true")

		assert.Equal(t, true, got["approved"])
	})
	t.Run("multiple pairs", func(t *testing.T) {
		got := preParseResult("num=2.5, flag=true")

		assert.Equal(t, 2.5, got["num"])
		assert.Equal(t, true, got["flag"])
	})
	t.Run("malformed pair without equals is skipped", func(t *testing.T) {
		got := preParseResult("a=1, badpair, b=2")

		assert.Equal(t, true, got["a"])
		assert.Equal(t, 2.0, got["b"])
		assert.Len(t, got, 2)
	})
}

func TestApplyParsedResult(t *testing.T) {
	t.Run("nil parsed does nothing", func(t *testing.T) {
		vars := map[string]any{"a": 1}

		applyParsedResult(nil, vars)

		assert.Len(t, vars, 1)
		assert.Equal(t, 1, vars["a"])
	})
	t.Run("applies parsed values to vars", func(t *testing.T) {
		parsed := map[string]any{"approved": true, "segment": "prime"}
		vars := map[string]any{"age": 25}

		applyParsedResult(parsed, vars)

		assert.Equal(t, 25, vars["age"])
		assert.Equal(t, true, vars["approved"])
		assert.Equal(t, "prime", vars["segment"])
	})
}
