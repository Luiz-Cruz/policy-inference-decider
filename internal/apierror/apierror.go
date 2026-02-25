package apierror

import "net/http"

const (
	CodeInvalidRequestBody = "invalid_request_body"
	CodeInvalidPolicyDOT   = "invalid_policy_dot"
	CodePolicyNoStartNode  = "policy_no_start_node"
	CodeInvalidCondition   = "invalid_condition"
	CodeInternalError      = "internal_error"
	CodeNotFound           = "not_found"
	CodeMethodNotAllowed   = "method_not_allowed"
)

const (
	msgInvalidRequestBody = "Invalid request body."
	msgInvalidPolicyDOT   = "Invalid policy DOT format."
	msgPolicyNoStartNode  = "Policy graph has no start node."
	msgInvalidCondition   = "Invalid condition in policy."
	msgInternalError      = "An internal error occurred."
	msgNotFound           = "Not found."
	msgMethodNotAllowed   = "Method not allowed."
)

type APIError struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"error"`
	Message   string `json:"message"`
}

func NewInvalidRequestBodyError() APIError {
	return APIError{Status: http.StatusBadRequest, ErrorCode: CodeInvalidRequestBody, Message: msgInvalidRequestBody}
}

func NewInvalidPolicyDotError() APIError {
	return APIError{Status: http.StatusBadRequest, ErrorCode: CodeInvalidPolicyDOT, Message: msgInvalidPolicyDOT}
}

func NewNoStartNodeError() APIError {
	return APIError{Status: http.StatusBadRequest, ErrorCode: CodePolicyNoStartNode, Message: msgPolicyNoStartNode}
}

func NewInvalidConditionError() APIError {
	return APIError{Status: http.StatusBadRequest, ErrorCode: CodeInvalidCondition, Message: msgInvalidCondition}
}

func NewInternalError() APIError {
	return APIError{Status: http.StatusInternalServerError, ErrorCode: CodeInternalError, Message: msgInternalError}
}

func NewNotFoundError() APIError {
	return APIError{Status: http.StatusNotFound, ErrorCode: CodeNotFound, Message: msgNotFound}
}

func NewMethodNotAllowedError() APIError {
	return APIError{Status: http.StatusMethodNotAllowed, ErrorCode: CodeMethodNotAllowed, Message: msgMethodNotAllowed}
}
