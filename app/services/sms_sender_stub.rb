# typed: strict
# frozen_string_literal: true

# Logs the SMS instead of sending it. Default adapter in development/test
# (and the current default everywhere until a real provider is chosen).
class SmsSenderStub < SmsSenderInterface
  extend T::Sig

  sig { override.params(phone: String, text: String).void }
  def self.deliver(phone:, text:)
    Rails.logger.info(Term::ANSIColor.green("[SmsSenderStub] to=#{phone} text=#{text}"))
  end
end
