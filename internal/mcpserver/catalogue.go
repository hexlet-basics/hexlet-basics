package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/internal/api"
	"hexletbasics/internal/oauthserver"
)

// listCoursesTool is named after the operation it mirrors (ADR-0018 names admin
// tools by operationId; the catalogue follows the same rule).
const listCoursesTool = "listCourses"

// addCatalogueTools registers the read-only catalogue tools. Each needs the
// catalog:read Scope.
func (s *Server) addCatalogueTools(server *mcp.Server) error {
	schema, err := responseSchema("listCourses")
	if err != nil {
		return fmt.Errorf("build %s output schema: %w", listCoursesTool, err)
	}
	// The contract returns a bare array, and MCP structured content is an object,
	// so the list sits under one key.
	output := map[string]any{
		"type":                 "object",
		"properties":           map[string]any{"courses": schema},
		"required":             []string{"courses"},
		"additionalProperties": false,
	}

	server.AddTool(&mcp.Tool{
		Name:         listCoursesTool,
		Title:        "List courses",
		Description:  "Lists the published Hexlet Basics courses in catalogue order, with their slug, name, duration and enrollment count.",
		InputSchema:  map[string]any{"type": "object", "additionalProperties": false},
		OutputSchema: output,
		Annotations:  &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if _, denied := requireScope(req, oauthserver.ScopeCatalogRead); denied != nil {
			return denied, nil
		}
		courses, err := s.deps.Catalogue.ListCourses(ctx)
		if err != nil {
			return nil, fmt.Errorf("list courses: %w", err)
		}
		if courses == nil {
			// An empty catalogue is an empty list: the schema has no null here.
			courses = []api.CourseCatalogItem{}
		}
		slugs := make([]string, len(courses))
		for i, c := range courses {
			slugs[i] = c.Slug
		}
		return toolResult(summarize("courses", slugs), map[string]any{"courses": courses})
	})
	return nil
}
