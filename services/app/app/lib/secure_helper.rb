class SecureHelper
  class << self
    def generate_token
      SecureRandom.hex(10)
    end
  end
end
