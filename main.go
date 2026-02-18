package main

import (
	"policy-inference-decider/internal/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Infer)
}
