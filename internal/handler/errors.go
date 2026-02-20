package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"policy-inference-decider/internal/policy"
)

const (
	CodeInvalidRequestBody = "invalid_request_body"
	CodeInvalidPolicyDOT   = "invalid_policy_dot"
	CodePolicyNoStartNode  = "policy_no_start_node"
	CodeInternalError      = "internal_error"
)

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Cause   []any  `json:"cause,omitempty"`
}

func jsonErrorResponse(code int, errorCode, message string, cause []any) events.APIGatewayProxyResponse {
	body := APIError{
		Error:   errorCode,
		Message: message,
		Status:  code,
		Cause:   cause,
	}
	b, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(b),
	}
}

func errorFromPolicy(err error) (code int, errorCode, message string, cause []any) {
	if errors.Is(err, policy.ErrNoStartNode) {
		return http.StatusBadRequest, CodePolicyNoStartNode, err.Error(), nil
	}
	return http.StatusInternalServerError, CodeInternalError, "An internal error occurred.", nil
}

func errorFromParseDOT(err error) (code int, errorCode, message string, cause []any) {
	if errors.Is(err, policy.ErrNoStartNode) {
		return http.StatusBadRequest, CodePolicyNoStartNode, err.Error(), nil
	}

	if errors.Is(err, policy.ErrInvalidPolicyDot) {
		return http.StatusBadRequest, CodeInvalidPolicyDOT, err.Error(), nil
	}

	return http.StatusBadRequest, CodeInvalidPolicyDOT, err.Error(), nil
}

func errorFromBindJSON(err error) (code int, errorCode, message string, cause []any) {
	return http.StatusBadRequest, CodeInvalidRequestBody, err.Error(), nil
}

type ErrorMapper func(err error) (code int, errorCode, message string, cause []any)

func Handle(err error, mapErr ErrorMapper) events.APIGatewayProxyResponse {
	code, errorCode, message, cause := mapErr(err)
	return jsonErrorResponse(code, errorCode, message, cause)
}
