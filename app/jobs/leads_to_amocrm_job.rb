# typed: strict

class LeadsToAmocrmJob < ApplicationJob
  prepend RailsEventStore::AsyncHandler

  retry_on Faraday::Error,
          wait: :polynomially_longer,
          attempts: 5

  def perform(event)
    n8n_client = DepsLocator.current.n8n_client
    serializer = WorkflowLeadSerializer.new(event)
    n8n_client.trigger_lead_created_workflow(serializer)
  end
end
