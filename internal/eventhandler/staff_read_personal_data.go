package eventhandler

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"hexletbasics/ent"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/events"
)

const staffReadPersonalDataHandlerName = "staff_read_personal_data_to_audit_v1"

// StaffReadPersonalData turns the fact that a staff member read personal data
// into an Audit Record (ADR-0020). It runs in the worker, so the request that
// published the fact never waits for this insert; cmd/server only publishes
// (ADR-0010).
//
// The record is written through the existing audit write path with the time the
// read happened, not the time of this delivery, and holds identifiers and a count
// only. Delivery is at least once, so a redelivered fact can leave a second
// record of the same read; an extra line in a trail is the safe direction to
// err, and the ADR already accepts that a record may be lost on a crash.
func StaffReadPersonalData(client *ent.Client) cqrs.EventHandler {
	return cqrs.NewEventHandler(
		staffReadPersonalDataHandlerName,
		func(ctx context.Context, event *events.StaffReadPersonalData) error {
			entry := audit.Entry{
				Actor: event.ActorUserID,
				Source: audit.Source{
					Channel:        audit.Channel(event.Channel),
					ConnectedAppID: event.ConnectedAppID,
					TokenID:        event.TokenID,
				},
				RequestID:  event.RequestID,
				Address:    event.Address,
				UserAgent:  event.UserAgent,
				Operation:  event.Operation,
				Locale:     event.Locale,
				Outcome:    audit.OutcomeSucceeded,
				OccurredAt: event.OccurredAt,
				EntityType: event.EntityType,
				Reads:      &audit.Reads{Count: event.Count, IDs: event.EntityIDs},
				EntityID:   "",
			}
			// A read of exactly one row is filed under it, like a change.
			if event.Count == 1 && len(event.EntityIDs) == 1 {
				entry.EntityID = event.EntityIDs[0]
			}

			err := audit.Write(ctx, client, entry, nil)
			if err != nil {
				return errs.Wrapf(err, "record read of personal data")
			}

			return nil
		},
	)
}
