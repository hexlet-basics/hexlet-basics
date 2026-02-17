
# typed: true
# frozen_string_literal: true

module EventConcern
  extend T::Sig
  extend T::Helpers

  requires_ancestor { ApplicationController }
  extend ActiveSupport::Concern

  # included do
  #   T.bind(self, T.any(
  #     T.class_of(Web::ApplicationController),
  #     T.class_of(Api::LeadsController)
  #   ))
  # end
  sig { params(event: ApplicationEvent).void }
  def js_event_now(event)
    return if I18n.locale == :es
    flash.now.inertia[:events] ||= []
    flash.now.inertia[:events] << event
  end

  sig { params(event: ApplicationEvent).void }
  def js_event(event)
    return if I18n.locale == :es
    flash.inertia[:events] ||= []
    flash.inertia[:events] << event
  end

  sig { params(events: T::Array[ApplicationEvent]).void }
  def js_events(events)
    events.each { js_event(it) }
  end

  def publish_event(...)
    EventSender.publish_event(...)
  end
end
