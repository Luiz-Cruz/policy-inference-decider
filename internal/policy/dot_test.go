package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type parseDOTScenario struct {
	t     *testing.T
	dot   string
	graph *Graph
	err   error
}

func TestParseDOT(t *testing.T) {
	validDOT := `digraph { start [result=""]; ok [result="approved=true"]; start -> ok [cond="age>=18"]; }`
	dotWithoutStart := `digraph { foo [result=""]; bar [result="x=1"]; foo -> bar [cond="true"]; }`
	invalidDOT := `digraph { start [result=]; }`

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"valid DOT returns graph with start": {run: func(t *testing.T) {
			s := startParseDOTScenario(t)
			s.givenDOT(validDOT)
			s.whenParseDOTIsExecuted()
			s.thenGraphHasStartAndNodes()
		}},
		"DOT without start node returns ErrNoStartNode": {run: func(t *testing.T) {
			s := startParseDOTScenario(t)
			s.givenDOT(dotWithoutStart)
			s.whenParseDOTIsExecuted()
			s.thenErrorIsErrNoStartNode()
		}},
		"invalid DOT syntax returns error": {run: func(t *testing.T) {
			s := startParseDOTScenario(t)
			s.givenDOT(invalidDOT)
			s.whenParseDOTIsExecuted()
			s.thenErrorIsNotErrNoStartNode()
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startParseDOTScenario(t *testing.T) *parseDOTScenario {
	return &parseDOTScenario{t: t}
}

func (s *parseDOTScenario) givenDOT(dot string) {
	s.dot = dot
}

func (s *parseDOTScenario) whenParseDOTIsExecuted() {
	s.graph, s.err = ParseDOT(s.dot)
}

func (s *parseDOTScenario) thenGraphHasStartAndNodes() {
	assert.NoError(s.t, s.err)
	assert.NotNil(s.t, s.graph)
	assert.Equal(s.t, "start", s.graph.Start)
	_, ok := s.graph.Nodes["start"]
	assert.True(s.t, ok, "missing node start")
	_, ok = s.graph.Nodes["ok"]
	assert.True(s.t, ok, "missing node ok")
	assert.GreaterOrEqual(s.t, len(s.graph.Edges), 1, "expected at least one edge")
}

func (s *parseDOTScenario) thenErrorIsErrNoStartNode() {
	assert.Error(s.t, s.err)
	assert.ErrorIs(s.t, s.err, ErrNoStartNode)
}

func (s *parseDOTScenario) thenErrorIsNotErrNoStartNode() {
	assert.Error(s.t, s.err)
	assert.NotErrorIs(s.t, s.err, ErrNoStartNode)
}
