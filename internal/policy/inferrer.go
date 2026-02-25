package policy

import "context"

type (
	Executor interface {
		Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error)
	}
	Parser interface {
		Parse(ctx context.Context, dot string) (*Graph, error)
	}
)
