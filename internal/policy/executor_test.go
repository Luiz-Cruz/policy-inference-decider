package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type executeScenario struct {
	t     *testing.T
	dot   string
	vars  map[string]any
	graph *Graph
	out   map[string]any
	err   error
}

func TestExecute(t *testing.T) {
	linearDOT := `digraph { start [result=""]; ok [result="approved=true"]; no [result="approved=false"]; start -> ok [cond="age>=18"]; start -> no [cond="age<18"]; }`
	cycleDOT := `digraph { start [result=""]; a [result="done=true"]; start -> a [cond="x==1"]; a -> a [cond="x==1"]; }`

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"single path applies node results": {run: func(t *testing.T) {
			s := startExecuteScenario(t)
			s.givenDOTAndVars(linearDOT, map[string]any{"age": 20})
			s.whenExecuteIsExecuted()
			s.thenOutputHasAgeAndApproved(20, true)
		}},
		"single path age under 18": {run: func(t *testing.T) {
			s := startExecuteScenario(t)
			s.givenDOTAndVars(linearDOT, map[string]any{"age": 15})
			s.whenExecuteIsExecuted()
			s.thenOutputApprovedIs(false)
		}},
		"graph with cycle terminates and returns": {run: func(t *testing.T) {
			s := startExecuteScenario(t)
			s.givenDOTAndVars(cycleDOT, map[string]any{"x": 1.0})
			s.whenExecuteIsExecuted()
			s.thenOutputHasXAndDone(1.0, true)
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startExecuteScenario(t *testing.T) *executeScenario {
	return &executeScenario{t: t}
}

func (s *executeScenario) givenDOTAndVars(dot string, vars map[string]any) {
	s.dot = dot
	s.vars = vars
}

func (s *executeScenario) whenExecuteIsExecuted() {
	var err error
	s.graph, err = DotParser{}.Parse(context.Background(), s.dot)
	if err != nil {
		s.err = err
		return
	}
	var resp InferResponse
	resp, s.err = GraphExecutor{}.Process(context.Background(), s.graph, s.vars)
	s.out = resp.Output
}

func (s *executeScenario) thenOutputHasAgeAndApproved(age int, approved bool) {
	assert.NoError(s.t, s.err)
	assert.Equal(s.t, age, s.out["age"])
	assert.Equal(s.t, approved, s.out["approved"])
}

func (s *executeScenario) thenOutputApprovedIs(approved bool) {
	assert.NoError(s.t, s.err)
	assert.Equal(s.t, approved, s.out["approved"])
}

func (s *executeScenario) thenOutputHasXAndDone(x float64, done bool) {
	assert.NoError(s.t, s.err)
	assert.Equal(s.t, x, s.out["x"])
	assert.Equal(s.t, done, s.out["done"])
}
