package policy

import "github.com/casbin/govaluate"

const StartNodeID = "start"

type (
	InferRequest struct {
		PolicyDOT string         `json:"policy_dot"`
		Input     map[string]any `json:"input"`
	}

	InferResponse struct {
		Output map[string]any `json:"output"`
	}

	Graph struct {
		Nodes   map[string]*Node
		Edges   []*Edge
		AdjList map[string][]*Edge
		Start   string
	}

	Node struct {
		ID           string
		Result       string
		ParsedResult map[string]any
	}

	Edge struct {
		From         string
		To           string
		Cond         string
		ValidCond    bool
		CompiledCond *govaluate.EvaluableExpression
	}
)
