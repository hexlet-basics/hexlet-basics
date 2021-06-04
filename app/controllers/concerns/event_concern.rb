# frozen_string_literal: true

module EventConcern
  extend ActiveSupport::Concern

  included do
    before_action :fire_js_events
    after_action :clean_js_events
  end

  def js_event(event_name, options = {})
    EventsMapping.exists!(event_name)
    session[:fired_events][event_name] = options
  end

  def clean_js_events
    # NOTE: могут быть двойные редиректы
    if response.successful?
      session[:fired_events] = {}
    end
  end

  def fire_js_events
    session[:fired_events] ||= {}
    gon.push({ fired_events: session[:fired_events] })
  end
end
