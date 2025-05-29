module EventConcern
  extend ActiveSupport::Concern

  included do |base|
    if base == Web::ApplicationController
      after_action :clean_js_events

      inertia_share do
        if session[:happened_event_ids]
          events = load_events_from_store(session[:happened_event_ids])
          {
            happenedEvents: events.map do |event|
              {
                id: event.event_id,
                type: event.event_type,
                data: event.data
              }
            end
          }
        else
          { happenedEvents: [] }
        end
      end
    end
  end

  def event_to_js(event)
    session[:happened_event_ids] ||= []
    session[:happened_event_ids] << event.event_id
  end

  def clean_js_events
    # Стираем только при HTML-ответе (не API)
    if response.successful?
      session[:happened_event_ids] = nil
    end
  end

  def publish_event(event, user)
    if user.guest?
      event_store.publish(event)
    else
      event_store.with_metadata(user_id: user.id) do
        event_store.publish(
          event,
          stream_name: "user-#{user.id}",
        )
      end
    end

    event_to_js(event)
  end

  def event_store
    Rails.configuration.event_store
  end

  def load_events_from_store(event_ids)
    event_ids.map { |id| event_store.read.event(id) }
  end
end
