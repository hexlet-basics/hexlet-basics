class ApplicationJob < ActiveJob::Base
  include Sentry::Rails::ActiveJobExtensions

  if Rails.env.production?
    retry_on StandardError, wait: :polynomially_longer
  end
  # Automatically retry jobs that encountered a deadlock
  # retry_on ActiveRecord::Deadlocked

  # Most jobs are safe to ignore if the underlying records are no longer available
  # discard_on ActiveJob::DeserializationError
end
