# typed: strict
# frozen_string_literal: true

class SendSmsJob < ApplicationJob
  extend T::Sig

  sig { params(phone: String, text: String).void }
  def perform(phone:, text:)
    DepsLocator.current.sms_sender.deliver(phone:, text:)
  end
end
