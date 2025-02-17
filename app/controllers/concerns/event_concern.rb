module EventConcern
  extend ActiveSupport::Concern

  included do |base|
    if base == Web::ApplicationController
      after_action :clean_js_events

      inertia_share do
        events = session[:events]
        { events: }
      end
    end
  end

  def event_to_js(event)
    session[:events] ||= []
    eventData = {
      id: event.event_id,
      type: event.event_type,
      data: event.data,
      metadata: event.metadata
    }
    session[:events] << eventData
    # EventsMapping.exists!(event_name)
    # session[:fired_events][event_name] = options
  end

  def clean_js_events
    # NOTE: могут быть двойные редиректы и стираем данные только
    # при отрисовке хтмл страницы, а не запросу по апи
    if response.successful? # && response.content_type.include?("text/html")
      session[:events] = nil
    end
  end

  # def publish_event(event, user)
  #   event_store.publish(event, stream_name: "user-#{user.id}")
  # end

  # def event_store
  #   Rails.configuration.event_store
  # end
end
