class LeadsToAmocrmJob < ApplicationJob
  prepend RailsEventStore::AsyncHandler

  def perform(event)
    n8n_client = ApplicationContainer[:n8n_client]
    payload    = WorkflowLeadSerializer.new(event).payload

    n8n_client.send(payload)
  end
end
