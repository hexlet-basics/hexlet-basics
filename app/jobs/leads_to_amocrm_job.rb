# typed: strict

class LeadsToAmocrmJob < ApplicationJob
  prepend RailsEventStore::AsyncHandler
  extend T::Sig

  retry_on Faraday::Error,
          wait: :polynomially_longer,
          attempts: 5

  def perform(event)
    n8n_client = T.let(
      ApplicationContainer[:n8n_client],
      N8nClient
    )

    serializer = T.let(
      WorkflowLeadSerializer.new(event),
      WorkflowLeadSerializer
    )

    n8n_client.trigger_lead_created_workflow(serializer)
  end
end
