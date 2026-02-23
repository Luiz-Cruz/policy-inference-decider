package handler

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"

	"policy-inference-decider/internal/apierror"
	"policy-inference-decider/internal/policy"
)

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func jsonErrorResponse(apiError apierror.APIError) events.APIGatewayProxyResponse {
	responseBody, _ := json.Marshal(apiError)
	return events.APIGatewayProxyResponse{
		StatusCode: apiError.Status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}
}

func errorFromPolicy(err error) apierror.APIError {
	if errors.Is(err, policy.ErrNoStartNode) {
		return apierror.NewNoStartNodeError()
	}
	if errors.Is(err, policy.ErrInvalidCondition) {
		return apierror.NewInvalidConditionError()
	}
	return apierror.NewInternalError()
}

func errorFromParseDOT(err error) apierror.APIError {
	if errors.Is(err, policy.ErrNoStartNode) {
		return apierror.NewNoStartNodeError()
	}
	return apierror.NewInvalidPolicyDotError()
}

func errorFromBindJSON(err error) apierror.APIError {
	return apierror.NewInvalidRequestBodyError()
}

type ErrorMapper func(err error) (apiError apierror.APIError)

func Handle(err error, mapErr ErrorMapper) events.APIGatewayProxyResponse {
	apiError := mapErr(err)
	return jsonErrorResponse(apiError)
}
