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

type Handler struct {
	parser   policy.Parser
	executor policy.Executor
}

func NewInferHandler(parser policy.Parser, executor policy.Executor) *Handler {
	return &Handler{parser: parser, executor: executor}
}

func (h *Handler) Infer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body policy.InferRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:bind_json] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		return Handle(err, errorFromBindJSON), nil
	}

	graph, err := h.parser.Parse(ctx, body.PolicyDOT)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:parse_dot] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		return Handle(err, errorFromParseDOT), nil
	}

	resp, err := h.executor.Process(ctx, graph, body.Input)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:execute] [request_id:%s] [err:%+v] ", req.RequestContext.RequestID, err))
		return Handle(err, errorFromPolicy), nil
	}

	responseBody, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}
