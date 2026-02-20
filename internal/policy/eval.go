package policy

import (
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

func EvalCondition(cond string, vars map[string]any) (bool, error) {
	if cond == "" {
		return true, nil
	}
	expr, err := govaluate.NewEvaluableExpression(cond)
	if err != nil {
		return false, err
	}
	evars := make(map[string]interface{})
	for k, v := range vars {
		evars[k] = v
	}
	result, err := expr.Evaluate(evars)
	if err != nil {
		return false, err
	}
	b, ok := result.(bool)
	if !ok {
		return false, nil
	}
	return b, nil
}

func ApplyResult(result string, vars map[string]any) error {
	result = strings.TrimSpace(result)
	if result == "" {
		return nil
	}
	pairs := strings.Split(result, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		valStr := strings.TrimSpace(kv[1])
		valStr = strings.Trim(valStr, "\"")
		var val any
		if v, err := strconv.ParseBool(valStr); err == nil {
			val = v
		} else if v, err := strconv.ParseFloat(valStr, 64); err == nil {
			val = v
		} else {
			val = valStr
		}
		vars[key] = val
	}
	return nil
}
