package metric

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	Namespace  = "PolicyInferenceDecider"
	MetricName = "policy_inference"

	CauseBindJSON = "bind_json"
	CauseParseDOT = "parse_dot"
	CauseExecute  = "execute"
)

var (
	clientOnce sync.Once
	cwClient   *cloudwatch.Client
	clientErr  error
)

func client(ctx context.Context) (*cloudwatch.Client, error) {
	clientOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			clientErr = err
			return
		}
		cwClient = cloudwatch.NewFromConfig(cfg)
	})
	return cwClient, clientErr
}

func IncrementSuccess(ctx context.Context) {
	go send(ctx, "success", "")
}

func IncrementError(ctx context.Context, cause string) {
	go send(ctx, "error", cause)
}

func send(ctx context.Context, result, cause string) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "metric panic", "error", r)
		}
	}()
	ctx = context.Background()
	c, err := client(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "cloudwatch client", "error", err)
		return
	}
	dims := []types.Dimension{
		{Name: aws.String("result"), Value: aws.String(result)},
	}
	if cause != "" {
		dims = append(dims, types.Dimension{Name: aws.String("cause"), Value: aws.String(cause)})
	}
	_, err = c.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(Namespace),
		MetricData: []types.MetricDatum{{
			MetricName: aws.String(MetricName),
			Value:      aws.Float64(1),
			Unit:       types.StandardUnitCount,
			Timestamp:  aws.Time(time.Now().UTC()),
			Dimensions: dims,
		}},
	})
	if err != nil {
		slog.ErrorContext(ctx, "PutMetricData", "error", err)
	}
}
