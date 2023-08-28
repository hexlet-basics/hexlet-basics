# frozen_string_literal: true

class EventsMapping
  # all events should goes here
  @mapping = {
    signed_up: '',
    signed_in: '',
    next_lesson: '',
    lesson_finished: '',
    language_finished: ''
  }

  def self.events
    @mapping.keys
  end

  def self.exists?(key)
    @mapping.key?(key.to_sym)
  end

  def self.exists!(key)
    raise "Event Name '#{key}' does not exists in EventMapper" unless exists?(key)
  end
end
