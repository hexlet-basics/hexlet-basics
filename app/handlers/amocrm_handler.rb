# frozen_string_literal: true
# typed: true

class AmocrmHandler < ApplicationJob
  extend T::Sig

  prepend RailsEventStore::AsyncHandler

  def perform(event)
    call(event)
  end

  sig { params(event: RailsEventStore::Event).void }
  def call(event)
    case event
    when LeadCreatedEvent
      lead = Lead.find_by(id: event.lead_id)
      return unless lead

      AmocrmService.create_lead(lead: lead)
    end
  end
end
