package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/ogen-go/ogen/middleware"

	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/events"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/permissions"
)

// publishTimeout bounds how long a read may wait to hand its fact to the
// outbox. The fact is one small insert; if the database cannot do that in this
// time the record is given up (and logged) rather than holding the response.
const publishTimeout = 2 * time.Second

// NewPersonalDataReadMiddleware records staff reads of personal data (ADR-0020
// addendum), the sibling of NewAuditMiddleware for the operations that change
// nothing. It is selected by the contract: an operation that declares
// `x-personal-data` (permissions.PersonalDataEntity) is recorded, every other
// operation, the catalogue and content included, passes through untouched, and
// a handler has no code of its own for it.
//
// Recording is asynchronous. After the handler has answered, this publishes a
// domain event (events.StaffReadPersonalData) carrying who, from where, which
// operation and the ids or the count of what the response held; a worker
// consumer writes the Audit Record. So the read never waits for the partitioned
// audit insert, and nothing in it can fail the read: a response is passed
// through unchanged whether or not the fact could be published, and a failure is
// logged. The publish itself is bounded by publishTimeout, on a context that a
// client hanging up does not cancel, because the read already happened.
//
// Only an answer that holds data is a read: a handler error, or a response that
// is an error body rather than a row or a list, records nothing. A refused call
// never gets here; the security handler records it (Security.HandleAdminSession).
// The MCP server dispatches through the same middleware with audit.WithSource
// set, so the same record carries the channel and Connected App.
func NewPersonalDataReadMiddleware(publisher events.StandalonePublisher, i18n *localization.Translator, errorHandler *APIErrorHandler) api.Middleware {
	logger := errorHandler.logger
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		entity, personal := permissions.PersonalDataEntity(req.OperationName)
		if !personal {
			return next(req)
		}
		resp, err := next(req)
		if err != nil {
			return resp, err
		}
		read, ok := touched(resp.Type)
		if !ok {
			return resp, nil
		}

		ctx, cancel := context.WithTimeout(context.WithoutCancel(req.Context), publishTimeout)
		defer cancel()
		entry := auditEntry(ctx, req.OperationName, i18n, audit.OutcomeSucceeded)
		fact := events.StaffReadPersonalData{
			ActorUserID:    entry.Actor,
			Channel:        string(entry.Source.Channel),
			ConnectedAppID: entry.Source.ConnectedAppID,
			TokenID:        entry.Source.TokenID,
			Operation:      entry.Operation,
			EntityType:     entity,
			EntityIDs:      read.IDs,
			Count:          read.Count,
			RequestID:      entry.RequestID,
			Address:        entry.Address,
			UserAgent:      entry.UserAgent,
			Locale:         entry.Locale,
			OccurredAt:     time.Now(),
		}
		if publishErr := publisher.PublishStandalone(ctx, fact); publishErr != nil {
			logger.ErrorContext(ctx, "publishing a read of personal data", "operation", req.OperationName, "err", publishErr)
		}
		return resp, nil
	}
}

// touched says what a response held, from its JSON form so no handler or
// response type needs to cooperate: a page (`items`) or a bare list is a count
// and the ids of its rows, an object with an `id` is one row. Anything else (an
// error body, an empty answer) is not a read of rows and reports false. The
// response is encoded a second time here; a page is small and bounded.
func touched(response any) (audit.Reads, bool) {
	raw, err := json.Marshal(response)
	if err != nil {
		return audit.Reads{}, false
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return audit.Reads{}, false
	}
	switch v := value.(type) {
	case []any:
		return readsOf(v), true
	case map[string]any:
		if items, ok := v["items"].([]any); ok {
			return readsOf(items), true
		}
		if id, ok := idOf(v); ok {
			return audit.Reads{Count: 1, IDs: []string{id}}, true
		}
	}
	return audit.Reads{}, false
}

func readsOf(rows []any) audit.Reads {
	reads := audit.Reads{Count: len(rows)}
	for _, row := range rows {
		if len(reads.IDs) == audit.MaxReadIDs {
			break
		}
		if object, ok := row.(map[string]any); ok {
			if id, ok := idOf(object); ok {
				reads.IDs = append(reads.IDs, id)
			}
		}
	}
	return reads
}

func idOf(object map[string]any) (string, bool) {
	switch id := object["id"].(type) {
	case json.Number:
		return id.String(), true
	case string:
		return id, id != ""
	}
	return "", false
}
