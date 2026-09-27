package assistant_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/assistant"
)

// The adapter hands over content deltas in order and reads the token usage
// from the final chunk, which the API only sends when the request asks for it.
func TestOpenAIStreamDeliversDeltasAndUsage(t *testing.T) {
	var request map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		w.Header().Set("Content-Type", "text/event-stream")
		for _, chunk := range []string{
			`{"id":"c","object":"chat.completion.chunk","created":1,"model":"m","choices":[{"index":0,"delta":{"role":"assistant","content":"Hel"}}]}`,
			`{"id":"c","object":"chat.completion.chunk","created":1,"model":"m","choices":[{"index":0,"delta":{"content":"lo"},"finish_reason":"stop"}]}`,
			`{"id":"c","object":"chat.completion.chunk","created":1,"model":"m","choices":[],"usage":{"prompt_tokens":42,"completion_tokens":2,"total_tokens":44}}`,
		} {
			_, _ = fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		_, _ = fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	t.Cleanup(server.Close)
	llm := assistant.NewOpenAI("token", "m", option.WithBaseURL(server.URL))

	var deltas []string
	usage, err := llm.Stream(t.Context(), []assistant.Turn{
		{Role: assistant.RoleSystem, Content: "Be brief."},
		{Role: assistant.RoleUser, Content: "Hi"},
	}, func(delta string) error {
		deltas = append(deltas, delta)
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"Hel", "lo"}, deltas)
	assert.Equal(t, assistant.Usage{InputTokens: 42, OutputTokens: 2}, usage)
	assert.Equal(t, map[string]any{"include_usage": true}, request["stream_options"])
	assert.Equal(t, "low", request["reasoning_effort"], "the first token must not wait on a long think")
}
