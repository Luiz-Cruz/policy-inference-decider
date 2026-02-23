package main

import (
	"policy-inference-decider/internal/handler"
	"policy-inference-decider/internal/policy"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	policyExecutor := policy.NewPolicyExecutor(&policy.GraphExecutor{})
	inferHandler := handler.New(policyExecutor)
	lambda.Start(inferHandler.Infer)
}
