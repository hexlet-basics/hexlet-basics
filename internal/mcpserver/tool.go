package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/internal/domain"
)

// toolSpec is what a tool declares; the rest of a call (scope, rate limit,
// timeout, size limit) is the same for every tool and lives in addTool.
type toolSpec struct {
	// name is the operationId of the operation the tool mirrors.
	name        string
	title       string
	description string
	// input is the JSON Schema of the arguments, an object.
	input map[string]any
	// output is the contract's schema of the result; addTool adds the optional
	// `truncation` property to it.
	output map[string]any
	scope  string
	class  Class
}

// call is what a tool's body is given: the subject, and its arguments.
type call struct {
	userID   int
	clientID string
	args     json.RawMessage
}

// decode reads the arguments into v. Unknown arguments are refused rather than
// ignored: an agent that misspells `perPage` would otherwise get a default page
// and never learn why.
func (c call) decode(v any) error {
	raw := bytes.TrimSpace(c.args)
	if len(raw) == 0 {
		raw = []byte("{}")
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return userError{errs.Wrapf(err, "invalid arguments")}
	}

	return nil
}

// result is what a tool's body returns.
type result struct {
	// text is the short summary a client that shows only text displays.
	text string
	// data is the structured content, an object the output schema describes.
	data any
	// list is set when data is a page of a list, so a cut can say which page
	// arguments continue it.
	list *listInfo
}

// userError is a failure the agent can fix or must be told about, answered as a
// tool error with this message; any other error is ours and answers as a failed
// request.
type userError struct{ error }

// errNotFound is what a tool answers for a record that does not exist or that a
// visitor may not read: the two are the same answer, as on the website.
func errNotFound(what string) error { return userError{notFoundError{what}} }

// notFoundError says what was not found.
type notFoundError struct{ what string }

func (e notFoundError) Error() string { return e.what + " not found" }

// addTool registers a tool behind the checkpoint every tool shares, in this
// order: the Scope and the rate limit of its Class (admit), a context that ends
// after the call timeout, the body, and the response size limit.
func (s *Server) addTool(server *mcp.Server, spec toolSpec, run func(ctx context.Context, c call) (result, error)) {
	tool := &mcp.Tool{
		Name:         spec.name,
		Title:        spec.title,
		Description:  spec.description,
		InputSchema:  spec.input,
		OutputSchema: withTruncation(outputShape(spec.output)),
		Annotations:  annotations(spec.class),
	}
	s.note(tool, toolNeeds{scope: spec.scope})
	server.AddTool(tool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, clientID, denied := s.admit(req, spec.scope, spec.class)
		if denied != nil {
			return denied, nil
		}

		// The timeout cancels the handler's context, which is what stops its
		// database queries; a call that overruns is told so in words.
		ctx, cancel := context.WithTimeout(ctx, s.limits.ToolTimeout)
		defer cancel()

		res, err := run(ctx, call{userID: userID, clientID: clientID, args: req.Params.Arguments})
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return toolError(fmt.Sprintf("%s took longer than %s and was cancelled; try a narrower request", spec.name, s.limits.ToolTimeout)), nil
		}

		var user userError
		switch {
		case errors.As(err, &user):
			return toolError(user.Error()), nil
		case errors.Is(domain.FromStore(err), domain.ErrNotFound):
			return toolError("not found"), nil
		case err != nil:
			return nil, errs.With("tool", spec.name).Wrapf(err, "run tool")
		}

		return s.finish(spec.name, res)
	})
}

// annotations marks what a client may assume about the Class.
func annotations(class Class) *mcp.ToolAnnotations {
	return hints(class != ClassWrite, false, class != ClassWrite)
}

// hints is a tool's annotations with all four hints set. The protocol's defaults
// (destructive and open-world both true) are the unsafe ones, so a tool that
// leaves one unset is described as the worst it could be, and a client may ask
// its user to confirm a read. Nothing here reaches outside the site's own
// database, so open-world is false everywhere.
func hints(readOnly, destructive, idempotent bool) *mcp.ToolAnnotations {
	closedWorld := false

	return &mcp.ToolAnnotations{ReadOnlyHint: readOnly, DestructiveHint: &destructive, IdempotentHint: idempotent, OpenWorldHint: &closedWorld}
}

// finish turns a body's result into the tool result, within the size limit.
func (s *Server) finish(name string, res result) (*mcp.CallToolResult, error) {
	raw, err := json.Marshal(res.data)
	if err != nil {
		return nil, errs.With("tool", name).Wrapf(err, "encode result")
	}
	// Round-trip through JSON so the content is exactly what the contract's
	// encoders produce (omitted optional fields, null for nullable ones), not
	// the Go shape of the generated types.
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, errs.With("tool", name).Wrapf(err, "decode result")
	}

	text := res.text
	if !fit(data, len(text), s.limits.MaxResponseBytes, res.list) {
		return toolError(fmt.Sprintf("the result of %s is larger than the %d byte response limit and cannot be shortened", name, s.limits.MaxResponseBytes)), nil
	}

	if _, cutDown := data[truncationKey]; cutDown {
		text += " (cut to fit the response size limit: see the truncation property)"
	}

	return &mcp.CallToolResult{
		Content:           []mcp.Content{&mcp.TextContent{Text: text}},
		StructuredContent: data,
	}, nil
}

// object is the input schema of a tool: a closed object.
func object(properties map[string]any, required ...string) map[string]any {
	schema := map[string]any{"type": "object", "properties": properties, "additionalProperties": false}
	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// slugArg is a string argument naming a record by its slug.
func slugArg(description string) map[string]any {
	return map[string]any{"type": "string", "minLength": 1, "description": description}
}
