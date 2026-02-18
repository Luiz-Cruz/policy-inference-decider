package policy

func Execute(g *Graph, input map[string]any) (map[string]any, error) {
	out := copyInputToOutput(input)
	visited := make(map[string]bool)
	current := g.Start
	for {
		applyNodeResult(g.Nodes[current], out)
		visited[current] = true
		next := findNextNode(current, g, out)
		if next == "" || visited[next] {
			break
		}
		current = next
	}
	return out, nil
}

func copyInputToOutput(input map[string]any) map[string]any {
	out := make(map[string]any, len(input))
	for k, v := range input {
		out[k] = v
	}
	return out
}

func applyNodeResult(node *Node, out map[string]any) {
	if node == nil {
		return
	}
	_ = ApplyResult(node.Result, out)
}

func findNextNode(current string, g *Graph, vars map[string]any) string {
	for _, e := range g.Edges {
		if e.From != current {
			continue
		}
		ok, err := EvalCondition(e.Cond, vars)
		if err != nil || !ok {
			continue
		}
		return e.To
	}
	return ""
}
