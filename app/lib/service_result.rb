# frozen_string_literal: true

# rubocop:disable Style/MissingRespondToMissing
class ServiceResult
  class << self
    def fail(payload = {})
      new(false, payload)
    end

    def success(payload = {})
      new(true, payload)
    end
  end

  attr_reader :payload

  def initialize(result_status, payload = {})
    @result_status = result_status
    @payload = payload
  end

  def success?
    @result_status == true
  end

  def fail?
    @result_status == false
  end

  def method_missing(method)
    return @payload[method] if @payload.key? method

    raise "method '#{method}' not found"
  end
end
# rubocop:enable Style/MissingRespondToMissing
