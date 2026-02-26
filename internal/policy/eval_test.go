package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalCondition(t *testing.T) {
	t.Run("empty cond is true", func(t *testing.T) {
		// Arrange
		cond := ""
		vars := map[string]any{"x": 1}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("age>=18 true", func(t *testing.T) {
		// Arrange
		cond := "age>=18"
		vars := map[string]any{"age": 20}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("age>=18 false", func(t *testing.T) {
		// Arrange
		cond := "age>=18"
		vars := map[string]any{"age": 15}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.False(t, got)
	})
	t.Run("x==1 true", func(t *testing.T) {
		// Arrange
		cond := "x==1"
		vars := map[string]any{"x": 1.0}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("invalid expr returns false and error", func(t *testing.T) {
		// Arrange
		cond := "invalid!!!"
		vars := map[string]any{}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("bare number is invalid condition", func(t *testing.T) {
		// Arrange
		cond := "1"
		vars := map[string]any{}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic plus returns invalid", func(t *testing.T) {
		// Arrange
		cond := "age+1>=18"
		vars := map[string]any{"age": 17}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic star returns invalid", func(t *testing.T) {
		// Arrange
		cond := "score*2>100"
		vars := map[string]any{"score": 60}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("allowed && and ||", func(t *testing.T) {
		// Arrange
		cond := "age>=18 && score>700"
		vars := map[string]any{"age": 25, "score": 720}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("allowed comparison with negative number", func(t *testing.T) {
		// Arrange
		cond := "score> -1"
		vars := map[string]any{"score": 100}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("allowed string literal", func(t *testing.T) {
		// Arrange
		cond := `role=="admin"`
		vars := map[string]any{"role": "admin"}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("allowed boolean literal", func(t *testing.T) {
		// Arrange
		cond := "active==true"
		vars := map[string]any{"active": true}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.NoError(t, err)
		assert.True(t, got)
	})
	t.Run("arithmetic plus with spaces invalid", func(t *testing.T) {
		// Arrange
		cond := "age + 1 >= 18"
		vars := map[string]any{"age": 17}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic star invalid", func(t *testing.T) {
		// Arrange
		cond := "score*2 == 1400"
		vars := map[string]any{"score": 700}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic slash invalid", func(t *testing.T) {
		// Arrange
		cond := "points / 2 > 100"
		vars := map[string]any{"points": 250}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arithmetic minus invalid", func(t *testing.T) {
		// Arrange
		cond := "balance - 50 > 0"
		vars := map[string]any{"balance": 100}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("triple equals invalid", func(t *testing.T) {
		// Arrange
		cond := `status === "ok"`
		vars := map[string]any{"status": "ok"}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("angle bracket not equals invalid", func(t *testing.T) {
		// Arrange
		cond := "age <> 18"
		vars := map[string]any{"age": 20}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("arrow equals invalid", func(t *testing.T) {
		// Arrange
		cond := "score => 800"
		vars := map[string]any{"score": 800}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("caret equals invalid", func(t *testing.T) {
		// Arrange
		cond := "active ^= true"
		vars := map[string]any{"active": true}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("parentheses not supported", func(t *testing.T) {
		// Arrange
		cond := "age>=18 && (score>700)"
		vars := map[string]any{"age": 25, "score": 720}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("parentheses around comparisons invalid", func(t *testing.T) {
		// Arrange
		cond := `(name == "João") || (age < 20)`
		vars := map[string]any{"name": "João", "age": 18}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("function call len invalid", func(t *testing.T) {
		// Arrange
		cond := "len(name) > 5"
		vars := map[string]any{"name": "Alice"}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("function call isAdult invalid", func(t *testing.T) {
		// Arrange
		cond := "isAdult(age) == true"
		vars := map[string]any{"age": 20}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("NULL value invalid", func(t *testing.T) {
		// Arrange
		cond := "height == NULL"
		vars := map[string]any{}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("undefined value invalid", func(t *testing.T) {
		// Arrange
		cond := "name == undefined"
		vars := map[string]any{}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("single ampersand invalid", func(t *testing.T) {
		// Arrange
		cond := "age>=18 & score>700"
		vars := map[string]any{"age": 25, "score": 720}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("triple pipe invalid", func(t *testing.T) {
		// Arrange
		cond := "age>=18 ||| score>700"
		vars := map[string]any{"age": 25, "score": 720}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("string without quotes invalid", func(t *testing.T) {
		// Arrange
		cond := "name == Joao"
		vars := map[string]any{"name": "Joao"}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
	t.Run("identifier starting with digit invalid", func(t *testing.T) {
		// Arrange
		cond := "1age == 10"
		vars := map[string]any{"1age": 10}

		// Act
		got, err := EvalCondition(cond, vars)

		// Assert
		assert.Error(t, err)
		assert.False(t, got)
	})
}

func TestApplyResult(t *testing.T) {
	t.Run("empty result does nothing", func(t *testing.T) {
		// Arrange
		result := ""
		vars := map[string]any{"a": 1}

		// Act
		ApplyResult(result, vars)

		// Assert
		assert.Len(t, vars, 1)
		assert.Equal(t, 1, vars["a"])
	})
	t.Run("key=value sets string", func(t *testing.T) {
		// Arrange
		result := "name=foo"
		vars := map[string]any{}

		// Act
		ApplyResult(result, vars)

		// Assert
		assert.Equal(t, "foo", vars["name"])
	})
	t.Run("key=true sets bool", func(t *testing.T) {
		// Arrange
		result := "approved=true"
		vars := map[string]any{}

		// Act
		ApplyResult(result, vars)

		// Assert
		assert.Equal(t, true, vars["approved"])
	})
	t.Run("multiple pairs", func(t *testing.T) {
		// Arrange
		result := "num=2.5, flag=true"
		vars := map[string]any{}

		// Act
		ApplyResult(result, vars)

		// Assert
		assert.Equal(t, 2.5, vars["num"])
		assert.Equal(t, true, vars["flag"])
	})
	t.Run("malformed pair without equals is skipped", func(t *testing.T) {
		// Arrange
		result := "a=1, badpair, b=2"
		vars := map[string]any{}

		// Act
		ApplyResult(result, vars)

		// Assert
		assert.Equal(t, true, vars["a"])
		assert.Equal(t, 2.0, vars["b"])
		assert.Len(t, vars, 2)
	})
}
