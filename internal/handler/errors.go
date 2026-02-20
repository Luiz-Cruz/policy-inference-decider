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
}

func jsonErrorResponse(code int, errorCode, message string) events.APIGatewayProxyResponse {
	body := APIError{
		Error:   errorCode,
		Message: message,
		Status:  code,
	}
	b, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(b),
	}
}

func errorFromPolicy(err error) (code int, errorCode, message string) {
	if errors.Is(err, policy.ErrNoStartNode) {
		return http.StatusBadRequest, CodePolicyNoStartNode, err.Error()
	}
	return http.StatusInternalServerError, CodeInternalError, "An internal error occurred."
}

func errorFromParseDOT(err error) (code int, errorCode, message string) {
	if errors.Is(err, policy.ErrNoStartNode) {
		return http.StatusBadRequest, CodePolicyNoStartNode, err.Error()
	}
	return http.StatusBadRequest, CodeInvalidPolicyDOT, err.Error()
}

func errorFromBindJSON(err error) (code int, errorCode, message string) {
	return http.StatusBadRequest, CodeInvalidRequestBody, err.Error()
}

type ErrorMapper func(err error) (code int, errorCode, message string)

func Handle(err error, mapErr ErrorMapper) events.APIGatewayProxyResponse {
	code, errorCode, message := mapErr(err)
	return jsonErrorResponse(code, errorCode, message)
}
