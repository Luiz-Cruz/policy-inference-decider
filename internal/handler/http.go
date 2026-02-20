package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"policy-inference-decider/internal/policy"

	"github.com/aws/aws-lambda-go/events"
)

func Infer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body policy.InferRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:bind_json] [err:%+v]", err))
		return Handle(err, errorFromBindJSON), nil
	}
	graph, err := policy.ParseDOT(body.PolicyDOT)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:parse_dot] [err:%+v]", err))
		return Handle(err, errorFromParseDOT), nil
	}
	out, err := policy.Execute(graph, body.Input)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:execute] [err:%+v]", err))
		return Handle(err, errorFromPolicy), nil
	}
	resp := policy.InferResponse{Output: out}
	b, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(b),
	}, nil
}
