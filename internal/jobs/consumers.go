package jobs

import "hexletbasics/internal/events"

// EventConsumers are the jobs domain facts queue (ADR-0031). The process that
// publishes passes them to events.NewPublisher; the worker runs their modules.
func EventConsumers() []events.Consumer {
	return []events.Consumer{amoCRMLead, recordPersonalDataRead}
}
