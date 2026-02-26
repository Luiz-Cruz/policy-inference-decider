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

type inferResponseBody struct {
	Output map[string]any `json:"output"`
}

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

func bodyFromInferRequest(r policy.InferRequest) string {
	b, _ := json.Marshal(r)
	return string(b)
}

func TestInfer(t *testing.T) {
	t.Run("success - approved true when age >= 18", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 20}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var out inferResponseBody
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &out))
		assert.Equal(t, inferResponseBody{Output: map[string]any{"age": float64(20), "approved": true}}, out)
	})
	t.Run("success - approved false when age < 18", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 15}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var out inferResponseBody
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &out))
		assert.Equal(t, inferResponseBody{Output: map[string]any{"age": float64(15), "approved": false}}, out)
	})
	t.Run("bad request - invalid JSON body returns APIError format", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		req := makeURLRequest("invalid", http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.NotEmpty(t, apiErr.Message)
		assert.Equal(t, http.StatusBadRequest, apiErr.Status)
		assert.Equal(t, apierror.CodeInvalidRequestBody, apiErr.Error)
	})
	t.Run("bad request - DOT without start node returns policy_no_start_node", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotNoStart, Input: map[string]any{"x": 1}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodePolicyNoStartNode, apiErr.Error)
	})
	t.Run("success - graph with cycle terminates and returns output", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotWithCycle, Input: map[string]any{"x": 1}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var out inferResponseBody
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &out))
		assert.Equal(t, inferResponseBody{Output: map[string]any{"x": float64(1), "done": true}}, out)
	})
	t.Run("bad request - invalid condition in edge returns invalid_condition", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: dotWithInvalidCond, Input: map[string]any{"x": 1}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodeInvalidCondition, apiErr.Error)
	})
	t.Run("bad request - invalid DOT format returns invalid_policy_dot", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: dothWithInvalidFormat, Input: map[string]any{"age": 25}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodeInvalidPolicyDOT, apiErr.Error)
	})
	t.Run("challenge example - Policy graph with age 25 score 720 returns approved and segment prime", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: policyChallengeDOT, Input: map[string]any{"age": 25, "score": 720}})
		req := makeURLRequest(body, http.MethodPost, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var out inferResponseBody
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &out))
		assert.Equal(t, inferResponseBody{Output: map[string]any{"age": float64(25), "score": float64(720), "approved": true, "segment": "prime"}}, out)
	})
	t.Run("not found when path is not /infer", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		req := makeURLRequest("", http.MethodPost, "/other")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodeNotFound, apiErr.Error)
	})
	t.Run("method not allowed when not POST", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 20}})
		req := makeURLRequest(body, http.MethodGet, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodeMethodNotAllowed, apiErr.Error)
	})
	t.Run("GET /ping returns pong", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		req := makeURLRequest("", http.MethodGet, "/ping")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "pong", resp.Body)
	})
	t.Run("unsupported method returns 405", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		body := bodyFromInferRequest(policy.InferRequest{PolicyDOT: exampleDOT, Input: map[string]any{"age": 20}})
		req := makeURLRequest(body, http.MethodPut, "/infer")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodeMethodNotAllowed, apiErr.Error)
	})
	t.Run("GET other path returns 404", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		req := makeURLRequest("", http.MethodGet, "/other")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var apiErr APIError
		require.NoError(t, json.Unmarshal([]byte(resp.Body), &apiErr))
		assert.Equal(t, apierror.CodeNotFound, apiErr.Error)
	})
	t.Run("pathFromRequest uses RawPath when HTTP.Path empty", func(t *testing.T) {
		// Arrange
		h := NewInferHandler(policy.NewDotParser(), policy.NewGraphExecutor())
		req := makeURLRequestWithRawPath("", http.MethodGet, "/ping")

		// Act
		resp, err := h.Infer(context.Background(), req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "pong", resp.Body)
	})
}
