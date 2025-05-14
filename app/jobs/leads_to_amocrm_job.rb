class LeadsToAmocrmJob < ActiveJob::Base
  prepend RailsEventStore::AsyncHandler

  def perform(event)
    # Implement
    # raise event.inspect
    # payload = WorkflowLeadSerializer.new(event...)
    # client.send(payload)
  end
end
