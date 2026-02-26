package policy

import "context"

type GraphExecutor struct{}

func NewGraphExecutor() *GraphExecutor {
	return &GraphExecutor{}
}

func (GraphExecutor) Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error) {
	out := copyInputToOutput(input)
	visited := make(map[string]bool)
	current := graph.Start
	for {
		if node := graph.Nodes[current]; node != nil {
			ApplyResult(node.Result, out)
		}
		visited[current] = true
		next, err := findNextNode(current, graph, out)
		if err != nil {
			return InferResponse{}, err
		}
		if next == "" || visited[next] {
			break
		}
		current = next
	}
	return InferResponse{Output: out}, nil
}

func copyInputToOutput(input map[string]any) map[string]any {
	out := make(map[string]any, len(input))
	for k, v := range input {
		out[k] = v
	}
	return out
}

// findNextNode returns the first outgoing edge from current whose condition evaluates to true (deterministic single path).
func findNextNode(current string, graph *Graph, vars map[string]any) (string, error) {
	for _, edge := range graph.Edges {
		if edge.From != current {
			continue
		}
		ok, err := EvalCondition(edge.Cond, vars)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}
		return edge.To, nil
	}
	return "", nil
}
