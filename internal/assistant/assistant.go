// Package assistant hosts the LLM access for the AI features: the lesson-review
// summarizer and the in-lesson chat (see the assistant port design — one
// shared client for both, synchronous calls, no persistence of provider
// internals).
package assistant

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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
// wiring (nil-skipped workers when absent), not validated here. Extra SDK
// options point a test at a local endpoint.
func NewOpenAI(token, model string, opts ...option.RequestOption) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(append([]option.RequestOption{option.WithAPIKey(token)}, opts...)...),
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

// Stream runs a chat completion over turns and hands each content delta to
// onDelta as it arrives. The SDK's accumulator collects the final usage, which
// the API only sends when asked for it (include_usage) — without it a streamed
// answer would carry no token counts. An error from onDelta (the learner went
// away) aborts the stream.
func (c *OpenAI) Stream(ctx context.Context, turns []Turn, onDelta func(string) error) (Usage, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(turns))
	for i, turn := range turns {
		switch turn.Role {
		case RoleSystem:
			messages[i] = openai.SystemMessage(turn.Content)
		case RoleAssistant:
			messages[i] = openai.AssistantMessage(turn.Content)
		default:
			messages[i] = openai.UserMessage(turn.Content)
		}
	}

	stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:         c.model,
		Messages:      messages,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{IncludeUsage: openai.Bool(true)},
	})
	// Close only releases the response body; the stream's own error is what
	// reports a failed answer.
	defer func() { _ = stream.Close() }()

	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		if len(chunk.Choices) == 0 || chunk.Choices[0].Delta.Content == "" {
			continue
		}
		if err := onDelta(chunk.Choices[0].Delta.Content); err != nil {
			return Usage{}, err
		}
	}
	if err := stream.Err(); err != nil {
		return Usage{}, oops.Wrapf(err, "openai chat completion stream")
	}
	return Usage{
		InputTokens:  int(acc.Usage.PromptTokens),
		OutputTokens: int(acc.Usage.CompletionTokens),
	}, nil
}
