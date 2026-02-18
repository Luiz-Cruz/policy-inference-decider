package policy

type (
	InferRequest struct {
		PolicyDOT string         `json:"policy_dot"`
		Input     map[string]any `json:"input"`
	}

	InferResponse struct {
		Output map[string]any `json:"output"`
	}

	Graph struct {
		Nodes map[string]*Node
		Edges []*Edge
		Start string
	}

	Node struct {
		ID     string
		Result string
	}

	Edge struct {
		From string
		To   string
		Cond string
	}
)
