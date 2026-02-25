package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"policy-inference-decider/internal/handler"
	"policy-inference-decider/internal/policy"
)

func main() {
	parser := &policy.DotParser{}
	executor := &policy.GraphExecutor{}
	inferHandler := handler.NewInferHandler(parser, executor)
	lambda.Start(inferHandler.Infer)
}
