# typed: true
# frozen_string_literal: true

class EventSender
  extend T::Sig

  class << self
    extend T::Sig

    sig { params(
      event: ApplicationEvent,
      user: T.nilable(User),
      request: T.nilable(ActionDispatch::Request)
    ).void }
    def publish_event(event, user, request = nil)
      if user.nil?
        event_store.publish(event)
      else
        event_with_metadata = event.class.new(
          data: event.data,
          event_id: event.event_id,
          metadata: event.metadata.to_h.merge(user_id: user.id),
        )

        event_store.publish(
          event_with_metadata,
          stream_name: "user-#{user.id}",
        )
      end
    end

    sig { returns(RailsEventStore::Client) }
    def event_store
      Rails.configuration.event_store
    end
  end
end
