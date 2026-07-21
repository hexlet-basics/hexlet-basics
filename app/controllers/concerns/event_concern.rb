
# typed: strict
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

  sig do
    params(
      event: ApplicationEvent,
      user: T.nilable(User),
      request: T.nilable(ActionDispatch::Request)
    ).void
  end
  def publish_event(event, user, request = nil)
    EventSender.publish_event(event, user, request)
  end
end
