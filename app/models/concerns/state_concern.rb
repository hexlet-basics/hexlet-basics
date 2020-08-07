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
      aasm(:state).fire state_event if state_event
    end
  end
end
