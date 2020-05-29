# frozen_string_literal: true

class Upload < ApplicationRecord
  include AASM

  aasm :state do
    state :pending, initial: true
    state :successful
    state :failed

    event :success do
      transitions from: :pending, to: :successful
    end

    event :fail do
      transitions from: :pending, to: :failed
    end
  end
end
