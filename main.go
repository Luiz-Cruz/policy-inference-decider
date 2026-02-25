package main

import (
	"policy-inference-decider/internal/handler"
	"policy-inference-decider/internal/policy"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	parser := &policy.DotParser{}
	executor := &policy.GraphExecutor{}
	inferHandler := handler.NewInferHandler(parser, executor)
	lambda.Start(inferHandler.Infer)
}
