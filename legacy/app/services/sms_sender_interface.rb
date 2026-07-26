# typed: strict
# frozen_string_literal: true

# Provider-agnostic SMS sending. Concrete adapters (Stub, Smsc, ...) implement `deliver`.
# Wired through DI (DepsLocator.current.sms_sender) so the provider can be swapped per env.
class SmsSenderInterface
  extend T::Sig
  extend T::Helpers

  abstract!

  sig { abstract.params(phone: String, text: String).void }
  def self.deliver(phone:, text:); end
end
