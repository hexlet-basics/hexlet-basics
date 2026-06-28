# typed: strict
# frozen_string_literal: true

class SecureHelper
  class << self
    extend T::Sig

    sig { returns(String) }
    def generate_token
      SecureRandom.hex(10)
    end
  end
end
