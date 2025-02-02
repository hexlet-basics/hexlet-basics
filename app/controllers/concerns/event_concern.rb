module EventConcern
  extend ActiveSupport::Concern

  included do
    # before_action :fire_js_events
    after_action :clean_js_events

    inertia_share do
      event = session[:event]
      { event: }
    end
  end

  def event_to_js(event)
    eventData = { event_type: event.event_type, data: event.data }
    session[:event] = eventData
    # EventsMapping.exists!(event_name)
    # session[:fired_events][event_name] = options
  end

  def clean_js_events
    # NOTE: могут быть двойные редиректы и стираем данные только
    # при отрисовке хтмл страницы, а не запросу по апи
    if response.successful? # && response.content_type.include?("text/html")
      session[:event] = nil
    end
  end

  # def fire_js_events
  #   session[:fired_events] ||= {}
  #
  #   # gon.push({ fired_events: session[:fired_events] })
  # end
end
