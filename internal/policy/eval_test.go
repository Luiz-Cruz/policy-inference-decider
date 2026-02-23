package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type evalConditionScenario struct {
	t    *testing.T
	cond string
	vars map[string]any
	got  bool
	err  error
}

func TestEvalCondition(t *testing.T) {
	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"empty cond is true": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("", map[string]any{"x": 1})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"age>=18 true": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age>=18", map[string]any{"age": 20})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"age>=18 false": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age>=18", map[string]any{"age": 15})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(false)
		}},
		"x==1 true": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("x==1", map[string]any{"x": 1.0})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"invalid expr returns false and error": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("invalid!!!", map[string]any{})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startEvalConditionScenario(t *testing.T) *evalConditionScenario {
	return &evalConditionScenario{t: t}
}

func (s *evalConditionScenario) givenCondAndVars(cond string, vars map[string]any) {
	s.cond = cond
	s.vars = vars
}

func (s *evalConditionScenario) whenEvalConditionIsExecuted() {
	s.got, s.err = EvalCondition(s.cond, s.vars)
}

func (s *evalConditionScenario) thenResultIs(expected bool) {
	assert.NoError(s.t, s.err)
	assert.Equal(s.t, expected, s.got)
}

func (s *evalConditionScenario) thenErrorAndResultIsFalse() {
	assert.Error(s.t, s.err)
	assert.False(s.t, s.got)
}

type applyResultScenario struct {
	t      *testing.T
	result string
	vars   map[string]any
	err    error
}

func TestApplyResult(t *testing.T) {
	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"empty result does nothing": {run: func(t *testing.T) {
			s := startApplyResultScenario(t)
			s.givenResultAndVars("", map[string]any{"a": 1})
			s.whenApplyResultIsExecuted()
			s.thenVarsUnchanged()
		}},
		"key=value sets string": {run: func(t *testing.T) {
			s := startApplyResultScenario(t)
			s.givenResultAndVars("name=foo", map[string]any{})
			s.whenApplyResultIsExecuted()
			s.thenVarEquals("name", "foo")
		}},
		"key=true sets bool": {run: func(t *testing.T) {
			s := startApplyResultScenario(t)
			s.givenResultAndVars("approved=true", map[string]any{})
			s.whenApplyResultIsExecuted()
			s.thenVarEquals("approved", true)
		}},
		"multiple pairs": {run: func(t *testing.T) {
			s := startApplyResultScenario(t)
			s.givenResultAndVars("num=2.5, flag=true", map[string]any{})
			s.whenApplyResultIsExecuted()
			s.thenVarEquals("num", 2.5)
			s.thenVarEquals("flag", true)
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startApplyResultScenario(t *testing.T) *applyResultScenario {
	return &applyResultScenario{t: t}
}

func (s *applyResultScenario) givenResultAndVars(result string, vars map[string]any) {
	s.result = result
	s.vars = vars
}

func (s *applyResultScenario) whenApplyResultIsExecuted() {
	ApplyResult(s.result, s.vars)
}

func (s *applyResultScenario) thenVarsUnchanged() {
	assert.NoError(s.t, s.err)
	assert.Len(s.t, s.vars, 1)
	assert.Equal(s.t, 1, s.vars["a"])
}

func (s *applyResultScenario) thenVarEquals(key string, expected any) {
	assert.NoError(s.t, s.err)
	assert.Equal(s.t, expected, s.vars[key])
}
