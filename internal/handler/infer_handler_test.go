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
		request  events.LambdaFunctionURLRequest
		response events.LambdaFunctionURLResponse
		err      error
		ctx      context.Context
	}
	inferResponseBody struct {
		Output map[string]any `json:"output"`
	}
)

func makeURLRequest(body, method, path string) events.LambdaFunctionURLRequest {
	return events.LambdaFunctionURLRequest{
		Body: body,
		RequestContext: events.LambdaFunctionURLRequestContext{
			RequestID: "test",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: method,
				Path:   path,
			},
		},
	}
}

func makeURLRequestWithRawPath(body, method, rawPath string) events.LambdaFunctionURLRequest {
	return events.LambdaFunctionURLRequest{
		Body:    body,
		RawPath: rawPath,
		RequestContext: events.LambdaFunctionURLRequestContext{
			RequestID: "test",
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: method,
				Path:   "",
			},
		},
	}
}

func TestInfer(t *testing.T) {
	bodyAge20 := bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 20}})
	bodyAge15 := bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 15}})
	validRequestAge20 := makeURLRequest(bodyAge20, http.MethodPost, "/infer")
	validRequestAge15 := makeURLRequest(bodyAge15, http.MethodPost, "/infer")
	invalidBodyRequest := makeURLRequest("invalid", http.MethodPost, "/infer")
	requestDotNoStart := makeURLRequest(bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotNoStart, Input: map[string]any{"x": 1}}), http.MethodPost, "/infer")
	requestWithCycle := makeURLRequest(bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotWithCycle, Input: map[string]any{"x": 1}}), http.MethodPost, "/infer")
	requestInvalidCond := makeURLRequest(bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotWithInvalidCond, Input: map[string]any{"x": 1}}), http.MethodPost, "/infer")
	requestInvalidFormat := makeURLRequest(bodyFromInferRequest(policy.InferRequest{PolicyDOT: dothWithInvalidFormat, Input: map[string]any{"age": 25}}), http.MethodPost, "/infer")
	requestChallengePolicy := makeURLRequest(bodyFromInferRequest(policy.InferRequest{PolicyDOT: policyChallengeDOT, Input: map[string]any{"age": 25, "score": 720}}), http.MethodPost, "/infer")
	requestPingWithRawPath := makeURLRequestWithRawPath("", http.MethodGet, "/ping")

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
		"not found when path is not /infer": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(makeURLRequest("", http.MethodPost, "/other"))
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t, http.StatusNotFound, apierror.CodeNotFound)
		}},
		"method not allowed when not POST": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(makeURLRequest(bodyAge20, http.MethodGet, "/infer"))
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t, http.StatusMethodNotAllowed, apierror.CodeMethodNotAllowed)
		}},
		"GET /ping returns pong": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(makeURLRequest("", http.MethodGet, "/ping"))
			s.whenInferIsExecuted()
			s.thenStatusOKWithPlainBody(t, "pong")
		}},
		"unsupported method returns 405": {run: func(t *testing.T) {
			s := startInferScenario()
			req := makeURLRequest(bodyAge20, http.MethodPut, "/infer")
			s.givenARequest(req)
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t, http.StatusMethodNotAllowed, apierror.CodeMethodNotAllowed)
		}},
		"GET other path returns 404": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(makeURLRequest("", http.MethodGet, "/other"))
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t, http.StatusNotFound, apierror.CodeNotFound)
		}},
		"pathFromRequest uses RawPath when HTTP.Path empty": {run: func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestPingWithRawPath)
			s.whenInferIsExecuted()
			s.thenStatusOKWithPlainBody(t, "pong")
		}},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc.run)
	}
}

func startInferScenario() *inferScenario {
	return &inferScenario{handler: NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())}
}

func (s *inferScenario) givenARequest(req events.LambdaFunctionURLRequest) {
	if req.RequestContext.HTTP.Path == "" && req.RawPath == "" {
		req.RequestContext.HTTP.Path = "/infer"
	}
	if req.RequestContext.HTTP.Method == "" {
		req.RequestContext.HTTP.Method = http.MethodPost
	}
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

func (s *inferScenario) thenStatusOKWithPlainBody(t *testing.T, wantBody string) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusOK, s.response.StatusCode)
	assert.Equal(t, wantBody, s.response.Body)
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
