# typed: strict

class EventRegistry
  extend T::Sig

  EVENTS_PATH = T.let(Rails.root.join("app", "events"), Pathname)

  sig { void }
  def initialize
    @event_names = T.let(scan_event_names, T::Array[String])
  end

  sig { returns(T::Array[String]) }
  def all
    @event_names
  end

  sig { returns(T::Array[[ String, String ]]) }
  def options_for_select
    @event_names.map do |event_name|
      [ humanize_event_name(event_name), event_name ]
    end
  end

  sig { params(event_name: T.untyped).returns(T::Boolean) }
  def valid?(event_name)
    @event_names.include?(event_name)
  end

  private

  sig { returns(T::Array[String]) }
  def scan_event_names
    Dir[EVENTS_PATH.join("*_event.rb")].map do |file|
      File.basename(file, ".rb").camelize
    end.sort
  end

  sig { params(class_name: String).returns(String) }
  def humanize_event_name(class_name)
    class_name.sub(/Event\z/, "").underscore.humanize
  end
end
