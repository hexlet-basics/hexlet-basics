package handlers

import (
	"context"
	"errors"
	"time"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/localization"
)

// AdminListAuditRecords lists the audit trail, newest first. Who may call it is
// the contract's business: the operation declares the capability Permission
// `audit.read`, which the security handler requires as an explicit grant even of
// an administrator, and `x-personal-data`, so the read middleware records that
// the trail was read (the records carry addresses and old values).
//
// Every query is bounded in time because the table is partitioned by month on
// occurred_at: audit.Filter resolves the window (30 days back from now by
// default, at most 92 days wide) and a request outside those rules is a 422
// naming the offending field, rather than a scan of every partition.
func (s *Server) AdminListAuditRecords(ctx context.Context, params api.AdminListAuditRecordsParams) (api.AdminListAuditRecordsRes, error) {
	filter := auditFilterOf(params)
	// One instant for the check, the count and the page, so the window cannot
	// shift between them.
	now := time.Now()
	if _, err := filter.Query(s.db, now); err != nil {
		return s.auditWindowProblem(ctx, err)
	}

	page, err := listPageWith(ctx, params.Page, params.PerPage,
		func() *ent.AuditRecordQuery {
			// Valid: the same filter and instant were accepted above.
			q, _ := filter.Query(s.db, now)
			return q
		},
		func(_ context.Context, rows []*ent.AuditRecord) ([]api.AuditRecord, error) {
			items := make([]api.AuditRecord, 0, len(rows))
			for _, row := range rows {
				item, err := apiconv.AuditRecord(row)
				if err != nil {
					return nil, err
				}
				items = append(items, item)
			}
			return items, nil
		},
		func(items []api.AuditRecord, total, page, perPage int32) *api.AuditRecordPage {
			return &api.AuditRecordPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// auditFilterOf maps the query parameters to the audit package's filter.
func auditFilterOf(params api.AdminListAuditRecordsParams) audit.Filter {
	var filter audit.Filter
	if v, ok := params.ActorUserId.Get(); ok {
		actor := int(v)
		filter.ActorUserID = &actor
	}
	if v, ok := params.Channel.Get(); ok {
		channel := audit.Channel(v)
		filter.Channel = &channel
	}
	if v, ok := params.Operation.Get(); ok {
		filter.Operation = &v
	}
	if v, ok := params.EntityType.Get(); ok {
		filter.EntityType = &v
	}
	if v, ok := params.EntityId.Get(); ok {
		filter.EntityID = &v
	}
	if v, ok := params.From.Get(); ok {
		filter.From = &v
	}
	if v, ok := params.To.Get(); ok {
		filter.To = &v
	}
	return filter
}

// auditWindowProblem turns a refused window or filter into the 422 the form
// shows next to the field.
func (s *Server) auditWindowProblem(ctx context.Context, err error) (api.AdminListAuditRecordsRes, error) {
	switch {
	case errors.Is(err, audit.ErrWindowReversed):
		return validationError("to", s.i18n.Text(ctx, localization.AuditWindowReversed)), nil
	case errors.Is(err, audit.ErrWindowTooWide):
		return validationError("from", s.i18n.Text(ctx, localization.AuditWindowTooWide)), nil
	case errors.Is(err, audit.ErrEntityIDWithoutType):
		return validationError("entityId", s.i18n.Text(ctx, localization.AuditEntityIDNeedsType)), nil
	}
	return nil, err
}
