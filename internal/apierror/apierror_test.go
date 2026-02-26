package apierror

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvalidRequestBodyError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewInvalidRequestBodyError()

		// Assert
		assert.Equal(t, http.StatusBadRequest, e.Status)
		assert.Equal(t, CodeInvalidRequestBody, e.ErrorCode)
		assert.Equal(t, "Invalid request body.", e.Message)
	})
}

func TestNewInvalidPolicyDotError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewInvalidPolicyDotError()

		// Assert
		assert.Equal(t, http.StatusBadRequest, e.Status)
		assert.Equal(t, CodeInvalidPolicyDOT, e.ErrorCode)
		assert.Equal(t, "Invalid policy DOT format.", e.Message)
	})
}

func TestNewNoStartNodeError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewNoStartNodeError()

		// Assert
		assert.Equal(t, http.StatusBadRequest, e.Status)
		assert.Equal(t, CodePolicyNoStartNode, e.ErrorCode)
		assert.Equal(t, "Policy graph has no start node.", e.Message)
	})
}

func TestNewInvalidConditionError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewInvalidConditionError()

		// Assert
		assert.Equal(t, http.StatusBadRequest, e.Status)
		assert.Equal(t, CodeInvalidCondition, e.ErrorCode)
		assert.Equal(t, "Invalid condition in policy.", e.Message)
	})
}

func TestNewInternalError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewInternalError()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, e.Status)
		assert.Equal(t, CodeInternalError, e.ErrorCode)
		assert.Equal(t, "An internal error occurred.", e.Message)
	})
}

func TestNewNotFoundError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewNotFoundError()

		// Assert
		assert.Equal(t, http.StatusNotFound, e.Status)
		assert.Equal(t, CodeNotFound, e.ErrorCode)
		assert.Equal(t, "Not found.", e.Message)
	})
}

func TestNewMethodNotAllowedError(t *testing.T) {
	t.Run("returns correct status and codes", func(t *testing.T) {
		// Act
		e := NewMethodNotAllowedError()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, e.Status)
		assert.Equal(t, CodeMethodNotAllowed, e.ErrorCode)
		assert.Equal(t, "Method not allowed.", e.Message)
	})
}
