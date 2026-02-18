package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const exampleDOT = `digraph { start [result=""]; ok [result="approved=true"]; no [result="approved=false"]; start -> ok [cond="age>=18"]; start -> no [cond="age<18"]; }`

const dotNoStart = `digraph { foo [result=""]; bar [result="x=1"]; foo -> bar [cond="true"]; }`

const dotWithCycle = `digraph { start [result=""]; a [result="done=true"]; start -> a [cond="x==1"]; a -> a [cond="x==1"]; }`

type inferScenario struct {
	request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
	err      error
	ctx      context.Context
}

func TestInfer(t *testing.T) {
	validRequestAge20 := events.APIGatewayProxyRequest{
		Body: mustMarshal(map[string]any{
			"policy_dot": exampleDOT,
			"input":      map[string]any{"age": 20},
		}),
	}
	validRequestAge15 := events.APIGatewayProxyRequest{
		Body: mustMarshal(map[string]any{
			"policy_dot": exampleDOT,
			"input":      map[string]any{"age": 15},
		}),
	}
	invalidBodyRequest := events.APIGatewayProxyRequest{Body: "invalid"}
	requestDotNoStart := events.APIGatewayProxyRequest{
		Body: mustMarshal(map[string]any{
			"policy_dot": dotNoStart,
			"input":      map[string]any{"x": 1},
		}),
	}
	requestWithCycle := events.APIGatewayProxyRequest{
		Body: mustMarshal(map[string]any{
			"policy_dot": dotWithCycle,
			"input":      map[string]any{"x": 1},
		}),
	}

	testCases := map[string]func(t *testing.T){
		"success - approved true when age >= 18": func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(validRequestAge20)
			s.whenInferIsExecuted()
			s.thenStatusOKAndApprovedTrue(t)
		},
		"success - approved false when age < 18": func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(validRequestAge15)
			s.whenInferIsExecuted()
			s.thenStatusOKAndApprovedFalse(t)
		},
		"bad request - invalid JSON body returns APIError format": func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(invalidBodyRequest)
			s.whenInferIsExecuted()
			s.thenBadRequestWithAPIError(t)
		},
		"bad request - DOT without start node returns policy_no_start_node": func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestDotNoStart)
			s.whenInferIsExecuted()
			s.thenBadRequestWithErrorCode(t, CodePolicyNoStartNode)
		},
		"success - graph with cycle terminates and returns output": func(t *testing.T) {
			s := startInferScenario()
			s.givenARequest(requestWithCycle)
			s.whenInferIsExecuted()
			s.thenStatusOKWithDoneTrue(t)
		},
	}

	t.Parallel()
	for name, tc := range testCases {
		t.Run(name, tc)
	}
}

func startInferScenario() *inferScenario {
	return &inferScenario{}
}

func (s *inferScenario) givenARequest(req events.APIGatewayProxyRequest) {
	s.request = req
	s.ctx = context.Background()
}

func (s *inferScenario) whenInferIsExecuted() {
	s.response, s.err = Infer(s.ctx, s.request)
}

func (s *inferScenario) thenStatusOKAndApprovedTrue(t *testing.T) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusOK, s.response.StatusCode)
	var out struct {
		Output map[string]any `json:"output"`
	}
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &out))
	assert.Equal(t, float64(20), out.Output["age"])
	assert.Equal(t, true, out.Output["approved"])
}

func (s *inferScenario) thenStatusOKAndApprovedFalse(t *testing.T) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusOK, s.response.StatusCode)
	var out struct {
		Output map[string]any `json:"output"`
	}
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &out))
	assert.Equal(t, false, out.Output["approved"])
}

func (s *inferScenario) thenBadRequestWithAPIError(t *testing.T) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusBadRequest, s.response.StatusCode)
	var apiErr APIError
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &apiErr))
	assert.NotEmpty(t, apiErr.Error)
	assert.NotEmpty(t, apiErr.Message)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status)
	assert.Equal(t, CodeInvalidRequestBody, apiErr.Error)
}

func (s *inferScenario) thenBadRequestWithErrorCode(t *testing.T, wantCode string) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusBadRequest, s.response.StatusCode)
	var apiErr APIError
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &apiErr))
	assert.Equal(t, wantCode, apiErr.Error)
}

func (s *inferScenario) thenStatusOKWithDoneTrue(t *testing.T) {
	require.NoError(t, s.err)
	assert.Equal(t, http.StatusOK, s.response.StatusCode)
	var out struct {
		Output map[string]any `json:"output"`
	}
	require.NoError(t, json.Unmarshal([]byte(s.response.Body), &out))
	assert.Equal(t, true, out.Output["done"])
}

func mustMarshal(v map[string]any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
