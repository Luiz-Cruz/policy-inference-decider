package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"policy-inference-decider/internal/handler"
	"policy-inference-decider/internal/policy"
)

func main() {
	cache := policy.NewLFUGraphCache(policy.DefaultCacheCapacity)
	parser := policy.NewCachedParser(policy.NewDotParser(), cache)
	executor := policy.NewGraphExecutor()
	inferHandler := handler.NewInferHandler(parser, executor)
	lambda.Start(inferHandler.Infer)
}
