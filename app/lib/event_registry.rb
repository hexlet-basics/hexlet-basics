class EventRegistry
  EVENTS_PATH = Rails.root.join("app", "events")

  def initialize
    @event_names = scan_event_names
  end

  def all
    @event_names
  end

  def options_for_select
    @event_names.map do |event_name|
      [ humanize_event_name(event_name), event_name ]
    end
  end

  def valid?(event_name)
    @event_names.include?(event_name)
  end

  private

  def scan_event_names
    Dir[EVENTS_PATH.join("*_event.rb")].map do |file|
      File.basename(file, ".rb").camelize
    end.sort
  end

  def humanize_event_name(class_name)
    class_name.sub(/Event\z/, "").underscore.humanize
  end
end
