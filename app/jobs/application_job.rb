# frozen_string_literal: true

class ApplicationJob < ActiveJob::Base
  # Automatically retry jobs that encountered a deadlock
  # retry_on ActiveRecord::Deadlocked
  retry_on Faraday::Error,
           wait: ->(execution_number) { (5 * 2**(execution_number-1)).seconds },
           attempts: 5

  # Most jobs are safe to ignore if the underlying records are no longer available
  # discard_on ActiveJob::DeserializationError
end
