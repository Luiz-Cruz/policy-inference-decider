package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"policy-inference-decider/internal/apierror"
	"policy-inference-decider/internal/policy"
)

type errorFromPolicyScenario struct {
	inputErr  error
	gotAPIErr apierror.APIError
}

func TestErrorFromPolicy(t *testing.T) {
	errNoStartNode := policy.ErrNoStartNode
	otherError := errors.New("some execution error")

	errInvalidCondition := policy.ErrInvalidCondition

	testCases := map[string]func(t *testing.T){
		"error - ErrNoStartNode then returns 400 and policy_no_start_node": func(t *testing.T) {
			s := startErrorFromPolicyScenario()
			s.givenAnError(errNoStartNode)
			s.whenErrorFromPolicyIsCalled()
			s.thenReturns(t, http.StatusBadRequest, apierror.CodePolicyNoStartNode, "Policy graph has no start node.")
		},
		"error - ErrInvalidCondition then returns 400 and invalid_condition": func(t *testing.T) {
			s := startErrorFromPolicyScenario()
			s.givenAnError(errInvalidCondition)
			s.whenErrorFromPolicyIsCalled()
			s.thenReturns(t, http.StatusBadRequest, apierror.CodeInvalidCondition, "Invalid condition in policy.")
		},
		"error - returns 500 and internal_error": func(t *testing.T) {
			s := startErrorFromPolicyScenario()
			s.givenAnError(otherError)
			s.whenErrorFromPolicyIsCalled()
			s.thenReturns(t, http.StatusInternalServerError, apierror.CodeInternalError, "An internal error occurred.")
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc)
	}
}

func startErrorFromPolicyScenario() *errorFromPolicyScenario {
	return &errorFromPolicyScenario{}
}

func (s *errorFromPolicyScenario) givenAnError(err error) {
	s.inputErr = err
}

func (s *errorFromPolicyScenario) whenErrorFromPolicyIsCalled() {
	s.gotAPIErr = errorFromPolicy(s.inputErr)
}

func (s *errorFromPolicyScenario) thenReturns(t *testing.T, wantCode int, wantErrorCode, wantMessage string) {
	assert.Equal(t, wantCode, s.gotAPIErr.Status)
	assert.Equal(t, wantErrorCode, s.gotAPIErr.ErrorCode)
	assert.Equal(t, wantMessage, s.gotAPIErr.Message)
}

type errorFromParseDOTScenario struct {
	inputErr  error
	gotAPIErr apierror.APIError
}

func TestErrorFromParseDOT(t *testing.T) {
	errNoStartNode := policy.ErrNoStartNode
	invalidDOTError := errors.New("syntax error: unexpected SEMI")

	testCases := map[string]func(t *testing.T){
		"test when ErrNoStartNode then returns 400 and policy_no_start_node": func(t *testing.T) {
			s := startErrorFromParseDOTScenario()
			s.givenAnError(errNoStartNode)
			s.whenErrorFromParseDOTIsCalled()
			s.thenReturns(t, http.StatusBadRequest, apierror.CodePolicyNoStartNode, "Policy graph has no start node.")
		},
		"test when other error then returns 400 and invalid_policy_dot": func(t *testing.T) {
			s := startErrorFromParseDOTScenario()
			s.givenAnError(invalidDOTError)
			s.whenErrorFromParseDOTIsCalled()
			s.thenReturns(t, http.StatusBadRequest, apierror.CodeInvalidPolicyDOT, "Invalid policy DOT format.")
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc)
	}
}

func startErrorFromParseDOTScenario() *errorFromParseDOTScenario {
	return &errorFromParseDOTScenario{}
}

func (s *errorFromParseDOTScenario) givenAnError(err error) {
	s.inputErr = err
}

func (s *errorFromParseDOTScenario) whenErrorFromParseDOTIsCalled() {
	s.gotAPIErr = errorFromParseDOT(s.inputErr)
}

func (s *errorFromParseDOTScenario) thenReturns(t *testing.T, wantCode int, wantErrorCode, wantMessage string) {
	assert.Equal(t, wantCode, s.gotAPIErr.Status)
	assert.Equal(t, wantErrorCode, s.gotAPIErr.ErrorCode)
	assert.Equal(t, wantMessage, s.gotAPIErr.Message)
}
