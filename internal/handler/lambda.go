package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Infer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(req.Body)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Lets Go!",
	}, nil
}
