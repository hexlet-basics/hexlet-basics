# typed: strict
# frozen_string_literal: true

module StateConcern
  extend T::Sig
  sig { params(base: T.untyped).returns(T.untyped) }
  def self.included(base)
    base.send :include, AASM
    base.send :include, InstanceMethods

    base.attribute :state_event, :string

    base.before_save :set_state
  end

  module InstanceMethods
    extend T::Sig
    sig { returns(T.untyped) }
    def set_state
      # mixed into AASM models with a dynamically-added `state_event` attribute,
      # neither of which Sorbet can see on this bare module — escape via T.unsafe.
      record = T.unsafe(self)
      record.aasm(:state).fire record.state_event.to_sym if record.state_event
    end
  end
end
