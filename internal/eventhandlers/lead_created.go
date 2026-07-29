// Package eventhandlers contains Watermill domain-event consumers.
package eventhandlers

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
)

const leadCreatedHandlerName = "lead_created_to_amocrm_v1"

type inserter interface {
	Insert(context.Context, river.JobArgs, *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// LeadCreated routes the domain fact to River; the Watermill consumer stays
// lightweight while River owns external-call retries and execution history.
func LeadCreated(client inserter) cqrs.EventHandler {
	return cqrs.NewEventHandler(
		leadCreatedHandlerName,
		func(ctx context.Context, event *events.LeadCreated) error {
			_, err := client.Insert(ctx, jobs.AmoCRMLeadArgs{Event: *event}, &river.InsertOpts{
				UniqueOpts: river.UniqueOpts{ByArgs: true},
			})
			return err
		},
	)
}
