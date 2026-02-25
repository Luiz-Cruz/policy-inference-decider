package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

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
		return Handle(err, errorFromBindJSON), nil
	}

	graph, err := h.policy.Parse(ctx, body.PolicyDOT)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:parse_dot] [request_id: %s] [err:%+v]", req.RequestContext.RequestID, err))
		return Handle(err, errorFromParseDOT), nil
	}

	resp, err := h.policy.Process(ctx, graph, body.Input)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:execute] [request_id:%s] [err:%+v] ", req.RequestContext.RequestID, err))
		return Handle(err, errorFromPolicy), nil
	}

	go func() {
		bg := context.Background()
		cfg, err := config.LoadDefaultConfig(bg)
		if err != nil {
			slog.ErrorContext(bg, "cloudwatch config", "error", err)
			return
		}
		client := cloudwatch.NewFromConfig(cfg)
		_, err = client.PutMetricData(bg, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String("PolicyInferenceDecider"),
			MetricData: []types.MetricDatum{{
				MetricName: aws.String("policy_inference"),
				Value:      aws.Float64(1),
				Unit:       types.StandardUnitCount,
				Timestamp:  aws.Time(time.Now().UTC()),
				Dimensions: []types.Dimension{
					{Name: aws.String("result"), Value: aws.String("success")},
				},
			}},
		})
		if err != nil {
			slog.ErrorContext(bg, "PutMetricData", "error", err)
		}
	}()

	responseBody, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}
