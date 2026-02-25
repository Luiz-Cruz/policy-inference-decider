package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"policy-inference-decider/internal/metric"
	"policy-inference-decider/internal/policy"
)

type Handler struct {
	policy policy.Inferrer
}

func NewInferHandler(policy policy.Inferrer) *Handler {
	return &Handler{policy: policy}
}

func (h *Handler) Infer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body policy.InferRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:bind_json] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		metric.IncrementError(ctx, metric.CauseBindJSON)
		return Handle(err, errorFromBindJSON), nil
	}

	graph, err := h.policy.Parse(ctx, body.PolicyDOT)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:parse_dot] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		metric.IncrementError(ctx, metric.CauseParseDOT)
		return Handle(err, errorFromParseDOT), nil
	}

	resp, err := h.policy.Process(ctx, graph, body.Input)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:execute] [request_id:%s] [err:%+v] ", req.RequestContext.RequestID, err))
		metric.IncrementError(ctx, metric.CauseExecute)
		return Handle(err, errorFromPolicy), nil
	}

	metric.IncrementSuccess(ctx)
	responseBody, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}
