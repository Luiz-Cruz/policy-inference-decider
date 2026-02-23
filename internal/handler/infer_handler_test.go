package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"policy-inference-decider/internal/apierror"
	"policy-inference-decider/internal/policy"
)

const exampleDOT = `digraph { start [result=""]; ok [result="approved=true"]; no [result="approved=false"]; start -> ok [cond="age>=18"]; start -> no [cond="age<18"]; }`

const dotNoStart = `digraph { foo [result=""]; bar [result="x=1"]; foo -> bar [cond="true"]; }`

const dotWithCycle = `digraph { start [result=""]; a [result="done=true"]; start -> a [cond="x==1"]; a -> a [cond="x==1"]; }`

const dotWithInvalidCond = `digraph { start [result=""]; end [result="x=1"]; start -> end [cond="invalid!!!"]; }`

const dothWithInvalidFormat = `digraph Policy { start [result=\"\"]; ok [result=\"approved=true\"]; no [result=\"approved=false\"]; start -> ok [cond=\"age>=18\"]; start -> no [cond=\"age<18\"; }`

const policyChallengeDOT = `digraph Policy { start [result=""] approved [result="approved=true,segment=prime"] rejected [result="approved=false"] review [result="approved=false,segment=manual"] start -> approved [cond="age>=18 && score>700"] start -> review [cond="age>=18 && score<=700"] start -> rejected [cond="age<18"] }`

type (
	inferScenario struct {
		handler  *Handler
		request  events.APIGatewayProxyRequest
		response events.APIGatewayProxyResponse
		err      error
		ctx      context.Context
	}
	inferResponseBody struct {
		Output map[string]any `json:"output"`
	}
)

func TestInfer(t *testing.T) {
	validRequestAge20 := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 20}}),
	}
	validRequestAge15 := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 15}}),
	}
	invalidBodyRequest := events.APIGatewayProxyRequest{Body: "invalid"}
	requestDotNoStart := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotNoStart, Input: map[string]any{"x": 1}}),
	}
	requestWithCycle := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotWithCycle, Input: map[string]any{"x": 1}}),
	}
	requestInvalidCond := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotWithInvalidCond, Input: map[string]any{"x": 1}}),
	}
	requestInvalidFormat := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: dothWithInvalidFormat, Input: map[string]any{"age": 25}}),
	}
	requestChallengePolicy := events.APIGatewayProxyRequest{
		Body: bodyFromInferRequest(policy.InferRequest{PolicyDOT: policyChallengeDOT, Input: map[string]any{"age": 25, "score": 720}}),
	}

	testCases := map[string]struct {
		run func(t *testing.T)
	}{
		"success - approved true when age >= 18": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(validRequestAge20)
			s.whenInferIsExecuted()
			s.thenStatusOKWithOutput(t, inferResponseBody{Output: map[string]any{"age": float64(20), "approved": true}})
		}},
		"success - approved false when age < 18": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(validRequestAge15)
			s.whenInferIsExecuted()
			s.thenStatusOKWithOutput(t, inferResponseBody{Output: map[string]any{"age": float64(15), "approved": false}})
		}},
		"bad request - invalid JSON body returns APIError format": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(invalidBodyRequest)
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t, http.StatusBadRequest, apierror.CodeInvalidRequestBody)
		}},
		"bad request - DOT without start node returns policy_no_start_node": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestDotNoStart)
			s.whenInferIsExecuted()
			s.thenBadRequestWithErrorCode(t, apierror.CodePolicyNoStartNode)
		}},
		"success - graph with cycle terminates and returns output": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestWithCycle)
			s.whenInferIsExecuted()
			s.thenStatusOKWithOutput(t, inferResponseBody{Output: map[string]any{"x": float64(1), "done": true}})
		}},
		"bad request - invalid condition in edge returns invalid_condition": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestInvalidCond)
			s.whenInferIsExecuted()
			s.thenBadRequestWithErrorCode(t, apierror.CodeInvalidCondition)
		}},
		"bad request - invalid DOT format returns invalid_policy_dot": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestInvalidFormat)
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t, http.StatusBadRequest, apierror.CodeInvalidPolicyDOT)
		}},
		"challenge example - Policy graph with age 25 score 720 returns approved and segment prime": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestChallengePolicy)
			s.whenInferIsExecuted()
			s.thenStatusOKWithOutput(t, inferResponseBody{Output: map[string]any{"age": float64(25), "score": float64(720), "approved": true, "segment": "prime"}})
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startInferScenario() *inferScenario {
	return &inferScenario{handler: NewInferHandler(policy.NewPolicyExecutor(&policy.GraphExecutor{}, &policy.DotParser{}))}
}

func (s *inferScenario) givenARequest(req events.APIGatewayProxyRequest) {
	s.request = req
	s.ctx = context.Background()
}

func (s *inferScenario) whenInferIsExecuted() {
	s.response, s.err = s.handler.Infer(s.ctx, s.request)
}

func (s *inferScenario) thenStatusOKWithOutput(t *testing.T, want inferResponseBody) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusOK, s.response.StatusCode)
	var out inferResponseBody
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &out))
	assert.Equal(t, want, out)
}

func (s *inferScenario) thenBadRequestWithAPIError(t *testing.T, wantStatus int, wantErrorCode string) {
	require.NoError(t, s.err)
	assert.Equal(t, wantStatus, s.response.StatusCode)
	var apiErr APIError
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &apiErr))
	assert.NotEmpty(t, apiErr.Message)
	assert.Equal(t, wantStatus, apiErr.Status)
	assert.Equal(t, wantErrorCode, apiErr.Error)
}

func (s *inferScenario) thenBadRequestWithErrorCode(t *testing.T, wantCode string) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusBadRequest, s.response.StatusCode)
	var apiErr APIError
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &apiErr))
	assert.Equal(t, wantCode, apiErr.Error)
}

func bodyFromInferRequest(r policy.InferRequest) string {
	b, _ := json.Marshal(r)
	return string(b)
}
