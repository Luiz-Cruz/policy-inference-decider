package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDOT(t *testing.T) {
	validDOT := `digraph { start [result=""]; ok [result="approved=true"]; start -> ok [cond="age>=18"]; }`
	dotWithoutStart := `digraph { foo [result=""]; bar [result="x=1"]; foo -> bar [cond="true"]; }`
	invalidDOT := `digraph { start [result=]; }`

	t.Run("valid DOT returns graph with start", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()

		// Act
		graph, err := parser.Parse(context.Background(), validDOT)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, graph)
		assert.Equal(t, "start", graph.Start)
		_, ok := graph.Nodes["start"]
		assert.True(t, ok, "missing node start")
		_, ok = graph.Nodes["ok"]
		assert.True(t, ok, "missing node ok")
		assert.GreaterOrEqual(t, len(graph.Edges), 1, "expected at least one edge")
	})
	t.Run("DOT without start node returns ErrNoStartNode", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()

		// Act
		_, err := parser.Parse(context.Background(), dotWithoutStart)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNoStartNode)
	})
	t.Run("invalid DOT syntax returns error", func(t *testing.T) {
		// Arrange
		parser := NewDotParser()

		// Act
		_, err := parser.Parse(context.Background(), invalidDOT)

		// Assert
		assert.Error(t, err)
		assert.NotErrorIs(t, err, ErrNoStartNode)
	})
}
