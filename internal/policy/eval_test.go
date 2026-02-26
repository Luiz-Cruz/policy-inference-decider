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
		"bare number is invalid condition": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("1", map[string]any{})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"arithmetic plus returns invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age+1>=18", map[string]any{"age": 17})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"arithmetic star returns invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("score*2>100", map[string]any{"score": 60})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"allowed && and ||": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age>=18 && score>700", map[string]any{"age": 25, "score": 720})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"allowed comparison with negative number": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("score> -1", map[string]any{"score": 100})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"allowed string literal": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars(`role=="admin"`, map[string]any{"role": "admin"})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"allowed boolean literal": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("active==true", map[string]any{"active": true})
			s.whenEvalConditionIsExecuted()
			s.thenResultIs(true)
		}},
		"arithmetic plus with spaces invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age + 1 >= 18", map[string]any{"age": 17})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"arithmetic star invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("score*2 == 1400", map[string]any{"score": 700})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"arithmetic slash invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("points / 2 > 100", map[string]any{"points": 250})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"arithmetic minus invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("balance - 50 > 0", map[string]any{"balance": 100})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"triple equals invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars(`status === "ok"`, map[string]any{"status": "ok"})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"angle bracket not equals invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age <> 18", map[string]any{"age": 20})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"arrow equals invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("score => 800", map[string]any{"score": 800})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"caret equals invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("active ^= true", map[string]any{"active": true})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"parentheses not supported": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age>=18 && (score>700)", map[string]any{"age": 25, "score": 720})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"parentheses around comparisons invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars(`(name == "João") || (age < 20)`, map[string]any{"name": "João", "age": 18})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"function call len invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("len(name) > 5", map[string]any{"name": "Alice"})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"function call isAdult invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("isAdult(age) == true", map[string]any{"age": 20})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"NULL value invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("height == NULL", map[string]any{})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"undefined value invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("name == undefined", map[string]any{})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"single ampersand invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age>=18 & score>700", map[string]any{"age": 25, "score": 720})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"triple pipe invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("age>=18 ||| score>700", map[string]any{"age": 25, "score": 720})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"string without quotes invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("name == Joao", map[string]any{"name": "Joao"})
			s.whenEvalConditionIsExecuted()
			s.thenErrorAndResultIsFalse()
		}},
		"identifier starting with digit invalid": {run: func(t *testing.T) {
			s := startEvalConditionScenario(t)
			s.givenCondAndVars("1age == 10", map[string]any{"1age": 10})
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
		"malformed pair without equals is skipped": {run: func(t *testing.T) {
			s := startApplyResultScenario(t)
			s.givenResultAndVars("a=1, badpair, b=2", map[string]any{})
			s.whenApplyResultIsExecuted()
			s.thenVarEquals("a", true)
			s.thenVarEquals("b", 2.0)
			assert.Len(s.t, s.vars, 2)
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
