# typed: true
# frozen_string_literal: true

module StateConcern
  def self.included(base)
    base.send :include, AASM
    base.send :include, InstanceMethods

    base.attribute :state_event, :string

    base.before_save :set_state
  end

  module InstanceMethods
    def set_state
      # mixed into AASM models with a dynamically-added `state_event` attribute,
      # neither of which Sorbet can see on this bare module — escape via T.unsafe.
      record = T.unsafe(self)
      record.aasm(:state).fire record.state_event.to_sym if record.state_event
    end
  end
end
