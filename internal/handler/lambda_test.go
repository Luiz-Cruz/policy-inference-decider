package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestInfer(t *testing.T) {
	ctx := context.Background()
	req := events.APIGatewayProxyRequest{
		Body: `{"policy_dot":"","input":{}}`,
	}
	resp, err := Infer(ctx, req)
	if err != nil {
		t.Fatalf("Infer: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != "Lets Go!" {
		t.Errorf("Body = %q; want %q", resp.Body, "Lets Go!")
	}
}
