package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"policy-inference-decider/internal/policy"
)

type errorFromPolicyScenario struct {
	inputErr     error
	gotCode      int
	gotErrorCode string
	gotMessage   string
	gotCause     []any
}

func TestErrorFromPolicy(t *testing.T) {
	errNoStartNode := policy.ErrNoStartNode
	otherError := errors.New("some execution error")

	testCases := map[string]func(t *testing.T){
		"error - ErrNoStartNode then returns 400 and policy_no_start_node": func(t *testing.T) {
			s := startErrorFromPolicyScenario()
			s.givenAnError(errNoStartNode)
			s.whenErrorFromPolicyIsCalled()
			s.thenReturns(t, http.StatusBadRequest, CodePolicyNoStartNode, errNoStartNode.Error(), nil)
		},
		"error - returns 500 and internal_error": func(t *testing.T) {
			s := startErrorFromPolicyScenario()
			s.givenAnError(otherError)
			s.whenErrorFromPolicyIsCalled()
			s.thenReturns(t, http.StatusInternalServerError, CodeInternalError, "An internal error occurred.", nil)
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
	s.gotCode, s.gotErrorCode, s.gotMessage, s.gotCause = errorFromPolicy(s.inputErr)
}

func (s *errorFromPolicyScenario) thenReturns(t *testing.T, wantCode int, wantErrorCode, wantMessage string, wantCause []any) {
	assert.Equal(t, wantCode, s.gotCode)
	assert.Equal(t, wantErrorCode, s.gotErrorCode)
	assert.Equal(t, wantMessage, s.gotMessage)
	if wantCause != nil || s.gotCause != nil {
		assert.Len(t, s.gotCause, len(wantCause))
	}
}

type errorFromParseDOTScenario struct {
	inputErr     error
	gotCode      int
	gotErrorCode string
	gotMessage   string
	gotCause     []any
}

func TestErrorFromParseDOT(t *testing.T) {
	errNoStartNode := policy.ErrNoStartNode
	invalidDOTError := errors.New("syntax error: unexpected SEMI")

	testCases := map[string]func(t *testing.T){
		"test when ErrNoStartNode then returns 400 and policy_no_start_node": func(t *testing.T) {
			s := startErrorFromParseDOTScenario()
			s.givenAnError(errNoStartNode)
			s.whenErrorFromParseDOTIsCalled()
			s.thenReturns(t, http.StatusBadRequest, CodePolicyNoStartNode, errNoStartNode.Error(), nil)
		},
		"test when other error then returns 400 and invalid_policy_dot": func(t *testing.T) {
			s := startErrorFromParseDOTScenario()
			s.givenAnError(invalidDOTError)
			s.whenErrorFromParseDOTIsCalled()
			s.thenReturns(t, http.StatusBadRequest, CodeInvalidPolicyDOT, invalidDOTError.Error(), nil)
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
	s.gotCode, s.gotErrorCode, s.gotMessage, s.gotCause = errorFromParseDOT(s.inputErr)
}

func (s *errorFromParseDOTScenario) thenReturns(t *testing.T, wantCode int, wantErrorCode, wantMessage string, wantCause []any) {
	assert.Equal(t, wantCode, s.gotCode)
	assert.Equal(t, wantErrorCode, s.gotErrorCode)
	assert.Equal(t, wantMessage, s.gotMessage)
	if wantCause != nil || s.gotCause != nil {
		assert.Len(t, s.gotCause, len(wantCause))
	}
}
