package policy

import (
	"strconv"
	"strings"

	"github.com/casbin/govaluate"
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

func parseKeyValue(pair string) (key, value string, ok bool) {
	kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
	if len(kv) != 2 {
		return "", "", false
	}
	key = strings.TrimSpace(kv[0])
	value = strings.TrimSpace(kv[1])
	value = strings.Trim(value, "\"")
	return key, value, true
}

func parseResultValue(valStr string) any {
	if v, err := strconv.ParseBool(valStr); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(valStr, 64); err == nil {
		return v
	}
	return valStr
}

func ApplyResult(result string, vars map[string]any) {
	result = strings.TrimSpace(result)
	if result == "" {
		return
	}
	pairs := strings.Split(result, ",")
	for _, pair := range pairs {
		key, valStr, ok := parseKeyValue(pair)
		if !ok {
			continue
		}
		vars[key] = parseResultValue(valStr)
	}
}
