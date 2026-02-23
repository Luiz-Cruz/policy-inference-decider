package policy

import "context"

type (
	Executor interface {
		Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error)
	}
	Parser interface {
		Parse(ctx context.Context, dot string) (*Graph, error)
	}
	PolicyInferrer interface {
		Parse(ctx context.Context, dot string) (*Graph, error)
		Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error)
	}
	policyExecutor struct {
		executor Executor
		parser   Parser
	}
)

func (p *policyExecutor) Parse(ctx context.Context, dot string) (*Graph, error) {
	return p.parser.Parse(ctx, dot)
}

func (p *policyExecutor) Process(ctx context.Context, graph *Graph, input map[string]any) (InferResponse, error) {
	return p.executor.Process(ctx, graph, input)
}

func NewPolicyExecutor(exec Executor, parser Parser) PolicyInferrer {
	return &policyExecutor{
		executor: exec,
		parser:   parser,
	}
}
