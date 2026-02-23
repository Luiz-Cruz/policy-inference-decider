package policy

import "context"

type (
	Executor interface {
		Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error)
	}
	PolicyExecutor struct {
		executor Executor
	}
)

func NewPolicyExecutor(exec Executor) *PolicyExecutor {
	return &PolicyExecutor{executor: exec}
}

func (e *PolicyExecutor) Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error) {
	return e.executor.Process(ctx, graph, input)
}
