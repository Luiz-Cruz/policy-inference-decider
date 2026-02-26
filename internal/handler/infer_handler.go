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
	"policy-inference-decider/internal/metric"
	"policy-inference-decider/internal/policy"
)

type Handler struct {
	parser   policy.Parser
	executor policy.Executor
}

func NewInferHandler(parser policy.Parser, executor policy.Executor) *Handler {
	return &Handler{parser: parser, executor: executor}
}

func (h *Handler) Infer(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	switch req.RequestContext.HTTP.Method {
	case http.MethodGet:
		return h.handleGet(req), nil
	case http.MethodPost:
		return h.handlePost(ctx, req), nil
	default:
		return jsonErrorResponseURL(apierror.NewMethodNotAllowedError()), nil
	}
}

func (h *Handler) handleGet(req events.LambdaFunctionURLRequest) events.LambdaFunctionURLResponse {
	path := strings.TrimSuffix(pathFromRequest(req), "/")
	switch path {
	case "/ping":
		return events.LambdaFunctionURLResponse{
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "pong",
		}
	case "/infer":
		return jsonErrorResponseURL(apierror.NewMethodNotAllowedError())
	default:
		return jsonErrorResponseURL(apierror.NewNotFoundError())
	}
}

func (h *Handler) handlePost(ctx context.Context, req events.LambdaFunctionURLRequest) events.LambdaFunctionURLResponse {
	path := strings.TrimSuffix(pathFromRequest(req), "/")
	if path != "/infer" {
		return jsonErrorResponseURL(apierror.NewNotFoundError())
	}
	return h.infer(ctx, req)
}

func pathFromRequest(req events.LambdaFunctionURLRequest) string {
	if path := req.RequestContext.HTTP.Path; path != "" {
		return path
	}
	return req.RawPath
}

func (h *Handler) infer(ctx context.Context, req events.LambdaFunctionURLRequest) events.LambdaFunctionURLResponse {
	var body policy.InferRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:bind_json] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		return HandleURL(err, errorFromBindJSON)
	}

	graph, err := h.parser.Parse(ctx, body.PolicyDOT)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:parse_dot] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		return HandleURL(err, errorFromParseDOT)
	}

	resp, err := h.executor.Process(ctx, graph, body.Input)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:execute] [request_id:%s] [err:%+v] ", req.RequestContext.RequestID, err))
		return HandleURL(err, errorFromPolicy)
	}

	metric.EmitSuccess()
	responseBody, _ := json.Marshal(resp)
	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}
}
