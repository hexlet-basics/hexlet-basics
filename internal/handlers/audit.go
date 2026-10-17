package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"

	"github.com/ogen-go/ogen/middleware"
	"go.opentelemetry.io/otel/trace"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/permissions"
	"hexletbasics/internal/store"
)

// requestIDHeader is the header a proxy or client may set to correlate a call.
const requestIDHeader = "X-Request-ID"

// NewAuditMiddleware is the unit of work of every admin write (ADR-0017
// addendum). The generated server hands a handler a context, not a
// transaction, so the transaction is opened here, around the handler:
//
//  1. an operation whose declared Permission is a write gets one transaction
//     from txStore; every other operation passes through untouched;
//  2. the transaction rides in the context (store.WithAmbientTx) so the
//     handler's plain `s.db` queries and the modules it calls join it, and an
//     audit Recorder rides with it so the ent hook notes each row change;
//  3. if the handler succeeds the Audit Record is written through the same
//     transaction and both commit together; if the handler fails, or the record
//     cannot be written, both roll back.
//
// It is selected by the contract (the operation's x-permission kind), not by a
// list, so a resource that declares its Permissions is audited with no code of
// its own, and the MCP server, which dispatches tool calls through the same
// api.Handler middleware, gets the same behaviour with a different
// audit.Source in the context.
//
// Refusals never get here (the security handler rejects them first); recording
// those and handler failures is ticket #845, which can use the same
// audit.Write with the other outcomes.
func NewAuditMiddleware(txStore store.Transactor, i18n *localization.Translator) api.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		requirement, declared := permissions.For(req.OperationName)
		if !declared || requirement.Kind != permissions.Write {
			return next(req)
		}

		var resp middleware.Response
		err := txStore.WithinTx(req.Context, func(tx *sql.Tx, client *ent.Client) error {
			ctx, recorder := audit.WithRecorder(store.WithAmbientTx(req.Context, tx, client))
			req.SetContext(ctx)

			var handlerErr error
			resp, handlerErr = next(req)
			if handlerErr != nil {
				return handlerErr
			}
			return audit.Write(ctx, client, auditEntry(ctx, req, i18n), recorder)
		})
		if err != nil {
			return middleware.Response{}, err
		}
		return resp, nil
	}
}

// auditEntry gathers who and from where. Everything comes from what the
// transport already put in the context or on the request; the MCP path, which
// has no http.Request, leaves address and user agent empty.
func auditEntry(ctx context.Context, req middleware.Request, i18n *localization.Translator) audit.Entry {
	entry := audit.Entry{
		Source:    audit.SourceFrom(ctx),
		Operation: req.OperationName,
		Locale:    i18n.Locale(ctx),
		Outcome:   audit.OutcomeSucceeded,
		RequestID: requestID(ctx, req),
	}
	if user, ok := AuthenticatedUser(ctx); ok {
		id := user.ID
		entry.Actor = &id
	}
	entry.Address, _ = ClientIP(ctx)
	if req.Raw != nil {
		entry.UserAgent = req.Raw.UserAgent()
	}
	return entry
}

// requestID prefers what the caller or proxy sent, then the OpenTelemetry trace
// id of the request (the id a log line or Sentry event carries), then a fresh
// random one, so a record always has something to look a request up by.
func requestID(ctx context.Context, req middleware.Request) string {
	if req.Raw != nil {
		if id := req.Raw.Header.Get(requestIDHeader); id != "" {
			return id
		}
	}
	if sc := trace.SpanContextFromContext(ctx); sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return rand.Text()
}
