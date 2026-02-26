package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"policy-inference-decider/internal/apierror"
	"policy-inference-decider/internal/policy"
)

func TestErrorFromPolicy(t *testing.T) {
	t.Run("error ErrNoStartNode then returns 400 and policy_no_start_node", func(t *testing.T) {
		// Arrange
		inputErr := policy.ErrNoStartNode

		// Act
		got := errorFromPolicy(inputErr)

		// Assert
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, apierror.CodePolicyNoStartNode, got.ErrorCode)
		assert.Equal(t, "Policy graph has no start node.", got.Message)
	})
	t.Run("error ErrInvalidCondition then returns 400 and invalid_condition", func(t *testing.T) {
		// Arrange
		inputErr := policy.ErrInvalidCondition

		// Act
		got := errorFromPolicy(inputErr)

		// Assert
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, apierror.CodeInvalidCondition, got.ErrorCode)
		assert.Equal(t, "Invalid condition in policy.", got.Message)
	})
	t.Run("error returns 500 and internal_error", func(t *testing.T) {
		// Arrange
		inputErr := errors.New("some execution error")

		// Act
		got := errorFromPolicy(inputErr)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, got.Status)
		assert.Equal(t, apierror.CodeInternalError, got.ErrorCode)
		assert.Equal(t, "An internal error occurred.", got.Message)
	})
}

func TestErrorFromParseDOT(t *testing.T) {
	t.Run("when ErrNoStartNode then returns 400 and policy_no_start_node", func(t *testing.T) {
		// Arrange
		inputErr := policy.ErrNoStartNode

		// Act
		got := errorFromParseDOT(inputErr)

		// Assert
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, apierror.CodePolicyNoStartNode, got.ErrorCode)
		assert.Equal(t, "Policy graph has no start node.", got.Message)
	})
	t.Run("when other error then returns 400 and invalid_policy_dot", func(t *testing.T) {
		// Arrange
		inputErr := errors.New("syntax error: unexpected SEMI")

		// Act
		got := errorFromParseDOT(inputErr)

		// Assert
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, apierror.CodeInvalidPolicyDOT, got.ErrorCode)
		assert.Equal(t, "Invalid policy DOT format.", got.Message)
	})
}
