package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	linearDOT := `digraph { start [result=""]; ok [result="approved=true"]; no [result="approved=false"]; start -> ok [cond="age>=18"]; start -> no [cond="age<18"]; }`
	cycleDOT := `digraph { start [result=""]; a [result="done=true"]; start -> a [cond="x==1"]; a -> a [cond="x==1"]; }`
	singleNodeDOT := `digraph { start [result="x=1"]; }`
	edgeToMissingNodeDOT := `digraph { start [result="done=true"]; start -> ghost [cond="x==1"]; }`

	t.Run("single path applies node results", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()
		executor := NewGraphExecutor()
		graph, err := parser.Parse(context.Background(), linearDOT)
		require.NoError(t, err)
		vars := map[string]any{"age": 20}

		// Act
		resp, err := executor.Process(context.Background(), graph, vars)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 20, resp.Output["age"])
		assert.Equal(t, true, resp.Output["approved"])
	})
	t.Run("single path age under 18", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()
		executor := NewGraphExecutor()
		graph, err := parser.Parse(context.Background(), linearDOT)
		require.NoError(t, err)
		vars := map[string]any{"age": 15}

		// Act
		resp, err := executor.Process(context.Background(), graph, vars)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, false, resp.Output["approved"])
	})
	t.Run("graph with cycle terminates and returns", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()
		executor := NewGraphExecutor()
		graph, err := parser.Parse(context.Background(), cycleDOT)
		require.NoError(t, err)
		vars := map[string]any{"x": 1.0}

		// Act
		resp, err := executor.Process(context.Background(), graph, vars)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1.0, resp.Output["x"])
		assert.Equal(t, true, resp.Output["done"])
	})
	t.Run("single node no edges empty input returns only node result", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()
		executor := NewGraphExecutor()
		graph, err := parser.Parse(context.Background(), singleNodeDOT)
		require.NoError(t, err)
		vars := map[string]any{}

		// Act
		resp, err := executor.Process(context.Background(), graph, vars)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"x": true}, resp.Output)
	})
	t.Run("edge to missing node applies start result then stops", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()
		executor := NewGraphExecutor()
		graph, err := parser.Parse(context.Background(), edgeToMissingNodeDOT)
		require.NoError(t, err)
		vars := map[string]any{"x": 1}

		// Act
		resp, err := executor.Process(context.Background(), graph, vars)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"done": true, "x": 1}, resp.Output)
	})
}
