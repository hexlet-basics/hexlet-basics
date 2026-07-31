// Package assistant hosts the LLM access for the AI features: the lesson-review
// summarizer today, the in-lesson chat assistant next (see the assistant port
// design — one shared client for both, synchronous calls, no persistence of
// provider internals).
package assistant

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/samber/oops"
)

// OpenAI is the thin completion adapter over the official SDK. The provider
// and env contract match the legacy Rails deployment (RubyLLM + OpenAI), so
// the secret carries over unchanged; the model is configuration, not data.
type OpenAI struct {
	client openai.Client
	model  openai.ChatModel
}

// NewOpenAI builds the shared client. The token is required by the callers'
// wiring (nil-skipped workers when absent), not validated here.
func NewOpenAI(token, model string) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(option.WithAPIKey(token)),
		model:  openai.ChatModel(model),
	}
}

// Complete runs a single system+user chat completion and returns the text.
func (c *OpenAI) Complete(ctx context.Context, instructions, prompt string) (string, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(instructions),
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return "", oops.Wrapf(err, "openai chat completion")
	}
	if len(resp.Choices) == 0 {
		return "", oops.Errorf("openai chat completion returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
