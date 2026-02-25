package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"policy-inference-decider/internal/apierror"
	"policy-inference-decider/internal/policy"
)

type Handler struct {
	parser   policy.Parser
	executor policy.Executor
}

func NewInferHandler(parser policy.Parser, executor policy.Executor) *Handler {
	return &Handler{parser: parser, executor: executor}
}

func (h *Handler) Infer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := strings.TrimSuffix(req.Path, "/")
	switch path {
	case "/infer":
		if req.HTTPMethod != http.MethodPost {
			return jsonErrorResponse(apierror.NewMethodNotAllowedError()), nil
		}
		return h.infer(ctx, req)
	case "/ping":
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "pong",
		}, nil
	default:
		return jsonErrorResponse(apierror.NewNotFoundError()), nil
	}
}

func (h *Handler) infer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
