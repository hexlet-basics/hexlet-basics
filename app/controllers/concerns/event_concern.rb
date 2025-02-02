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
    eventData = { event_type: event.event_type, data: event.data }
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
end
