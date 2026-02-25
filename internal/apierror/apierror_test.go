package apierror

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvalidRequestBodyError(t *testing.T) {
	e := NewInvalidRequestBodyError()
	assert.Equal(t, http.StatusBadRequest, e.Status)
	assert.Equal(t, CodeInvalidRequestBody, e.ErrorCode)
	assert.Equal(t, "Invalid request body.", e.Message)
}

func TestNewInvalidPolicyDotError(t *testing.T) {
	e := NewInvalidPolicyDotError()
	assert.Equal(t, http.StatusBadRequest, e.Status)
	assert.Equal(t, CodeInvalidPolicyDOT, e.ErrorCode)
	assert.Equal(t, "Invalid policy DOT format.", e.Message)
}

func TestNewNoStartNodeError(t *testing.T) {
	e := NewNoStartNodeError()
	assert.Equal(t, http.StatusBadRequest, e.Status)
	assert.Equal(t, CodePolicyNoStartNode, e.ErrorCode)
	assert.Equal(t, "Policy graph has no start node.", e.Message)
}

func TestNewInvalidConditionError(t *testing.T) {
	e := NewInvalidConditionError()
	assert.Equal(t, http.StatusBadRequest, e.Status)
	assert.Equal(t, CodeInvalidCondition, e.ErrorCode)
	assert.Equal(t, "Invalid condition in policy.", e.Message)
}

func TestNewInternalError(t *testing.T) {
	e := NewInternalError()
	assert.Equal(t, http.StatusInternalServerError, e.Status)
	assert.Equal(t, CodeInternalError, e.ErrorCode)
	assert.Equal(t, "An internal error occurred.", e.Message)
}

func TestNewNotFoundError(t *testing.T) {
	e := NewNotFoundError()
	assert.Equal(t, http.StatusNotFound, e.Status)
	assert.Equal(t, CodeNotFound, e.ErrorCode)
	assert.Equal(t, "Not found.", e.Message)
}

func TestNewMethodNotAllowedError(t *testing.T) {
	e := NewMethodNotAllowedError()
	assert.Equal(t, http.StatusMethodNotAllowed, e.Status)
	assert.Equal(t, CodeMethodNotAllowed, e.ErrorCode)
	assert.Equal(t, "Method not allowed.", e.Message)
}
