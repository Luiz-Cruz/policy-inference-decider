package policy

import (
	"regexp"
	"strconv"
	"strings"
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

func EvalEdgeCondition(edge *Edge, vars map[string]any) (bool, error) {
	if edge.Cond == "" {
		return true, nil
	}
	if !edge.ValidCond || edge.CompiledCond == nil {
		return false, ErrInvalidCondition
	}
	result, err := edge.CompiledCond.Evaluate(vars)
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

func preParseResult(result string) map[string]any {
	result = strings.TrimSpace(result)
	if result == "" {
		return nil
	}
	parsed := make(map[string]any)
	pairs := strings.Split(result, ",")
	for _, pair := range pairs {
		key, valStr, ok := parseKeyValue(pair)
		if !ok {
			continue
		}
		parsed[key] = parseResultValue(valStr)
	}
	return parsed
}

func applyParsedResult(parsed map[string]any, vars map[string]any) {
	for k, v := range parsed {
		vars[k] = v
	}
}
