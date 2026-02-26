package policy

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/casbin/govaluate"
)

const arithmeticCondChars = "+*/"

var validCondRegex = regexp.MustCompile(
	`^\s*[a-zA-Z_]\w*\s*(==|!=|>=|<=|>|<)\s*(-?\d+(?:\.\d+)?|"[^"]*"|true|false)(\s*(?:&&|\|\|)\s*[a-zA-Z_]\w*\s*(==|!=|>=|<=|>|<)\s*(-?\d+(?:\.\d+)?|"[^"]*"|true|false))*\s*$`,
)

func isValidCond(cond string) bool {
	if strings.ContainsAny(cond, arithmeticCondChars) {
		return false
	}
	return validCondRegex.MatchString(cond)
}

func EvalCondition(cond string, vars map[string]any) (bool, error) {
	if cond == "" {
		return true, nil
	}
	if !isValidCond(cond) {
		return false, ErrInvalidCondition
	}
	expr, err := govaluate.NewEvaluableExpression(cond)
	if err != nil {
		return false, ErrInvalidCondition
	}
	evars := make(map[string]interface{})
	for k, v := range vars {
		evars[k] = v
	}
	result, err := expr.Evaluate(evars)
	if err != nil {
		return false, ErrInvalidCondition
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
	if value, err := strconv.ParseBool(valStr); err == nil {
		return value
	}
	if value, err := strconv.ParseFloat(valStr, 64); err == nil {
		return value
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
