package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"

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

// NewAuditMiddleware is the unit of work of every admin write (ADR-0018
// addendum). The generated server hands a handler a context, not a
// transaction, so the transaction is opened here, around the handler:
//
//  1. an operation the contract marks as an admin write (permissions.IsAdminWrite:
//     the admin session scheme plus a write Permission, or, for an operation
//     that declares only an administrators-only or any-staff marker, an unsafe
//     HTTP method) gets one transaction from txStore; every other operation
//     passes through untouched;
//  2. the transaction rides in the context (store.WithAmbientTx) so the
//     handler's plain `s.db` queries and the modules it calls join it, and an
//     audit Recorder rides with it so the ent hook notes each row change;
//  3. if the handler succeeds the Audit Record is written through the same
//     transaction and both commit together; if the handler fails, or the record
//     cannot be written, both roll back.
//
// A call that did not succeed still leaves its record, outcome failed, written
// in a short transaction of its own after the rollback and holding no changes:
// the changes did not happen, and the error text or the row values could name
// what the caller was never allowed to read. A handler that answers with a
// validation error (a response, not an error) is recorded the same way, because
// nothing was changed. Refusals are recorded by the security handler, which
// rejects them first (Security.HandleAdminSession).
//
// An operation that declares `x-unit-of-work: standalone` calls a service
// outside the database (blob storage), and a connection held in a transaction
// across that call is exactly what we avoid, so it runs without one: its own
// writes commit as they happen, the recorder still notes them, and the record is
// written afterwards in its own transaction. The cost is a window between the
// change and its record in which a crash loses the record (ADR-0018 addendum).
//
// It is selected by the contract, not by a list, so a resource that declares
// itself in TypeSpec is audited with no code of its own, and the MCP server,
// which dispatches tool calls through the same api.Handler middleware, gets the
// same behaviour with a different audit.Source in the context.
func NewAuditMiddleware(txStore store.Transactor, i18n *localization.Translator, errorHandler *APIErrorHandler) api.Middleware {
	logger := errorHandler.logger
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		if !permissions.IsAdminWrite(req.OperationName) {
			return next(req)
		}
		base := req.Context

		var resp middleware.Response
		var err error
		if permissions.Standalone(req.OperationName) {
			resp, err = runStandalone(txStore, i18n, req, next)
		} else {
			resp, err = runInTransaction(txStore, i18n, req, next)
		}

		if errors.Is(err, errNotAChange) {
			err = nil
		} else if err == nil {
			return resp, nil
		}
		// Not a success: record it outside the transaction that was rolled back.
		// The caller may be gone (a cancelled request), which must not stop the
		// record of what was attempted.
		ctx := context.WithoutCancel(base)
		if recordErr := RecordOutcome(ctx, txStore, i18n, req.OperationName, audit.OutcomeFailed); recordErr != nil {
			logger.ErrorContext(ctx, "recording a failed admin call", "operation", req.OperationName, "err", recordErr)
		}
		return resp, err
	}
}

// errNotAChange is how the transaction closure says "roll back, the handler
// answered with a validation error": the response goes to the caller as it is,
// and the call is recorded as failed.
var errNotAChange = errors.New("handler refused the input")

func runInTransaction(txStore store.Transactor, i18n *localization.Translator, req middleware.Request, next middleware.Next) (middleware.Response, error) {
	var resp middleware.Response
	err := txStore.WithinTx(req.Context, func(tx *sql.Tx, client *ent.Client) error {
		ctx, recorder := audit.WithRecorder(store.WithAmbientTx(req.Context, tx, client))
		req.SetContext(ctx)

		var handlerErr error
		resp, handlerErr = next(req)
		if handlerErr != nil {
			return handlerErr
		}
		if refusedInput(resp) {
			return errNotAChange
		}
		return audit.Write(ctx, client, auditEntry(ctx, req.OperationName, i18n, audit.OutcomeSucceeded), recorder)
	})
	return resp, err
}

func runStandalone(txStore store.Transactor, i18n *localization.Translator, req middleware.Request, next middleware.Next) (middleware.Response, error) {
	ctx, recorder := audit.WithRecorder(req.Context)
	req.SetContext(ctx)

	resp, err := next(req)
	if err != nil {
		return resp, err
	}
	if refusedInput(resp) {
		return resp, errNotAChange
	}
	entry := auditEntry(ctx, req.OperationName, i18n, audit.OutcomeSucceeded)
	err = txStore.WithinTx(ctx, func(_ *sql.Tx, client *ent.Client) error {
		return audit.Write(ctx, client, entry, recorder)
	})
	return resp, err
}

// refusedInput reports a handler answer that is a validation failure: ogen
// returns it as a response with a nil error, and nothing was changed.
func refusedInput(resp middleware.Response) bool {
	_, ok := resp.Type.(*api.ValidationError)
	return ok
}

// RecordOutcome writes an Audit Record that holds no changes, for an attempt
// that did not complete: a refusal or a failure. It runs in a transaction of its
// own, because the one the attempt used is gone or never existed, and takes no
// field values by design, so the record cannot show a caller what they were not
// allowed to read.
func RecordOutcome(ctx context.Context, txStore store.Transactor, i18n *localization.Translator, operation string, outcome audit.Outcome) error {
	entry := auditEntry(ctx, operation, i18n, outcome)
	return txStore.WithinTx(ctx, func(_ *sql.Tx, client *ent.Client) error {
		return audit.Write(ctx, client, entry, nil)
	})
}

// auditEntry gathers who and from where. Everything comes from what the
// transport already put in the context (WithClientIP carries the address, user
// agent and request id the HTTP request had), so the MCP path, which has no
// http.Request, leaves address and user agent empty.
func auditEntry(ctx context.Context, operation string, i18n *localization.Translator, outcome audit.Outcome) audit.Entry {
	info, _ := ctx.Value(requestInfoKey{}).(requestInfo)
	entry := audit.Entry{
		Source:    audit.SourceFrom(ctx),
		Operation: operation,
		Locale:    i18n.Locale(ctx),
		Outcome:   outcome,
		RequestID: requestID(ctx, info),
		UserAgent: info.userAgent,
	}
	if user, ok := AuthenticatedUser(ctx); ok {
		id := user.ID
		entry.Actor = &id
	}
	entry.Address, _ = ClientIP(ctx)
	return entry
}

// requestID prefers what the caller or proxy sent, then the OpenTelemetry trace
// id of the request (the id a log line or Sentry event carries), then a fresh
// random one, so a record always has something to look a request up by.
func requestID(ctx context.Context, info requestInfo) string {
	if info.requestID != "" {
		return info.requestID
	}
	if sc := trace.SpanContextFromContext(ctx); sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return rand.Text()
}
